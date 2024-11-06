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
//
// It allowing developers to mock functions or struct methods for unit tests
//
//	refer to README.md for more details and examples
type Mocker interface {
	// Mock allows one to mock either a function or a struct method visible to the current package
	//
	//   expectFunc pass in the pointer to the function to be mocked
	//   returns an Expecter instance to allow setting up parameter expectations
	Mock(expectFunc interface{}) Expecter
	// Stub allows one to stub either a function or a struct method visible to the current package
	//
	//   expectFunc pass in the pointer to the function to be mocked
	//   returns a Returner instance to allow setting up return expectations
	Stub(expectFunc interface{}) Returner
}

// Expecter is the interface for setting up parameter expectations
//
//	refer to README.md for more details and examples
type Expecter interface {
	// Expects allows one to setup a list of parameters to be verified during a function or a struct method call
	//
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
//
//	refer to README.md for more details and examples
type Returner interface {
	// Returns allows one to setup a list of values to be returned after a function or a struct method call
	//
	//   values pass in the list of values to be returned,
	//     just like how they are normally returned from the original function or struct method
	//   returns a Counter instance to allow setting up execution expectations
	Returns(values ...any) Counter
}

// Returner is the interface for setting up execution expectations
//
//	refer to README.md for more details and examples
type Counter interface {
	// SideEffect allows one to setup a callback function that is called during expectation verification
	//   note that there is only one side effect for each mock or stub, and the newest overrides previous ones
	//
	//   callback pass in the customized callback function with an integer parameter `index`
	//     this parameter indicates the number of executions done so far including the current one
	//     and `params` are the exact arguments passed into the underlying function or struct method
	//   returns the same Counter instance to allow setting up further execution expectations
	SideEffect(callback func(index int, params ...interface{})) Counter
	// Once allows one to quickly setup only once execution for the current mock or stub
	//   this is equivalent to call Times(1)
	Once() Mocker
	// Once allows one to quickly setup only twice executions for the current mock or stub
	//   this is equivalent to call Times(2)
	Twice() Mocker
	// Times allows one to setup the number of executions for the current mock or stub
	//
	//   count pass in the number of executions expected, and must be a positive number
	Times(count int) Mocker
}

type mockEntry struct {
	parameters []interface{}
	returns    []interface{}
	callback   func(int, ...interface{})
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
//
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

type parameter struct {
	isAnything bool
	matchFunc  func(value interface{}) bool
}

// Anything creates a parameter matcher that simply bypasses the check
func Anything() *parameter {
	return &parameter{
		isAnything: true,
	}
}

// Matches creates a parameter matcher using the provided match function
//
//	matchFunc pass in the function that customizes the check for a particular parameter
//	  the original parameter is wrapped into an interface and is given as `value` here
//	  returning false would cause the corresponding test to fail
func Matches(matchFunc func(value interface{}) bool) *parameter {
	return &parameter{
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

func (m *mocker) getFuncPointer(expectFunc interface{}) (uintptr, string) {
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

func (m *mocker) doComparison(name string, calls int, index int, expect interface{}, actual reflect.Value) {
	m.tester.Helper()
	var param, ok = expect.(*parameter)
	if !ok {
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
		return
	}
	if param.isAnything {
		return
	}
	if param.matchFunc != nil {
		if !param.matchFunc(actual.Interface()) {
			m.tester.Errorf(
				"[%v] Parameter mismatch at call #%v parameter #%v: matchFunc failed on actual %v",
				name,
				calls,
				index,
				actual.Interface(),
			)
		}
	}
}

func (m *mocker) compareNormalParameters(name string, calls int, expects []interface{}, actuals []reflect.Value) {
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

func (m *mocker) compareVariadicParameters(name string, calls int, expects []interface{}, actuals []reflect.Value) {
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

func (m *mocker) constructReturns(name string, calls int, funcType reflect.Type, returns []interface{}) []reflect.Value {
	m.tester.Helper()
	var count = funcType.NumOut()
	if count != len(returns) {
		m.tester.Fatalf(
			"[%v] Invalid number of returns at call #%v: expect %v, actual %v",
			name,
			calls,
			count,
			len(returns),
		)
		return nil
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
					m.tester.Fatalf(
						"[%v] Unepxected number of calls: expect %v, actual %v",
						name,
						entry.expect,
						entry.actual,
					)
					entry.verified = true
					return nil
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
			if mock.callback != nil {
				var params = []interface{}{}
				for _, arg := range args {
					params = append(params, arg.Interface())
				}
				mock.callback(entry.actual, params...)
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
//
//	expectFunc pass in the pointer to the function to be mocked
//	returns an Expecter instance to allow setting up parameter expectations
func (m *mocker) Mock(expectFunc interface{}) Expecter {
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
//
//	expectFunc pass in the pointer to the function to be mocked
//	returns a Returner instance to allow setting up return expectations
func (m *mocker) Stub(expectFunc interface{}) Returner {
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
//
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
//
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
//
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

// SideEffect allows one to setup a callback function that is called during expectation verification
//
//	note that there is only one side effect for each mock or stub, and the newest overrides previous ones
//
//	callback pass in the customized callback function with an integer parameter `index`
//	  this parameter indicates the number of executions done so far including the current one
//	  and `params` are the exact arguments passed into the underlying function or struct method
//	returns the same Counter instance to allow setting up further execution expectations
func (m *mocker) SideEffect(callback func(index int, params ...interface{})) Counter {
	m.tester.Helper()
	if m.current == nil || m.temp == nil {
		m.tester.Fatalf(
			"Unexpected call to SideEffect without setting up an anticipated function or method",
		)
		return m
	}
	m.temp.callback = callback
	return m
}

// Once allows one to quickly setup only once execution for the current mock or stub
//
//	this is equivalent to call Times(1)
func (m *mocker) Once() Mocker {
	m.tester.Helper()
	return m.Times(1)
}

// Once allows one to quickly setup only twice executions for the current mock or stub
//
//	this is equivalent to call Times(2)
func (m *mocker) Twice() Mocker {
	m.tester.Helper()
	return m.Times(2)
}

// Times allows one to setup the number of executions for the current mock or stub
//
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
