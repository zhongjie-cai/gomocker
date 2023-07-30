package mocker

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"testing"
	"unsafe"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/agiledragon/gomonkey/v2/creflect"
	"github.com/stretchr/testify/assert"
)

type Mocker interface {
	ExpectFunc(expectFunc any, count int, mockFunc any)
	ExpectMethod(targetStruct any, expectMethod string, count int, mockMethod any)
	CalledFunc(expectFunc any) int
	CalledMethod(targetStruct any, expectMethod string) int
}

type mocker struct {
	t       *testing.T
	patches *gomonkey.Patches
	called  map[uintptr]int
	expect  map[uintptr]int
	names   map[uintptr]string
	locker  *sync.Mutex
}

func NewMocker(t *testing.T) Mocker {
	var m = &mocker{
		t:		 t,
		patches: gomonkey.NewPatches(),
		called:  make(map[uintptr]int),
		expect:  make(map[uintptr]int),
		names:   make(map[uintptr]string),
		locker:  &sync.Mutex{},
	}
	m.t.Cleanup(m.verifyAll)
	return m
}

type funcValue struct {
	_ uintptr
	p unsafe.Pointer
}

func getReflectPointer(value reflect.Value) uintptr {
	return *(*uintptr)((*funcValue)(unsafe.Pointer(&value)).p)
}

func getFuncPointer(expectFunc any) (uintptr, string) {
	var value = reflect.ValueOf(expectFunc)
	var funcPtr = getReflectPointer(value)
	var pointer = value.Pointer()
	var funcForPC = runtime.FuncForPC(pointer)
	var name = funcForPC.Name()
	var file, _ = funcForPC.FileLine(pointer)
	return funcPtr, fmt.Sprint(file, ".", name)
}

func getTypeName(typeValue reflect.Type) string {
	switch (typeValue.Kind()) {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Pointer, reflect.Slice:
		return fmt.Sprint(typeValue.Elem().PkgPath(), "/", typeValue.Elem().Name())
	}
		return fmt.Sprint(typeValue.PkgPath(), "/", typeValue.Name())
}

func (m *mocker) ExpectFunc(expectFunc any, count int, mockFunc any) {
	m.locker.Lock()
	defer m.locker.Unlock()
	var funcPtr, name = getFuncPointer(expectFunc)
	m.expect[funcPtr] = count
	m.names[funcPtr] = name
	m.patches.ApplyFunc(expectFunc, mockFunc)
}

func (m *mocker) getMethodPointer(targetStruct any, expectMethod string) (uintptr, string) {
	var typeValue, ok = targetStruct.(reflect.Type)
	if !ok {
		typeValue = reflect.TypeOf(targetStruct)
	}
	var typeName = getTypeName(typeValue)
	var method, success = typeValue.MethodByName(expectMethod)
	if !success {
		m.t.Skipf(
			"Method [%v] cannot be located for given target struct [%v]. Skipping the current test.",
			expectMethod,
			typeName,
		)
		return 0, ""
	}
	return getReflectPointer(method.Func), fmt.Sprint(typeName, ".", expectMethod)
}

func (m *mocker) getPrivateMethodPointer(targetStruct any, expectMethod string) (uintptr, string) {
	var typeValue, ok = targetStruct.(reflect.Type)
	if !ok {
		typeValue = reflect.TypeOf(targetStruct)
	}
	var typeName = getTypeName(typeValue)
	var funcPtr, success = creflect.MethodByName(typeValue, expectMethod)
	if !success {
		m.t.Skipf(
			"Method [%v] cannot be located for given target struct [%v]. Skipping the current test.",
			expectMethod,
			typeName,
		)
		return 0, ""
	}
	return *(*uintptr)(funcPtr), fmt.Sprint(typeName, ".", expectMethod)
}

func isPrivateMethod(methodName string) bool {
	var firstChar = methodName[0]
	return firstChar >= 'a' && firstChar <= 'z'
}

func (m *mocker) ExpectMethod(targetStruct any, expectMethod string, count int, mockMethod any) {
	m.locker.Lock()
	defer m.locker.Unlock()
	var funcPtr uintptr
	var name string
	if isPrivateMethod(expectMethod) {
		funcPtr, name = m.getPrivateMethodPointer(targetStruct, expectMethod)
		m.patches.ApplyPrivateMethod(targetStruct, expectMethod, mockMethod)
	} else {
		funcPtr, name = m.getMethodPointer(targetStruct, expectMethod)
		m.patches.ApplyMethod(targetStruct, expectMethod, mockMethod)
	}
	m.expect[funcPtr] = count
	m.names[funcPtr] = name
}

func (m *mocker) calledFunc(funcPtr uintptr) int {
	var count, found = m.called[funcPtr]
	if !found {
		count = 1
	} else {
		count++
	}
	m.called[funcPtr] = count
	return count
}

func (m *mocker) CalledFunc(expectFunc any) int {
	m.locker.Lock()
	defer m.locker.Unlock()
	var funcPtr, name = getFuncPointer(expectFunc)
	m.names[funcPtr] = name
	return m.calledFunc(funcPtr)
}

func (m *mocker) CalledMethod(targetStruct any, expectMethod string) int {
	m.locker.Lock()
	defer m.locker.Unlock()
	var funcPtr uintptr
	var name string
	if isPrivateMethod(expectMethod) {
		funcPtr, name = m.getPrivateMethodPointer(targetStruct, expectMethod)
	} else {
		funcPtr, name = m.getMethodPointer(targetStruct, expectMethod)
	}
	m.names[funcPtr] = name
	return m.calledFunc(funcPtr)
}

func (m *mocker) verifyAll() {
	for funcPtr, name := range m.names {
		var expectCount, isExpected = m.expect[funcPtr]
		var callCount, isCalled = m.called[funcPtr]
		if isExpected {
			if isCalled {
				assert.Equal(m.t, expectCount, callCount, "Expected call to method [%v] for %v times but actually called %v times", name, expectCount, callCount)
			} else {
				assert.Fail(m.t, fmt.Sprintf("Expected call to method [%v] for %v times but not called at all", name, expectCount))
			}
		} else {
			if isCalled {
				assert.Fail(m.t, fmt.Sprintf("Unexpected call to method [%v] for %v times", name, callCount))
			} else {
				// success
			}
		}
	}
	m.patches.Reset()
}
