package gomocker

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"testing"
	"unsafe"

	"github.com/agiledragon/gomonkey/v2"
)

// Mocker is the major interface for mocker library
// It allowing developers to mock functions or struct methods for unit tests
//	refer to README.md for more details and examples
type Mocker interface {
	// Mock allows one to mock either a function or a struct method visible to the current package
	//   expectFunc pass in the pointer to the function to be mocked
	//   returns an Expecter instance to allow setting up parameter expectations
	Mock(expectFunc any) Expecter
	// Stub allows one to stub either a function or a struct method visible to the current package
	//   expectFunc pass in the pointer to the function to be mocked
	//   returns a Returner instance to allow setting up return expectations
	Stub(expectFunc any) Returner
}

// Expecter is the interface for setting up parameter expectations
//	refer to README.md for more details and examples
type Expecter interface {
	// Expects allows one to setup a list of parameters to be verified during a function or a struct method call
	//   parameters pass in the list of parameters to be verified,
	//     just like how they are normally passed into the original function or struct method
	//   returns a Returner instance to allow setting up return expectations
	Expects(parameters ...any) Returner
	// NotCalled verifies that no call is expected to the underlying function or struct method
	//   the underlying function or struct method cannot be mocked or stubbed again in the same test
	//   this completes the current Mock sequence, as well as overrides any previous mock or stub
	NotCalled()
}

// Returner is the interface for setting up return expectations
//	refer to README.md for more details and examples
type Returner interface {
	// Returns allows one to setup a list of values to be returned after a function or a struct method call
	//   values pass in the list of values to be returned,
	//     just like how they are normally returned from the original function or struct method
	//   returns a Counter instance to allow setting up execution expectations
	Returns(values ...any) Counter
}

// Returner is the interface for setting up execution expectations
//	refer to README.md for more details and examples
type Counter interface {
	// SideEffects allows one to setup callback functions that is called during expectation verification
	//   note that callbacks are added to the mocked entry and are executed sequentially before return
	//   callbacks pass in are generated by either `GeneralCallback` or `ParameterizedCallback`
	//   returns the same Counter instance to allow setting up further execution expectations
	SideEffects(callbacks ...callback) Counter
	// Once allows one to quickly setup only once execution for the current mock or stub
	//   this is equivalent to call Times(1)
	Once() Mocker
	// Once allows one to quickly setup only twice executions for the current mock or stub
	//   this is equivalent to call Times(2)
	Twice() Mocker
	// Times allows one to setup the number of executions for the current mock or stub
	//   count pass in the number of executions expected, and must be a positive number
	Times(count int) Mocker
}

type mockEntry struct {
	parameters []any
	returns    []any
	callbacks  []callback
}

type funcEntry struct {
	name     string
	stub     bool
	expect   int
	actual   int
	nocall   bool
	verified bool
	mocks    []*mockEntry
}

type mocker struct {
	tester  testing.TB
	patches patcher
	entries map[uintptr]*funcEntry
	locker  sync.Locker
	current *funcEntry
	temp    *mockEntry
}

type patcher interface {
	ApplyCore(target, double reflect.Value) *gomonkey.Patches
	Reset()
}

// NewMocker creates a new instance of mocker using the provided tester interface
//	tester simply pass in the Golang testing struct from a test method
func NewMocker(tester testing.TB) Mocker {
	var m = &mocker{
		tester:  tester,
		patches: gomonkey.NewPatches(),
		entries: make(map[uintptr]*funcEntry),
		locker:  &sync.Mutex{},
	}
	m.tester.Cleanup(m.verifyAll)
	m.tester.Helper()
	return m
}

type parameter interface {
	compare(any) bool
}

type anything struct{}

func (a *anything) compare(any) bool { return true }

// Anything creates a parameter matcher that simply bypasses the check
func Anything() parameter {
	return &anything{}
}

type matching[T any] struct {
	matchFunc func(value T) bool
}

func (m *matching[T]) compare(value any) bool {
	return m.matchFunc(value.(T))
}

// Matches creates a parameter matcher using the provided match function
//	matchFunc pass in the function that customizes the check for a particular parameter
//	  the original parameter is provided in the matchFunc's `value` param
//	  returning false would cause the corresponding test to fail
func Matches[T any](matchFunc func(value T) bool) parameter {
	return &matching[T]{
		matchFunc: matchFunc,
	}
}

type funcValue struct {
	_ uintptr
	p unsafe.Pointer
}

func (m *mocker) getReflectPointer(value reflect.Value) uintptr {
	m.tester.Helper()
	return *(*uintptr)((*funcValue)(unsafe.Pointer(&value)).p)
}

func (m *mocker) getFuncPointer(expectFunc any) (uintptr, string) {
	m.tester.Helper()
	var value = reflect.ValueOf(expectFunc)
	var funcPtr = m.getReflectPointer(value)
	var pointer = value.Pointer()
	var funcForPC = runtime.FuncForPC(pointer)
	var name = funcForPC.Name()
	var file, _ = funcForPC.FileLine(pointer)
	return funcPtr, fmt.Sprint(file, ".", name)
}

func (m *mocker) recover(name string) {
	m.tester.Helper()
	var result = recover()
	if result == nil {
		return
	}
	var message string
	var err, ok = result.(error)
	if ok {
		message = err.Error()
	} else {
		message = fmt.Sprint(result)
	}
	m.tester.Errorf("[%v] Mocker panicing recovered: %v", name, message)
}

func (m *mocker) doComparison(name string, calls int, index int, expect any, actual reflect.Value) {
	m.tester.Helper()
	var param, ok = expect.(parameter)
	if ok {
		if !param.compare(actual.Interface()) {
			m.tester.Errorf(
				"[%v] Parameter mismatch at call #%v parameter #%v: matchFunc failed on actual %v",
				name,
				calls,
				index,
				actual.Interface(),
			)
		}
		return
	}
	if expect == nil {
		if actual.IsValid() && !actual.IsNil() {
			m.tester.Errorf(
				"[%v] Parameter mismatch at call #%v parameter #%v: expect %v, actual %v",
				name,
				calls,
				index,
				expect,
				actual.Interface(),
			)
		}
	} else if !reflect.DeepEqual(actual.Interface(), expect) {
		m.tester.Errorf(
			"[%v] Parameter mismatch at call #%v parameter #%v: expect %v, actual %v",
			name,
			calls,
			index,
			expect,
			actual.Interface(),
		)
	}
}

func (m *mocker) compareNormalParameters(name string, calls int, expects []any, actuals []reflect.Value) {
	m.tester.Helper()
	if len(expects) != len(actuals) {
		m.tester.Errorf(
			"[%v] Invalid number of parameters at call #%v: expect %v, actual %v",
			name,
			calls,
			len(expects),
			len(actuals),
		)
		return
	}
	for index, actual := range actuals {
		m.doComparison(name, calls, index+1, expects[index], actual)
	}
}

func (m *mocker) compareVariadicParameters(name string, calls int, expects []any, actuals []reflect.Value) {
	m.tester.Helper()
	for index, actual := range actuals {
		if index != len(actuals)-1 {
			m.doComparison(name, calls, index+1, expects[index], actual)
		} else {
			if actual.Len() != len(expects)-index {
				m.tester.Errorf(
					"[%v] Invalid number of variadic parameters at call #%v: expect %v, actual %v",
					name,
					calls,
					len(expects)-index,
					actual.Len(),
				)
				return
			}
			for i := index; i < len(expects); i++ {
				var expect = expects[i]
				var item = actual.Index(i - index)
				m.doComparison(name, calls, index+1, expect, item)
			}
		}
	}
}

func (m *mocker) returnZeros(funcType reflect.Type) []reflect.Value {
	var count = funcType.NumOut()
	var rets = make([]reflect.Value, 0, count)
	for i := 0; i < count; i++ {
		rets = append(rets, reflect.Zero(funcType.Out(i)))
	}
	return rets
}

func (m *mocker) constructReturns(name string, calls int, funcType reflect.Type, returns []any) []reflect.Value {
	m.tester.Helper()
	var count = funcType.NumOut()
	if count != len(returns) {
		m.tester.Errorf(
			"[%v] Invalid number of returns at call #%v: expect %v, actual %v",
			name,
			calls,
			count,
			len(returns),
		)
		return m.returnZeros(funcType)
	}
	var rets = []reflect.Value{}
	for i, ret := range returns {
		if ret == nil {
			rets = append(rets, reflect.Zero(funcType.Out(i)))
		} else {
			rets = append(rets, reflect.ValueOf(ret))
		}
	}
	return rets
}

func (m *mocker) makeFunc(name string, funcPtr uintptr, funcType reflect.Type) reflect.Value {
	m.tester.Helper()
	return reflect.MakeFunc(
		funcType,
		func(args []reflect.Value) []reflect.Value {
			m.tester.Helper()
			defer m.recover(name)
			var entry, found = m.entries[funcPtr]
			if !found {
				m.tester.Fatalf(
					"The underlying function or method %v was never setup",
					name,
				)
				return nil
			}
			entry.actual++
			if entry.actual > entry.expect || entry.actual > len(entry.mocks) {
				if !entry.stub {
					m.tester.Errorf(
						"[%v] Unepxected number of calls: expect %v, actual %v",
						name,
						entry.expect,
						entry.actual,
					)
					entry.verified = true
					return m.returnZeros(funcType)
				}
				entry.actual = len(entry.mocks)
			}
			var mock = entry.mocks[entry.actual-1]
			if !entry.stub {
				if funcType.IsVariadic() {
					m.compareVariadicParameters(name, entry.actual, mock.parameters, args)
				} else {
					m.compareNormalParameters(name, entry.actual, mock.parameters, args)
				}
			}
			for _, callback := range mock.callbacks {
				switch callback.(type) {
				case *generalCallback:
					callback.execute(entry.actual, 0, nil)
				default:
					for paramIndex, arg := range args {
						callback.execute(entry.actual, paramIndex+1, arg.Interface())
					}
				}
			}
			return m.constructReturns(name, entry.actual, funcType, mock.returns)
		},
	)
}

func (m *mocker) setup(name string, stub bool, funcPtr uintptr) {
	m.tester.Helper()
	if m.current != nil || m.temp != nil {
		m.tester.Fatalf(
			"A former setup for function or method [%v] was incomplete."+
				" Did you miss calling the Once/Twice/Times method in the end?",
			name,
		)
		return
	}
	var entry, found = m.entries[funcPtr]
	if found {
		if entry.stub != stub {
			if entry.stub {
				m.tester.Fatalf(
					"A former setup for function or method [%v] was a Stub but current setup is a Mock."+
						" We do not support mixing Stub and Mock for the same function or method at the moment.",
					name,
				)
			} else {
				m.tester.Fatalf(
					"A former setup for function or method [%v] was a Mock but current setup is a Stub."+
						" We do not support mixing Stub and Mock for the same function or method at the moment.",
					name,
				)
			}
			return
		}
		if entry.nocall {
			m.tester.Fatalf("A former setup for function or method [%v] was to be not called,"+
				" therefore no more Mock or Stub can be setup for it now.",
				name,
			)
			return
		}
		m.current = entry
		m.temp = &mockEntry{}
		return
	}
	entry = &funcEntry{
		name:   name,
		stub:   stub,
		actual: 0,
		mocks:  make([]*mockEntry, 0),
	}
	m.entries[funcPtr] = entry
	m.current = entry
	m.temp = &mockEntry{}
}

// Mock allows one to mock either a function or a struct method visible to the current package
//	expectFunc pass in the pointer to the function to be mocked
//	returns an Expecter instance to allow setting up parameter expectations
func (m *mocker) Mock(expectFunc any) Expecter {
	m.tester.Helper()
	m.locker.Lock()
	defer m.locker.Unlock()
	var funcPtr, name = m.getFuncPointer(expectFunc)
	var funcType = reflect.TypeOf(expectFunc)
	m.setup(name, false, funcPtr)
	m.patches.ApplyCore(
		reflect.ValueOf(expectFunc),
		m.makeFunc(name, funcPtr, funcType),
	)
	return m
}

// Stub allows one to stub either a function or a struct method visible to the current package
//	expectFunc pass in the pointer to the function to be mocked
//	returns a Returner instance to allow setting up return expectations
func (m *mocker) Stub(expectFunc any) Returner {
	m.tester.Helper()
	m.locker.Lock()
	defer m.locker.Unlock()
	var funcPtr, name = m.getFuncPointer(expectFunc)
	var funcType = reflect.TypeOf(expectFunc)
	m.setup(name, true, funcPtr)
	m.patches.ApplyCore(
		reflect.ValueOf(expectFunc),
		m.makeFunc(name, funcPtr, funcType),
	)
	return m
}

// Expects allows one to setup a list of parameters to be verified during a function or a struct method call
//	parameters pass in the list of parameters to be verified,
//	  just like how they are normally passed into the original function or struct method
//	returns a Returner instance to allow setting up return expectations
func (m *mocker) Expects(parameters ...any) Returner {
	m.tester.Helper()
	if m.current == nil || m.temp == nil {
		m.tester.Fatalf(
			"Unexpected call to Expects without setting up an anticipated function or method",
		)
		return m
	}
	m.temp.parameters = parameters
	return m
}

// NotCalled verifies that no call is expected to the underlying function or struct method
//	the underlying function or struct method cannot be mocked or stubbed again in the same test
//	this completes the current Mock sequence, as well as overrides any previous mock or stub
func (m *mocker) NotCalled() {
	m.tester.Helper()
	if m.current == nil || m.temp == nil {
		m.tester.Fatalf(
			"Unexpected call to NotCalled without setting up an anticipated function or method",
		)
		return
	}
	m.current.nocall = true
	m.current.expect = 0
	m.current.mocks = []*mockEntry{{}}
	m.temp = nil
	m.current = nil
}

// Returns allows one to setup a list of values to be returned after a function or a struct method call
//	values pass in the list of values to be returned,
//	  just like how they are normally returned from the original function or struct method
//	returns a Counter instance to allow setting up execution expectations
func (m *mocker) Returns(values ...any) Counter {
	m.tester.Helper()
	if m.current == nil || m.temp == nil {
		m.tester.Fatalf(
			"Unexpected call to Returns without setting up an anticipated function or method",
		)
		return m
	}
	m.temp.returns = values
	return m
}

type callback interface {
	execute(int, int, any)
}

type generalCallback struct {
	callIndex    int
	callbackFunc func()
}

func (c *generalCallback) execute(callIndex int, paramIndex int, paramValue any) {
	if c.callIndex > 0 && c.callIndex != callIndex {
		return
	}
	c.callbackFunc()
}

// GeneralSideEffect creates a general side-effect with callback function to be registered
//  callIndex limits when the callback should happen based on the current number of call
//    if a non-positive number is provided, the side effect is happening for each call
//  this value is 1-based, i.e. the first call would yield 1 for callIndex, etc.
//    callbackFunc is the function to be called when the side effect takes place
func GeneralSideEffect(callIndex int, callbackFunc func()) callback {
	return &generalCallback{
		callIndex:    callIndex,
		callbackFunc: callbackFunc,
	}
}

type parameterizedCallback[T any] struct {
	callIndex    int
	paramIndex   int
	callbackFunc func(T)
}

func (c *parameterizedCallback[T]) execute(callIndex int, paramIndex int, paramValue any) {
	if c.callIndex > 0 && c.callIndex != callIndex {
		return
	}
	if c.paramIndex != paramIndex {
		return
	}
	c.callbackFunc(paramValue.(T))
}

// ParamSideEffect creates a parameterized side-effect with callback function to be registered
//  callIndex limits when the side effect should happen based on the current number of call
//    if a non-positive number is provided, the side effect is happening for each call
//    this value is 1-based, i.e. the first call would yield 1 for callIndex, etc.
//  paramIndex identifies which parameter to be operated for the side effect
//    this value is 1-based, i.e. the first parameter would yield 1 for paramIndex, etc.
//  callbackFunc is the function to be called when the side effect takes place
//    to which the parameter of the paramIndex is passed as the parameter to the function
func ParamSideEffect[T any](callIndex int, paramIndex int, callbackFunc func(T)) callback {
	return &parameterizedCallback[T]{
		callIndex:    callIndex,
		paramIndex:   paramIndex,
		callbackFunc: callbackFunc,
	}
}

// SideEffects allows one to setup callback functions that are called during expectation verification
//	note that callbacks are added to the mocked entry and are executed sequentially before return
//	callbacks pass in are generated by either `GeneralCallback` or `ParameterizedCallback`
//	returns the same Counter instance to allow setting up further execution expectations
func (m *mocker) SideEffects(callbacks ...callback) Counter {
	m.tester.Helper()
	if m.current == nil || m.temp == nil {
		m.tester.Fatalf(
			"Unexpected call to SideEffect without setting up an anticipated function or method",
		)
		return m
	}
	m.temp.callbacks = append(m.temp.callbacks, callbacks...)
	return m
}

// Once allows one to quickly setup only once execution for the current mock or stub
//	this is equivalent to call Times(1)
func (m *mocker) Once() Mocker {
	m.tester.Helper()
	return m.Times(1)
}

// Once allows one to quickly setup only twice executions for the current mock or stub
//	this is equivalent to call Times(2)
func (m *mocker) Twice() Mocker {
	m.tester.Helper()
	return m.Times(2)
}

// Times allows one to setup the number of executions for the current mock or stub
//	count pass in the number of executions expected, and must be a positive number
func (m *mocker) Times(count int) Mocker {
	m.tester.Helper()
	if m.current == nil || m.temp == nil {
		m.tester.Fatalf(
			"Unexpected call to Times without setting up an anticipated function or method",
		)
		return m
	}
	if count < 0 {
		m.tester.Fatalf(
			"function or method [%v] cannot be mocked for negative [%v] times",
			m.current.name,
			count,
		)
		return m
	} else if count == 0 {
		m.tester.Fatalf(
			"function or method [%v] cannot be mocked for zero times using Times method."+
				" Try using NotCalled method instead.",
			m.current.name,
		)
		return m
	}
	m.current.expect += count
	for i := 0; i < count; i++ {
		m.current.mocks = append(m.current.mocks, m.temp)
	}
	m.temp = nil
	m.current = nil
	return m
}

func (m *mocker) verifyAll() {
	m.tester.Helper()
	for _, entry := range m.entries {
		if entry.verified {
			continue
		}
		if !entry.stub && entry.expect != entry.actual {
			m.tester.Errorf(
				"[%v] Unepxected number of calls: expect %v, actual %v",
				entry.name,
				entry.expect,
				entry.actual,
			)
		}
	}
	m.entries = make(map[uintptr]*funcEntry)
	m.patches.Reset()
}
