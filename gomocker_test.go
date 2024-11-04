package gomocker

import (
	"errors"
	"math/rand"
	"reflect"
	"sync"
	"testing"
)

func assertEquals(t *testing.T, expect interface{}, actual interface{}, message string) {
	t.Helper()
	if expect == actual {
		return
	}
	t.Errorf(
		"%v: expect %v, actual %v",
		message,
		expect,
		actual,
	)
}

func TestMocker_ShouldStubFunctionOnce(t *testing.T) {
	// arrange
	var foo = func(bar int) int {
		return 0
	}
	var dummyBar = rand.Intn(100)
	var dummyResult = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// expect
	m.Stub(foo).Returns(dummyResult).Once()

	// SUT + act
	var result = foo(dummyBar)

	// assert
	assertEquals(t, dummyResult, result, "foo call result different")
}

func TestMocker_ShouldStubFunctionTwice(t *testing.T) {
	// arrange
	var foo = func(bar int) int {
		return 0
	}
	var dummyBar = rand.Intn(100)
	var dummyResult = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// expect
	m.Stub(foo).Returns(dummyResult).Twice()

	// SUT + act
	var result1 = foo(dummyBar)
	var result2 = foo(dummyBar)

	// assert
	assertEquals(t, dummyResult, result1, "foo call result 1 different")
	assertEquals(t, dummyResult, result2, "foo call result 2 different")
}

func TestMocker_ShouldStubFunctionMultipleTimes(t *testing.T) {
	// arrange
	var foo = func(bar int) int {
		return 0
	}
	var dummyBar = rand.Intn(100)
	var dummyResult = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// expect
	m.Stub(foo).Returns(dummyResult).Times(5)

	// SUT + act
	var result1 = foo(dummyBar)
	var result2 = foo(dummyBar)
	var result3 = foo(dummyBar)
	var result4 = foo(dummyBar)
	var result5 = foo(dummyBar)

	// assert
	assertEquals(t, dummyResult, result1, "foo call result 1 different")
	assertEquals(t, dummyResult, result2, "foo call result 2 different")
	assertEquals(t, dummyResult, result3, "foo call result 3 different")
	assertEquals(t, dummyResult, result4, "foo call result 4 different")
	assertEquals(t, dummyResult, result5, "foo call result 5 different")
}

func TestMocker_ShouldMockFunctionOnce(t *testing.T) {
	// arrange
	var foo = func(bar int) int {
		return 0
	}
	var dummyBar = rand.Intn(100)
	var dummyResult = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// expect
	m.Mock(foo).Expects(dummyBar).Returns(dummyResult).Once()

	// SUT + act
	var result = foo(dummyBar)

	// assert
	assertEquals(t, dummyResult, result, "foo call result different")
}

func TestMocker_ShouldMockFunctionTwice(t *testing.T) {
	// arrange
	var foo = func(bar int) int {
		return 0
	}
	var dummyBar = rand.Intn(100)
	var dummyResult = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// expect
	m.Mock(foo).Expects(dummyBar).Returns(dummyResult).Twice()

	// SUT + act
	var result1 = foo(dummyBar)
	var result2 = foo(dummyBar)

	// assert
	assertEquals(t, dummyResult, result1, "foo call result 1 different")
	assertEquals(t, dummyResult, result2, "foo call result 2 different")
}

func TestMocker_ShouldMockFunctionMultipleTimes(t *testing.T) {
	// arrange
	var foo = func(bar int) int {
		return 0
	}
	var dummyBar = rand.Intn(100)
	var dummyResult = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// expect
	m.Mock(foo).Expects(dummyBar).Returns(dummyResult).Times(5)

	// SUT + act
	var result1 = foo(dummyBar)
	var result2 = foo(dummyBar)
	var result3 = foo(dummyBar)
	var result4 = foo(dummyBar)
	var result5 = foo(dummyBar)

	// assert
	assertEquals(t, dummyResult, result1, "foo call result 1 different")
	assertEquals(t, dummyResult, result2, "foo call result 2 different")
	assertEquals(t, dummyResult, result3, "foo call result 3 different")
	assertEquals(t, dummyResult, result4, "foo call result 4 different")
	assertEquals(t, dummyResult, result5, "foo call result 5 different")
}

func TestMocker_ShouldMockFunctionReturningInterfaceType(t *testing.T) {
	// arrange
	type i interface {
		do()
	}
	var foo = func(bar int) i {
		return nil
	}
	var dummyBar = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// expect
	m.Mock(foo).Expects(dummyBar).Returns(nil).Once()

	// SUT + act
	var result = foo(dummyBar)

	// assert
	assertEquals(t, nil, result, "foo call result different")
}

func TestMocker_ShouldStubPublicFunction(t *testing.T) {
	// arrange
	var dummyMessage = "some random message"
	var dummyResult = errors.New("some other message")

	// mock
	var m = NewMocker(t)

	// expect
	m.Stub(errors.New).Returns(dummyResult).Once()

	// SUT + act
	var result = errors.New(dummyMessage)

	// assert
	assertEquals(t, dummyResult, result, "errors.New call result different")
}

func TestMocker_ShouldMockPublicFunction(t *testing.T) {
	// arrange
	var dummyMessage = "some random message"
	var dummyResult = errors.New("some other message")

	// mock
	var m = NewMocker(t)

	// expect
	m.Mock(errors.New).Expects(dummyMessage).Returns(dummyResult).Once()

	// SUT + act
	var result = errors.New(dummyMessage)

	// assert
	assertEquals(t, dummyResult, result, "errors.New call result different")
}

func TestMocker_ShouldStubVariadicFunction(t *testing.T) {
	// arrange
	var foo = func(int, ...int) int {
		return 0
	}
	var dummyBar = rand.Intn(100)
	var dummyBaz = rand.Intn(100)
	var dummyBam = rand.Intn(100)
	var dummyBat = rand.Intn(100)
	var dummyResult = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// expect
	m.Stub(foo).Returns(dummyResult).Once()

	// SUT + act
	var result = foo(dummyBar, dummyBaz, dummyBam, dummyBat)

	// assert
	assertEquals(t, dummyResult, result, "foo call result different")
}

func TestMocker_ShouldMockVariadicFunction(t *testing.T) {
	// arrange
	var foo = func(int, ...int) int {
		return 0
	}
	var dummyBar = rand.Intn(100)
	var dummyBaz = rand.Intn(100)
	var dummyBam = rand.Intn(100)
	var dummyBat = rand.Intn(100)
	var dummyResult = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// expect
	m.Mock(foo).Expects(dummyBar, dummyBaz, dummyBam, dummyBat).Returns(dummyResult).Once()

	// SUT + act
	var result = foo(dummyBar, dummyBaz, dummyBam, dummyBat)

	// assert
	assertEquals(t, dummyResult, result, "foo call result different")
}

func TestMocker_ShouldStubFunctionWithSideEffect(t *testing.T) {
	// arrange
	var foo = func(int) int {
		return 0
	}
	var dummyValue = rand.Intn(100)
	var dummyResult = rand.Intn(100)
	var dummySideEffect = false

	// mock
	var m = NewMocker(t)

	// expect
	m.Stub(foo).Returns(dummyResult).SideEffect(func(index int, params ...interface{}) {
		dummySideEffect = true
		assertEquals(t, 1, index, "foo call side effect index different")
		assertEquals(t, 1, len(params), "foo call side effect params count different")
		assertEquals(t, dummyValue, params[0], "foo call side effect param 1 different")
	}).Once()

	// SUT + act
	var result = foo(dummyValue)

	// assert
	assertEquals(t, dummyResult, result, "foo call result different")
	assertEquals(t, true, dummySideEffect, "foo call side effect different")
}

func TestMocker_ShouldMockFunctionWithSideEffect(t *testing.T) {
	// arrange
	var foo = func(int) int {
		return 0
	}
	var dummyValue = rand.Intn(100)
	var dummyResult = rand.Intn(100)
	var dummySideEffect = false

	// mock
	var m = NewMocker(t)

	// expect
	m.Mock(foo).Expects(dummyValue).Returns(dummyResult).SideEffect(func(index int, params ...interface{}) {
		dummySideEffect = true
		assertEquals(t, 1, index, "foo call side effect index different")
		assertEquals(t, 1, len(params), "foo call side effect params count different")
		assertEquals(t, dummyValue, params[0], "foo call side effect param 1 different")
	}).Once()

	// SUT + act
	var result = foo(dummyValue)

	// assert
	assertEquals(t, dummyResult, result, "foo call result different")
	assertEquals(t, true, dummySideEffect, "foo call side effect different")
}

func TestMocker_ShouldStubFunctionDifferently(t *testing.T) {
	// arrange
	var foo = func(bar int) int {
		return 0
	}
	var dummyBar1 = rand.Intn(100)
	var dummyBar2 = rand.Intn(100)
	var dummyResult1 = rand.Intn(100)
	var dummyResult2 = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// expect
	m.Stub(foo).Returns(dummyResult1).Once()
	m.Stub(foo).Returns(dummyResult2).Once()

	// SUT + act
	var result1 = foo(dummyBar1)
	var result2 = foo(dummyBar2)

	// assert
	assertEquals(t, dummyResult1, result1, "foo call result 1 different")
	assertEquals(t, dummyResult2, result2, "foo call result 2 different")
}

func TestMocker_ShouldMockFunctionDifferently(t *testing.T) {
	// arrange
	var foo = func(bar int) int {
		return 0
	}
	var dummyBar1 = rand.Intn(100)
	var dummyBar2 = rand.Intn(100)
	var dummyResult1 = rand.Intn(100)
	var dummyResult2 = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// expect
	m.Mock(foo).Expects(dummyBar1).Returns(dummyResult1).Once()
	m.Mock(foo).Expects(dummyBar2).Returns(dummyResult2).Once()

	// SUT + act
	var result1 = foo(dummyBar1)
	var result2 = foo(dummyBar2)

	// assert
	assertEquals(t, dummyResult1, result1, "foo call result 1 different")
	assertEquals(t, dummyResult2, result2, "foo call result 2 different")
}

func TestMocker_ShouldMockFunctionWithParameterStyle(t *testing.T) {
	// arrange
	var foo = func(int, int, int) int {
		return 0
	}
	var dummyBar = rand.Intn(100)
	var dummyBaz = rand.Intn(100)
	var dummyBam = rand.Intn(100)
	var dummyResult = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// expect
	m.Mock(foo).Expects(dummyBar, Anything(), Matches(func(value interface{}) bool { return value == dummyBam })).Returns(dummyResult).Once()

	// SUT + act
	var result = foo(dummyBar, dummyBaz, dummyBam)

	// assert
	assertEquals(t, dummyResult, result, "foo call result different")
}

type testObject struct {
}

func (o *testObject) Foo(bar int) int {
	return 0
}

func TestMocker_ShouldStubStructMethod(t *testing.T) {
	// arrange
	var dummyBar = rand.Intn(100)
	var dummyResult = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// expect
	m.Stub((*testObject).Foo).Returns(dummyResult).Once()

	// SUT + act
	var sut = &testObject{}

	// act
	var result = sut.Foo(dummyBar)

	// assert
	assertEquals(t, dummyResult, result, "testObject.Foo call result different")
}

func TestMocker_ShouldMockStructMethod(t *testing.T) {
	// arrange
	var dummyBar = rand.Intn(100)
	var dummyResult = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// SUT
	var sut = &testObject{}

	// expect
	m.Mock((*testObject).Foo).Expects(sut, dummyBar).Returns(dummyResult).Once()

	// act
	var result = sut.Foo(dummyBar)

	// assert
	assertEquals(t, dummyResult, result, "testObject.Foo call result different")
}

func TestMocker_ShouldStubStructMethodWithSideEffects(t *testing.T) {
	// arrange
	var dummyBar = rand.Intn(100)
	var dummyResult = rand.Intn(100)
	var dummySideEffect = false

	// mock
	var m = NewMocker(t)

	// SUT + act
	var sut = &testObject{}

	// expect
	m.Stub((*testObject).Foo).Returns(dummyResult).SideEffect(func(index int, params ...interface{}) {
		dummySideEffect = true
		assertEquals(t, 1, index, "foo call side effect index different")
		assertEquals(t, 2, len(params), "foo call side effect params count different")
		assertEquals(t, sut, params[0], "foo call side effect param 1 different")
		assertEquals(t, dummyBar, params[1], "foo call side effect param 2 different")
	}).Once()

	// act
	var result = sut.Foo(dummyBar)

	// assert
	assertEquals(t, dummyResult, result, "testObject.Foo call result different")
	assertEquals(t, true, dummySideEffect, "foo call side effect different")
}

func TestMocker_ShouldMockStructMethodWithSideEffect(t *testing.T) {
	// arrange
	var dummyBar = rand.Intn(100)
	var dummyResult = rand.Intn(100)
	var dummySideEffect = false

	// mock
	var m = NewMocker(t)

	// SUT
	var sut = &testObject{}

	// expect
	m.Mock((*testObject).Foo).Expects(sut, dummyBar).Returns(dummyResult).SideEffect(func(index int, params ...interface{}) {
		dummySideEffect = true
		assertEquals(t, 1, index, "foo call side effect index different")
		assertEquals(t, 2, len(params), "foo call side effect params count different")
		assertEquals(t, sut, params[0], "foo call side effect param 1 different")
		assertEquals(t, dummyBar, params[1], "foo call side effect param 2 different")
	}).Once()

	// act
	var result = sut.Foo(dummyBar)

	// assert
	assertEquals(t, dummyResult, result, "testObject.Foo call result different")
	assertEquals(t, true, dummySideEffect, "foo call side effect different")
}

func TestMocker_ShouldStubInterfaceMethod(t *testing.T) {
	// arrange
	type TestInterface interface {
		Foo(int) int
	}
	var foo = func(i TestInterface, bar int) int {
		return i.Foo(bar)
	}
	type testInterface struct {
		TestInterface
	}
	var dummyTestObject = &testInterface{}
	var dummyBar = rand.Intn(100)
	var dummyResult = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// expect
	m.Stub((*testInterface).Foo).Returns(dummyResult).Once()

	// SUT + act
	var result = foo(dummyTestObject, dummyBar)

	// assert
	assertEquals(t, dummyResult, result, "foo call result different")
}

func TestMocker_ShouldMockInterfaceMethod(t *testing.T) {
	// arrange
	type TestInterface interface {
		Foo(int) int
	}
	var foo = func(i TestInterface, bar int) int {
		return i.Foo(bar)
	}
	type testInterface struct {
		TestInterface
	}
	var dummyTestObject = &testInterface{}
	var dummyBar = rand.Intn(100)
	var dummyResult = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// expect
	m.Mock((*testInterface).Foo).Expects(dummyTestObject, dummyBar).Returns(dummyResult).Once()

	// SUT + act
	var result = foo(dummyTestObject, dummyBar)

	// assert
	assertEquals(t, dummyResult, result, "foo call result different")
}

func TestMocker_ShouldStubInterfaceMethodWithSideEffect(t *testing.T) {
	// arrange
	type TestInterface interface {
		Foo(int) int
	}
	var foo = func(i TestInterface, bar int) int {
		return i.Foo(bar)
	}
	type testInterface struct {
		TestInterface
	}
	var dummyTestObject = &testInterface{}
	var dummyBar = rand.Intn(100)
	var dummyResult = rand.Intn(100)
	var dummySideEffect = false

	// mock
	var m = NewMocker(t)

	// expect
	m.Stub((*testInterface).Foo).Returns(dummyResult).SideEffect(func(index int, params ...interface{}) {
		dummySideEffect = true
		assertEquals(t, 1, index, "foo call side effect index different")
		assertEquals(t, 2, len(params), "foo call side effect params count different")
		assertEquals(t, dummyTestObject, params[0], "foo call side effect param 1 different")
		assertEquals(t, dummyBar, params[1], "foo call side effect param 2 different")
	}).Once()

	// SUT + act
	var result = foo(dummyTestObject, dummyBar)

	// assert
	assertEquals(t, dummyResult, result, "foo call result different")
	assertEquals(t, true, dummySideEffect, "foo call side effect different")
}

func TestMocker_ShouldMockInterfaceMethodWithSideEffect(t *testing.T) {
	// arrange
	type TestInterface interface {
		Foo(int) int
	}
	var foo = func(i TestInterface, bar int) int {
		return i.Foo(bar)
	}
	type testInterface struct {
		TestInterface
	}
	var dummyTestObject = &testInterface{}
	var dummyBar = rand.Intn(100)
	var dummyResult = rand.Intn(100)
	var dummySideEffect = false

	// mock
	var m = NewMocker(t)

	// expect
	m.Mock((*testInterface).Foo).Expects(dummyTestObject, dummyBar).Returns(dummyResult).SideEffect(func(index int, params ...interface{}) {
		dummySideEffect = true
		assertEquals(t, 1, index, "foo call side effect index different")
		assertEquals(t, 2, len(params), "foo call side effect params count different")
		assertEquals(t, dummyTestObject, params[0], "foo call side effect param 1 different")
		assertEquals(t, dummyBar, params[1], "foo call side effect param 2 different")
	}).Once()

	// SUT + act
	var result = foo(dummyTestObject, dummyBar)

	// assert
	assertEquals(t, dummyResult, result, "foo call result different")
	assertEquals(t, true, dummySideEffect, "foo call side effect different")
}

type tester struct {
	testing.TB
	t      *testing.T
	errorf func(string, ...interface{})
	fatalf func(string, ...interface{})
}

func (t *tester) Errorf(format string, args ...interface{}) {
	t.errorf(format, args...)
}

func (t *tester) Fatalf(format string, args ...interface{}) {
	t.fatalf(format, args...)
}

func (t *tester) Cleanup(f func()) {
	t.t.Cleanup(f)
}

func (t *tester) Helper() {
	t.t.Helper()
}

func TestMocker_ShouldReportTestFailureWhenMockFunctionIsNotCalledButExpected(t *testing.T) {
	// arrange
	var foo = func(bar int) int {
		return 0
	}
	var tester = &tester{t: t}

	// mock
	var m = NewMocker(tester)

	// expect
	tester.errorf = func(format string, args ...interface{}) {
		assertEquals(t, "[%v] Unepxected number of calls: expect %v, actual %v", format, "tester.Errorf called with different message")
		assertEquals(t, 3, len(args), "tester.Errorf called with different number of args")
		assertEquals(t, 1, args[1], "tester.Errorf called with different argument 2")
		assertEquals(t, 0, args[2], "tester.Errorf called with different argument 3")
	}
	m.Mock(foo).Expects().Returns().Once()
}

func TestMocker_ShouldReportTestFailureWhenMockFunctionIsCalledButNotExpected(t *testing.T) {
	defer func() {
		recover()
	}()

	// arrange
	var foo = func() {}
	var tester = &tester{t: t}

	// mock
	var m = NewMocker(tester)

	// expect
	tester.fatalf = func(format string, args ...interface{}) {
		assertEquals(t, "[%v] Unepxected number of calls: expect %v, actual %v", format, "tester.Fatalf called with different message")
		assertEquals(t, 3, len(args), "tester.Fatalf called with different number of args")
		assertEquals(t, 0, args[1], "tester.Fatalf called with different argument 2")
		assertEquals(t, 1, args[2], "tester.Fatalf called with different argument 3")
	}

	// SUT
	m.Mock(foo).NotCalled()

	// act
	foo()
}

func TestMocker_ShouldReportTestFailureWhenMockFunctionPanicsWithErrorInExecution(t *testing.T) {
	defer func() {
		recover()
	}()

	// arrange
	var foo = func() {}
	var tester = &tester{t: t}

	// mock
	var m = NewMocker(tester)

	// expect
	tester.errorf = func(format string, args ...interface{}) {
		assertEquals(t, "[%v] Mocker panicing recovered: %v", format, "tester.Errorf called with different message")
		assertEquals(t, 2, len(args), "tester.Errorf called with different number of args")
		assertEquals(t, "paniced", args[1], "tester.Errorf called with different argument 2")
	}
	m.Mock(foo).Expects().Returns().SideEffect(func(index int, params ...interface{}) {
		panic(errors.New("paniced"))
	}).Once()

	// SUT + act
	foo()
}

func TestMocker_ShouldReportTestFailureWhenMockFunctionPanicsWithMessageInExecution(t *testing.T) {
	defer func() {
		recover()
	}()

	// arrange
	var foo = func() {}
	var tester = &tester{t: t}

	// mock
	var m = NewMocker(tester)

	// expect
	tester.errorf = func(format string, args ...interface{}) {
		assertEquals(t, "[%v] Mocker panicing recovered: %v", format, "tester.Errorf called with different message")
		assertEquals(t, 2, len(args), "tester.Errorf called with different number of args")
		assertEquals(t, "paniced", args[1], "tester.Errorf called with different argument 2")
	}
	m.Mock(foo).Expects().Returns().SideEffect(func(index int, params ...interface{}) {
		panic("paniced")
	}).Once()

	// SUT + act
	foo()
}

func TestMocker_ShouldReportTestFailureWhenMockFunctionParameterValueNotEqual(t *testing.T) {
	// arrange
	var foo = func(_ int) {}
	var tester = &tester{t: t}
	var dummyBar = rand.Intn(100)

	// mock
	var m = NewMocker(tester)

	// expect
	tester.errorf = func(format string, args ...interface{}) {
		assertEquals(t, "[%v] Parameter mismatch at call #%v parameter #%v: expect %v, actual %v", format, "tester.Errorf called with different message")
		assertEquals(t, 5, len(args), "tester.Errorf called with different number of args")
		assertEquals(t, 1, args[1], "tester.Errorf called with different argument 2")
		assertEquals(t, 1, args[2], "tester.Errorf called with different argument 3")
		assertEquals(t, dummyBar+1, args[3], "tester.Errorf called with different argument 4")
		assertEquals(t, dummyBar, args[4], "tester.Errorf called with different argument 5")
	}
	m.Mock(foo).Expects(dummyBar + 1).Returns().Once()

	// SUT + act
	foo(dummyBar)
}

func TestMocker_ShouldReportTestFailureWhenMockFunctionParameterValueMismatch(t *testing.T) {
	// arrange
	var foo = func(_ int) {}
	var tester = &tester{t: t}
	var dummyBar = rand.Intn(100)

	// mock
	var m = NewMocker(tester)

	// expect
	tester.errorf = func(format string, args ...interface{}) {
		assertEquals(t, "[%v] Parameter mismatch at call #%v parameter #%v: matchFunc failed on actual %v", format, "tester.Errorf called with different message")
		assertEquals(t, 4, len(args), "tester.Errorf called with different number of args")
		assertEquals(t, 1, args[1], "tester.Errorf called with different argument 2")
		assertEquals(t, 1, args[2], "tester.Errorf called with different argument 3")
		assertEquals(t, dummyBar, args[3], "tester.Errorf called with different argument 4")
	}
	m.Mock(foo).Expects(Matches(func(value interface{}) bool { return false })).Returns().Once()

	// SUT + act
	foo(dummyBar)
}

func TestMocker_ShouldReportTestFailureWhenMockFunctionNormalParameterCountMismatch(t *testing.T) {
	// arrange
	var foo = func(_ int) {}
	var tester = &tester{t: t}
	var dummyBar = rand.Intn(100)

	// mock
	var m = NewMocker(tester)

	// expect
	tester.errorf = func(format string, args ...interface{}) {
		assertEquals(t, "[%v] Invalid number of parameters at call #%v: expect %v, actual %v", format, "tester.Errorf called with different message")
		assertEquals(t, 4, len(args), "tester.Errorf called with different number of args")
		assertEquals(t, 1, args[1], "tester.Errorf called with different argument 2")
		assertEquals(t, 0, args[2], "tester.Errorf called with different argument 3")
		assertEquals(t, 1, args[3], "tester.Errorf called with different argument 4")
	}
	m.Mock(foo).Expects().Returns().Once()

	// SUT + act
	foo(dummyBar)
}

func TestMocker_ShouldReportTestFailureWhenMockFunctionVariadicParameterCountMismatch(t *testing.T) {
	// arrange
	var foo = func(_ ...int) {}
	var tester = &tester{t: t}
	var dummyBar = rand.Intn(100)

	// mock
	var m = NewMocker(tester)

	// expect
	tester.errorf = func(format string, args ...interface{}) {
		assertEquals(t, "[%v] Invalid number of variadic parameters at call #%v: expect %v, actual %v", format, "tester.Errorf called with different message")
		assertEquals(t, 4, len(args), "tester.Errorf called with different number of args")
		assertEquals(t, 1, args[1], "tester.Errorf called with different argument 2")
		assertEquals(t, 0, args[2], "tester.Errorf called with different argument 3")
		assertEquals(t, 1, args[3], "tester.Errorf called with different argument 4")
	}
	m.Mock(foo).Expects().Returns().Once()

	// SUT + act
	foo(dummyBar)
}

func TestMocker_ShouldReportTestFailureWhenMockFunctionReturnCountMismatch(t *testing.T) {
	defer func() {
		recover()
	}()

	// arrange
	var foo = func() int { return 0 }
	var tester = &tester{t: t}

	// mock
	var m = NewMocker(tester)

	// expect
	tester.fatalf = func(format string, args ...interface{}) {
		assertEquals(t, "[%v] Invalid number of returns at call #%v: expect %v, actual %v", format, "tester.Fatalf called with different message")
		assertEquals(t, 4, len(args), "tester.Fatalf called with different number of args")
		assertEquals(t, 1, args[1], "tester.Fatalf called with different argument 2")
		assertEquals(t, 1, args[2], "tester.Fatalf called with different argument 3")
		assertEquals(t, 0, args[3], "tester.Fatalf called with different argument 4")
	}

	// SUT
	m.Mock(foo).Expects().Returns().Once()

	// act
	foo()
}

func TestMocker_ShouldHandleEntryNotFoundScenarioWhenMakeFunc(t *testing.T) {
	defer func() {
		recover()
	}()

	// arrange
	var foo = func() int { return 0 }
	var dummyName = "some name"
	var dummyFuncPtr = uintptr(rand.Intn(100))
	var dummyFuncType = reflect.TypeOf(foo)
	var tester = &tester{t: t}

	// mock
	var m = NewMocker(tester).(*mocker)

	// expect
	tester.fatalf = func(format string, args ...interface{}) {
		assertEquals(t, "The underlying function or method %v was never setup", format, "tester.Fatalf called with different message")
		assertEquals(t, 1, len(args), "tester.Fatalf called with different number of args")
		assertEquals(t, dummyName, args[0], "tester.Fatalf called with different argument 1")
	}

	// SUT
	var sut = m.makeFunc(dummyName, dummyFuncPtr, dummyFuncType)

	// act
	sut.Call([]reflect.Value{})
}

func TestMocker_ShouldHandleEntryCountMismatchScenarioWhenMakeFuncMock(t *testing.T) {
	defer func() {
		recover()
	}()

	// arrange
	var foo = func() int { return 0 }
	var dummyName = "some name"
	var dummyFuncPtr = uintptr(rand.Intn(100))
	var dummyFuncType = reflect.TypeOf(foo)
	var tester = &tester{t: t}

	// mock
	var m = NewMocker(tester).(*mocker)

	// stub
	m.entries[dummyFuncPtr] = &funcEntry{
		name: dummyName,
	}

	// expect
	tester.fatalf = func(format string, args ...interface{}) {
		assertEquals(t, "[%v] Unepxected number of calls: expect %v, actual %v", format, "tester.Fatalf called with different message")
		assertEquals(t, 3, len(args), "tester.Fatalf called with different number of args")
		assertEquals(t, dummyName, args[0], "tester.Fatalf called with different argument 1")
		assertEquals(t, 0, args[1], "tester.Fatalf called with different argument 2")
		assertEquals(t, 1, args[2], "tester.Fatalf called with different argument 3")
	}

	// SUT
	var sut = m.makeFunc(dummyName, dummyFuncPtr, dummyFuncType)

	// act
	sut.Call([]reflect.Value{})
}

func TestMocker_ShouldHandleEntryCountMismatchScenarioWhenMakeFuncStub(t *testing.T) {
	// arrange
	var foo = func() {}
	var dummyName = "some name"
	var dummyFuncPtr = uintptr(rand.Intn(100))
	var dummyFuncType = reflect.TypeOf(foo)
	var tester = &tester{t: t}

	// mock
	var m = &mocker{
		tester: tester,
		entries: map[uintptr]*funcEntry{
			dummyFuncPtr: {
				name:   dummyName,
				stub:   true,
				expect: 0,
				actual: 0,
				mocks: []*mockEntry{
					{},
				},
			},
		},
		locker: &sync.Mutex{},
	}

	// expect
	tester.errorf = func(format string, args ...interface{}) {
		assertEquals(t, "[%v] Unepxected number of calls: expect %v, actual %v", format, "tester.Errorf called with different message")
		assertEquals(t, 3, len(args), "tester.Errorf called with different number of args")
		assertEquals(t, dummyName, args[0], "tester.Errorf called with different argument 1")
		assertEquals(t, 0, args[1], "tester.Errorf called with different argument 2")
		assertEquals(t, 1, args[2], "tester.Errorf called with different argument 3")
	}

	// SUT
	var sut = m.makeFunc(dummyName, dummyFuncPtr, dummyFuncType)

	// act
	sut.Call([]reflect.Value{})
}

func TestMocker_ShouldReportErrorIfAFormerSetupWasIncompleteWhenCallingANewSetup1(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyStub = false
	var dummyFuncPtr = uintptr(rand.Intn(100))
	var tester = &tester{t: t}

	// expect
	tester.fatalf = func(format string, args ...interface{}) {
		assertEquals(t, "A former setup for function or method [%v] was incomplete. Did you miss calling the Once/Twice/Times method in the end?", format, "tester.Fatalf called with different message")
		assertEquals(t, 1, len(args), "tester.Fatalf called with different number of args")
		assertEquals(t, dummyName, args[0], "tester.Fatalf called with different argument 1")
	}

	// SUT
	var m = &mocker{
		tester:  tester,
		current: &funcEntry{},
	}

	// act
	m.setup(dummyName, dummyStub, dummyFuncPtr)
}

func TestMocker_ShouldReportErrorIfAFormerSetupWasIncompleteWhenCallingANewSetup2(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyStub = false
	var dummyFuncPtr = uintptr(rand.Intn(100))
	var tester = &tester{t: t}

	// expect
	tester.fatalf = func(format string, args ...interface{}) {
		assertEquals(t, "A former setup for function or method [%v] was incomplete. Did you miss calling the Once/Twice/Times method in the end?", format, "tester.Fatalf called with different message")
		assertEquals(t, 1, len(args), "tester.Fatalf called with different number of args")
		assertEquals(t, dummyName, args[0], "tester.Fatalf called with different argument 1")
	}

	// SUT
	var m = &mocker{
		tester: tester,
		temp:   &mockEntry{},
	}

	// act
	m.setup(dummyName, dummyStub, dummyFuncPtr)
}

func TestMocker_ShouldReportErrorIfAFormerSetupWasMockButCurrentSetupIsStub(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyStub = false
	var dummyFuncPtr = uintptr(rand.Intn(100))
	var tester = &tester{t: t}

	// expect
	tester.fatalf = func(format string, args ...interface{}) {
		assertEquals(t, "A former setup for function or method [%v] was a Stub but current setup is a Mock. We do not support mixing Stub and Mock for the same function or method at the moment.", format, "tester.Fatalf called with different message")
		assertEquals(t, 1, len(args), "tester.Fatalf called with different number of args")
		assertEquals(t, dummyName, args[0], "tester.Fatalf called with different argument 1")
	}

	// SUT
	var m = &mocker{
		tester: tester,
		entries: map[uintptr]*funcEntry{
			dummyFuncPtr: {stub: true},
		},
	}

	// act
	m.setup(dummyName, dummyStub, dummyFuncPtr)
}

func TestMocker_ShouldReportErrorIfAFormerSetupWasStubButCurrentSetupIsMock(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyStub = true
	var dummyFuncPtr = uintptr(rand.Intn(100))
	var tester = &tester{t: t}

	// expect
	tester.fatalf = func(format string, args ...interface{}) {
		assertEquals(t, "A former setup for function or method [%v] was a Mock but current setup is a Stub. We do not support mixing Stub and Mock for the same function or method at the moment.", format, "tester.Fatalf called with different message")
		assertEquals(t, 1, len(args), "tester.Fatalf called with different number of args")
		assertEquals(t, dummyName, args[0], "tester.Fatalf called with different argument 1")
	}

	// SUT
	var m = &mocker{
		tester: tester,
		entries: map[uintptr]*funcEntry{
			dummyFuncPtr: {stub: false},
		},
	}

	// act
	m.setup(dummyName, dummyStub, dummyFuncPtr)
}

func TestMocker_ShouldReportErrorIfAFormerSetupToBeNotCalledButMockAgain(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyStub = false
	var dummyFuncPtr = uintptr(rand.Intn(100))
	var tester = &tester{t: t}

	// expect
	tester.fatalf = func(format string, args ...interface{}) {
		assertEquals(t, "A former setup for function or method [%v] was to be not called, therefore no more Mock or Stub can be setup for it now.", format, "tester.Fatalf called with different message")
		assertEquals(t, 1, len(args), "tester.Fatalf called with different number of args")
		assertEquals(t, dummyName, args[0], "tester.Fatalf called with different argument 1")
	}

	// SUT
	var m = &mocker{
		tester: tester,
		entries: map[uintptr]*funcEntry{
			dummyFuncPtr: {
				nocall: true,
				stub:   false,
			},
		},
	}

	// act
	m.setup(dummyName, dummyStub, dummyFuncPtr)
}

func TestMocker_ShouldReportErrorIfAFormerSetupToBeNotCalledButStubAgain(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyStub = true
	var dummyFuncPtr = uintptr(rand.Intn(100))
	var tester = &tester{t: t}

	// expect
	tester.fatalf = func(format string, args ...interface{}) {
		assertEquals(t, "A former setup for function or method [%v] was to be not called, therefore no more Mock or Stub can be setup for it now.", format, "tester.Fatalf called with different message")
		assertEquals(t, 1, len(args), "tester.Fatalf called with different number of args")
		assertEquals(t, dummyName, args[0], "tester.Fatalf called with different argument 1")
	}

	// SUT
	var m = &mocker{
		tester: tester,
		entries: map[uintptr]*funcEntry{
			dummyFuncPtr: {
				nocall: true,
				stub:   true,
			},
		},
	}

	// act
	m.setup(dummyName, dummyStub, dummyFuncPtr)
}

func TestMocker_ShouldReportErrorIfNoFormerSetupWhenCallingExpects(t *testing.T) {
	// arrange
	var tester = &tester{t: t}

	// expect
	tester.fatalf = func(format string, args ...interface{}) {
		assertEquals(t, "Unexpected call to Expects without setting up an anticipated function or method", format, "tester.Fatalf called with different message")
		assertEquals(t, 0, len(args), "tester.Fatalf called with different number of args")
	}

	// SUT
	var m = &mocker{
		tester: tester,
	}

	// act
	m.Expects()
}

func TestMocker_ShouldReportErrorIfNoFormerSetupWhenCallingNotCalled(t *testing.T) {
	// arrange
	var tester = &tester{t: t}

	// expect
	tester.fatalf = func(format string, args ...interface{}) {
		assertEquals(t, "Unexpected call to NotCalled without setting up an anticipated function or method", format, "tester.Fatalf called with different message")
		assertEquals(t, 0, len(args), "tester.Fatalf called with different number of args")
	}

	// SUT
	var m = &mocker{
		tester: tester,
	}

	// act
	m.NotCalled()
}

func TestMocker_ShouldReportErrorIfNoFormerSetupWhenCallingReturns(t *testing.T) {
	// arrange
	var tester = &tester{t: t}

	// expect
	tester.fatalf = func(format string, args ...interface{}) {
		assertEquals(t, "Unexpected call to Returns without setting up an anticipated function or method", format, "tester.Fatalf called with different message")
		assertEquals(t, 0, len(args), "tester.Fatalf called with different number of args")
	}

	// SUT
	var m = &mocker{
		tester: tester,
	}

	// act
	m.Returns()
}

func TestMocker_ShouldReportErrorIfNoFormerSetupWhenCallingSideEffect(t *testing.T) {
	// arrange
	var tester = &tester{t: t}

	// expect
	tester.fatalf = func(format string, args ...interface{}) {
		assertEquals(t, "Unexpected call to SideEffect without setting up an anticipated function or method", format, "tester.Fatalf called with different message")
		assertEquals(t, 0, len(args), "tester.Fatalf called with different number of args")
	}

	// SUT
	var m = &mocker{
		tester: tester,
	}

	// act
	m.SideEffect(nil)
}

func TestMocker_ShouldReportErrorIfNoFormerSetupWhenCallingTimes(t *testing.T) {
	// arrange
	var tester = &tester{t: t}

	// expect
	tester.fatalf = func(format string, args ...interface{}) {
		assertEquals(t, "Unexpected call to Times without setting up an anticipated function or method", format, "tester.Fatalf called with different message")
		assertEquals(t, 0, len(args), "tester.Fatalf called with different number of args")
	}

	// SUT
	var m = &mocker{
		tester: tester,
	}

	// act
	m.Times(0)
}

func TestMocker_ShouldReportErrorIfCountIsNegativeWhenCallingTimes(t *testing.T) {
	// arrange
	var tester = &tester{t: t}
	var dummyName = "some name"
	var dummyCount = -1 - rand.Intn(100)

	// expect
	tester.fatalf = func(format string, args ...interface{}) {
		assertEquals(t, "function or method [%v] cannot be mocked for negative [%v] times", format, "tester.Fatalf called with different message")
		assertEquals(t, 2, len(args), "tester.Fatalf called with different number of args")
		assertEquals(t, dummyName, args[0], "tester.Fatalf called with different argument 1")
		assertEquals(t, dummyCount, args[1], "tester.Fatalf called with different argument 2")
	}

	// SUT
	var m = &mocker{
		tester: tester,
		current: &funcEntry{
			name: dummyName,
		},
		temp: &mockEntry{},
	}

	// act
	m.Times(dummyCount)
}

func TestMocker_ShouldReportErrorIfCountIsZeroWhenCallingTimes(t *testing.T) {
	// arrange
	var tester = &tester{t: t}
	var dummyName = "some name"

	// expect
	tester.fatalf = func(format string, args ...interface{}) {
		assertEquals(t, "function or method [%v] cannot be mocked for zero times using Times method. Try using NotCalled method instead.", format, "tester.Fatalf called with different message")
		assertEquals(t, 1, len(args), "tester.Fatalf called with different number of args")
		assertEquals(t, dummyName, args[0], "tester.Fatalf called with different argument 1")
	}

	// SUT
	var m = &mocker{
		tester: tester,
		current: &funcEntry{
			name: dummyName,
		},
		temp: &mockEntry{},
	}

	// act
	m.Times(0)
}
