package gomocker

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"testing"
	"unsafe"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/agiledragon/gomonkey/v2/creflect"
)

// Mocker is the interface for mocker library
//
// It allowing developers to mock either functions or struct methods according to unit test needs
// refer to README.md for examples on how to use this library
type Mocker interface {
	// ExpectFunc allows one to mock either a public or private function visible to the current package
	//
	//   expectFunc pass in the pointer to the function to be mocked
	//   count indicates the number of calls for the expectFunc during test execution; zero can be provided but must be the first expectation; negative values are treated as zeros
	//   mockFunc pass in the pointer to the function to be actually called during test execution
	//   returns the mocker instance itself to allow fluent calls to it
	ExpectFunc(expectFunc interface{}, count int, mockFunc interface{}) Mocker

	// FuncCalledCount checks the number of calls made to the expected function as of the time this method is executed
	//
	//	expectFunc pass in the pointer to the function to be mocked
	//	returns the number of calls made to the expected function
	FuncCalledCount(expectFunc interface{}) int

	// ExpectMethod allows one to mock either a public or private method associated to a struct or interface visible to the current package
	//
	//   targetStruct pass in the pointer to the struct or interface instance to be mocked
	//   expectMethod pass in the name of the method to be mocked
	//   count indicates the number of calls for the expectFunc during test execution; zero can be provided but must be the first expectation; negative values are treated as zeros
	//   mockFunc pass in the pointer to the function to be actually called during test execution;
	//     due to language specs, one additional parameter is expected as the first parameter in the method signature, reflecting the struct pointer or value itself
	//   returns the mocker instance itself to allow fluent calls to it
	ExpectMethod(targetStruct interface{}, expectMethod string, count int, mockMethod interface{}) Mocker

	// MethodCalledCount checks the number of calls made to the expected method as of the time this method is executed
	//
	//	targetStruct pass in the pointer to the struct or interface instance to be mocked
	//	expectMethod pass in the name of the method to be mocked
	//	returns the number of calls made to the expected method
	MethodCalledCount(targetStruct interface{}, expectMethod string) int
}

type funcEntry struct {
	name   string
	expect int
	actual int
	method []reflect.Value
}

type mocker struct {
	tester  testing.TB
	patches patcher
	entries map[uintptr]*funcEntry
	locker  sync.Locker
}

type patcher interface {
	ApplyCore(target, double reflect.Value) *gomonkey.Patches
	ApplyCoreOnlyForPrivateMethod(target unsafe.Pointer, double reflect.Value) *gomonkey.Patches
	Reset()
}

// NewMocker creates a new instance of mocker using the provided tester interface
//
//	tester: simply pass in the Golang testing struct from a test method
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
	var result = recover()
	if result == nil {
		return
	}
	m.tester.Helper()
	var message string
	var err, ok = result.(error)
	if ok {
		message = err.Error()
	} else {
		message = fmt.Sprint(result)
	}
	m.tester.Errorf("[%v] Mocker panicing recovered: %v", name, message)
}

func (m *mocker) getTypeName(typeValue reflect.Type) string {
	m.tester.Helper()
	switch typeValue.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		return fmt.Sprint(typeValue.Elem().PkgPath(), ".", typeValue.Elem().Name())
	}
	return fmt.Sprint(typeValue.PkgPath(), ".", typeValue.Name())
}

func (m *mocker) makeFunc(name string, funcPtr uintptr, mockFunc interface{}) reflect.Value {
	m.tester.Helper()
	var funcType = reflect.TypeOf(mockFunc)
	return reflect.MakeFunc(
		funcType,
		func(args []reflect.Value) []reflect.Value {
			m.tester.Helper()
			defer m.recover(name)
			var entry, found = m.entries[funcPtr]
			var funcValue = reflect.ValueOf(mockFunc)
			if !found {
				entry = &funcEntry{
					name:   name,
					expect: 0,
					actual: 1,
					method: make([]reflect.Value, 0),
				}
				m.entries[funcPtr] = entry
			} else {
				entry.actual++
				if entry.actual <= entry.expect {
					funcValue = entry.method[entry.actual-1]
				}
			}
			if funcType.IsVariadic() {
				return funcValue.CallSlice(args)
			}
			return funcValue.Call(args)
		},
	)
}

func (m *mocker) setupExpect(name string, funcPtr uintptr, count int, mockFunc interface{}) {
	m.tester.Helper()
	var mockValue = reflect.ValueOf(mockFunc)
	var entry, found = m.entries[funcPtr]
	if found {
		if count <= 0 {
			m.tester.Errorf(
				"[%v] Expect count must be greather than 0 when already mocked with other expectations",
				name,
			)
		} else {
			entry.expect += count
			for i := 0; i < count; i++ {
				entry.method = append(entry.method, mockValue)
			}
		}
		return
	}
	var method = []reflect.Value{}
	if count <= 0 {
		method = append(method, mockValue)
	} else {
		for i := 0; i < count; i++ {
			method = append(method, mockValue)
		}
	}
	entry = &funcEntry{
		name:   name,
		expect: count,
		actual: 0,
		method: method,
	}
	m.entries[funcPtr] = entry
}

// ExpectFunc allows one to mock either a public or private function visible to the current package
//
//	expectFunc pass in the pointer to the function to be mocked
//	count indicates the number of calls for the expectFunc during test execution; zero can be provided but must be the first expectation; negative values are treated as zeros
//	mockFunc pass in the pointer to the function to be actually called during test execution
//	returns the mocker instance itself to allow fluent calls to it
func (m *mocker) ExpectFunc(expectFunc interface{}, count int, mockFunc interface{}) Mocker {
	m.tester.Helper()
	m.locker.Lock()
	defer m.locker.Unlock()
	var funcPtr, name = m.getFuncPointer(expectFunc)
	m.setupExpect(name, funcPtr, count, mockFunc)
	m.patches.ApplyCore(
		reflect.ValueOf(expectFunc),
		m.makeFunc(name, funcPtr, mockFunc),
	)
	return m
}

// FuncCalledCount checks the number of calls made to the expected function as of the time this method is executed
//
//	expectFunc pass in the pointer to the function to be mocked
//	returns the number of calls made to the expected function
func (m *mocker) FuncCalledCount(expectFunc interface{}) int {
	m.tester.Helper()
	m.locker.Lock()
	defer m.locker.Unlock()
	var funcPtr, _ = m.getFuncPointer(expectFunc)
	var entry, found = m.entries[funcPtr]
	if !found {
		return 0
	}
	return entry.actual
}

func (m *mocker) getMethodPointer(targetStruct interface{}, expectMethod string) (uintptr, reflect.Value, string) {
	m.tester.Helper()
	var typeValue, ok = targetStruct.(reflect.Type)
	if !ok {
		typeValue = reflect.TypeOf(targetStruct)
	}
	var typeName = m.getTypeName(typeValue)
	var method, success = typeValue.MethodByName(expectMethod)
	if !success {
		m.tester.Errorf(
			"Method [%v] cannot be located for given target struct [%v]",
			expectMethod,
			typeName,
		)
		return 0, reflect.Value{}, ""
	}
	return m.getReflectPointer(method.Func), method.Func, fmt.Sprint(typeName, ".", expectMethod)
}

func (m *mocker) getPrivateMethodPointer(targetStruct interface{}, expectMethod string) (uintptr, string) {
	m.tester.Helper()
	var typeValue, ok = targetStruct.(reflect.Type)
	if !ok {
		typeValue = reflect.TypeOf(targetStruct)
	}
	var typeName = m.getTypeName(typeValue)
	var funcPtr, success = creflect.MethodByName(typeValue, expectMethod)
	if !success {
		m.tester.Errorf(
			"Method [%v] cannot be located for given target struct [%v]",
			expectMethod,
			typeName,
		)
		return 0, ""
	}
	return *(*uintptr)(funcPtr), fmt.Sprint(typeName, ".", expectMethod)
}

func (m *mocker) isPrivateMethod(methodName string) bool {
	m.tester.Helper()
	var firstChar = methodName[0]
	return firstChar >= 'a' && firstChar <= 'z'
}

// ExpectMethod allows one to mock either a public or private method associated to a struct or interface visible to the current package
//
//	targetStruct pass in the pointer to the struct or interface instance to be mocked
//	expectMethod pass in the name of the method to be mocked
//	count indicates the number of calls for the expectFunc during test execution; zero can be provided but must be the first expectation; negative values are treated as zeros
//	mockFunc pass in the pointer to the function to be actually called during test execution;
//	  due to language specs, one additional parameter is expected as the first parameter in the method signature, reflecting the struct pointer or value itself
//	returns the mocker instance itself to allow fluent calls to it
func (m *mocker) ExpectMethod(targetStruct interface{}, expectMethod string, count int, mockMethod interface{}) Mocker {
	m.tester.Helper()
	m.locker.Lock()
	defer m.locker.Unlock()
	var funcPtr uintptr
	var name string
	if m.isPrivateMethod(expectMethod) {
		funcPtr, name = m.getPrivateMethodPointer(targetStruct, expectMethod)
		if funcPtr == 0 {
			return m
		}
		var target = unsafe.Pointer(&funcPtr)
		m.patches.ApplyCoreOnlyForPrivateMethod(target, m.makeFunc(name, funcPtr, mockMethod))
	} else {
		var funcValue reflect.Value
		funcPtr, funcValue, name = m.getMethodPointer(targetStruct, expectMethod)
		if funcPtr == 0 {
			return m
		}
		m.patches.ApplyCore(funcValue, m.makeFunc(name, funcPtr, mockMethod))
	}
	m.setupExpect(name, funcPtr, count, mockMethod)
	return m
}

// MethodCalledCount checks the number of calls made to the expected method as of the time this method is executed
//
//	targetStruct pass in the pointer to the struct or interface instance to be mocked
//	expectMethod pass in the name of the method to be mocked
//	returns the number of calls made to the expected method
func (m *mocker) MethodCalledCount(targetStruct interface{}, expectMethod string) int {
	m.tester.Helper()
	m.locker.Lock()
	defer m.locker.Unlock()
	var funcPtr uintptr
	if m.isPrivateMethod(expectMethod) {
		funcPtr, _ = m.getPrivateMethodPointer(targetStruct, expectMethod)
	} else {
		funcPtr, _, _ = m.getMethodPointer(targetStruct, expectMethod)
	}
	var entry, found = m.entries[funcPtr]
	if !found {
		return 0
	}
	return entry.actual
}

func (m *mocker) verifyAll() {
	m.tester.Helper()
	for _, entry := range m.entries {
		if entry.expect != entry.actual {
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
