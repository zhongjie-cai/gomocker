package gomocker

import (
	"errors"
	"math/rand"
	"reflect"
	"testing"
)

func foo1(bar int) int {
	return bar * 2
}

func TestMocker_ShouldMockPrivateFunctionOnce(t *testing.T) {
	// arrange
	var dummyBar = rand.Intn(100)
	var dummyResult = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// expect
	m.ExpectFunc(foo1, 1, func(bar int) int {
		assertEquals(t, dummyBar, bar, "foo1 call parameter different")
		return dummyResult
	})

	// SUT + act
	var result = foo1(dummyBar)

	// assert
	assertEquals(t, dummyResult, result, "foo1 call result different")
}

func TestMocker_ShouldMockPrivateFunctionMultipleTimes(t *testing.T) {
	// arrange
	var dummyBar1 = rand.Intn(100)
	var dummyBar2 = rand.Intn(100)
	var dummyResult1 = rand.Intn(100)
	var dummyResult2 = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// expect
	m.ExpectFunc(foo1, 1, func(bar int) int {
		assertEquals(t, dummyBar1, bar, "foo1 1st call parameter different")
		return dummyResult1
	}).ExpectFunc(foo1, 1, func(bar int) int {
		assertEquals(t, dummyBar2, bar, "foo1 2nd call parameter different")
		return dummyResult2
	})

	// SUT + act
	var result1 = foo1(dummyBar1)
	var result2 = foo1(dummyBar2)

	// assert
	assertEquals(t, dummyResult1, result1, "foo1 1st call result different")
	assertEquals(t, dummyResult2, result2, "foo1 2nd call result different")
}

func TestMocker_ShouldMockPrivateFunctionMultipleTimes_UseCounter(t *testing.T) {
	// arrange
	var dummyBar1 = rand.Intn(100)
	var dummyBar2 = rand.Intn(100)
	var dummyResult1 = rand.Intn(100)
	var dummyResult2 = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// expect
	m.ExpectFunc(foo1, 2, func(bar int) int {
		if m.FuncCalledCount(foo1) == 1 {
			assertEquals(t, dummyBar1, bar, "foo1 1st call parameter different")
			return dummyResult1
		} else if m.FuncCalledCount(foo1) == 2 {
			assertEquals(t, dummyBar2, bar, "foo1 2nd call parameter different")
			return dummyResult2
		}
		t.Errorf("foo1 called more than twice")
		return -1
	})

	// SUT + act
	var result1 = foo1(dummyBar1)
	var result2 = foo1(dummyBar2)

	// assert
	assertEquals(t, dummyResult1, result1, "foo1 1st call result different")
	assertEquals(t, dummyResult2, result2, "foo1 2nd call result different")
	assertEquals(t, 0, m.FuncCalledCount(assertEquals), "not mocked function count calculated wrongly")
}

func Foo1(bar int) int {
	return bar * 2
}

func TestMocker_ShouldMockPublicFunctionOnce(t *testing.T) {
	// arrange
	var dummyBar = rand.Intn(100)
	var dummyResult = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// expect
	m.ExpectFunc(Foo1, 1, func(bar int) int {
		assertEquals(t, dummyBar, bar, "Foo1 call parameter different")
		return dummyResult
	})

	// SUT + act
	var result = Foo1(dummyBar)

	// assert
	assertEquals(t, dummyResult, result, "Foo1 call result different")
}

func TestMocker_ShouldMockPublicFunctionMultipleTimes(t *testing.T) {
	// arrange
	var dummyBar1 = rand.Intn(100)
	var dummyBar2 = rand.Intn(100)
	var dummyResult1 = rand.Intn(100)
	var dummyResult2 = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// expect
	m.ExpectFunc(Foo1, 1, func(bar int) int {
		assertEquals(t, dummyBar1, bar, "Foo1 1st call parameter different")
		return dummyResult1
	}).ExpectFunc(Foo1, 1, func(bar int) int {
		assertEquals(t, dummyBar2, bar, "Foo1 2nd call parameter different")
		return dummyResult2
	})

	// SUT + act
	var result1 = Foo1(dummyBar1)
	var result2 = Foo1(dummyBar2)

	// assert
	assertEquals(t, dummyResult1, result1, "Foo1 1st call result different")
	assertEquals(t, dummyResult2, result2, "Foo1 2nd call result different")
}

func TestMocker_ShouldMockPublicFunctionMultipleTimes_UseCounter(t *testing.T) {
	// arrange
	var dummyBar1 = rand.Intn(100)
	var dummyBar2 = rand.Intn(100)
	var dummyResult1 = rand.Intn(100)
	var dummyResult2 = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// expect
	m.ExpectFunc(Foo1, 2, func(bar int) int {
		if m.FuncCalledCount(Foo1) == 1 {
			assertEquals(t, dummyBar1, bar, "Foo1 1st call parameter different")
			return dummyResult1
		} else if m.FuncCalledCount(Foo1) == 2 {
			assertEquals(t, dummyBar2, bar, "Foo1 2nd call parameter different")
			return dummyResult2
		}
		t.Errorf("foo1 called more than twice")
		return -1
	})

	// SUT + act
	var result1 = Foo1(dummyBar1)
	var result2 = Foo1(dummyBar2)

	// assert
	assertEquals(t, dummyResult1, result1, "Foo1 1st call result different")
	assertEquals(t, dummyResult2, result2, "Foo1 2nd call result different")
	assertEquals(t, 0, m.FuncCalledCount(assertEquals), "not mocked function count calculated wrongly")
}

type foo2 struct {
	self int
}

func (f *foo2) bar1(bar int) int {
	return f.self * bar
}

func TestMocker_ShouldMockPrivateFunctionOnStructPointerReceiver(t *testing.T) {
	// arrange
	var dummySelf = rand.Intn(100)
	var dummyBar = rand.Intn(100)
	var dummyResult = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// SUT
	var sut = &foo2{
		self: dummySelf,
	}

	// expect
	m.ExpectFunc((*foo2).bar1, 1, func(obj *foo2, bar int) int {
		assertEquals(t, sut, obj, "(*foo2).bar1 call instance different")
		assertEquals(t, dummyBar, bar, "(*foo2).bar1 call parameter different")
		return dummyResult
	})

	// act
	var result = sut.bar1(dummyBar)

	// assert
	assertEquals(t, dummyResult, result, "(*foo2).bar1 call result different")
}

func TestMocker_ShouldMockPrivateFunctionOnStructPointerReceiverMultipleTimes(t *testing.T) {
	// arrange
	var dummySelf = rand.Intn(100)
	var dummyBar1 = rand.Intn(100)
	var dummyBar2 = rand.Intn(100)
	var dummyResult1 = rand.Intn(100)
	var dummyResult2 = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// SUT
	var sut = &foo2{
		self: dummySelf,
	}

	// expect
	m.ExpectFunc((*foo2).bar1, 1, func(obj *foo2, bar int) int {
		assertEquals(t, sut, obj, "(*foo2).bar1 call instance different")
		assertEquals(t, dummyBar1, bar, "(*foo2).bar1 call parameter different")
		return dummyResult1
	}).ExpectFunc((*foo2).bar1, 1, func(obj *foo2, bar int) int {
		assertEquals(t, sut, obj, "(*foo2).bar1 call instance different")
		assertEquals(t, dummyBar2, bar, "(*foo2).bar1 call parameter different")
		return dummyResult2
	})

	// act
	var result1 = sut.bar1(dummyBar1)
	var result2 = sut.bar1(dummyBar2)

	// assert
	assertEquals(t, dummyResult1, result1, "(*foo2).bar1 1st call result different")
	assertEquals(t, dummyResult2, result2, "(*foo2).bar1 2nd call result different")
}

func TestMocker_ShouldMockPrivateFunctionOnStructPointerReceiverMultipleTimes_UseCounter(t *testing.T) {
	// arrange
	var dummySelf = rand.Intn(100)
	var dummyBar1 = rand.Intn(100)
	var dummyBar2 = rand.Intn(100)
	var dummyResult1 = rand.Intn(100)
	var dummyResult2 = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// SUT
	var sut = &foo2{
		self: dummySelf,
	}

	// expect
	m.ExpectFunc((*foo2).bar1, 2, func(obj *foo2, bar int) int {
		if m.FuncCalledCount((*foo2).bar1) == 1 {
			assertEquals(t, sut, obj, "(*foo2).bar1 call instance different")
			assertEquals(t, dummyBar1, bar, "(*foo2).bar1 call parameter different")
			return dummyResult1
		} else if m.FuncCalledCount((*foo2).bar1) == 2 {
			assertEquals(t, sut, obj, "(*foo2).bar1 call instance different")
			assertEquals(t, dummyBar2, bar, "(*foo2).bar1 call parameter different")
			return dummyResult2
		}
		t.Errorf("(*foo2).bar1 called more than twice")
		return -1
	})

	// act
	var result1 = sut.bar1(dummyBar1)
	var result2 = sut.bar1(dummyBar2)

	// assert
	assertEquals(t, dummyResult1, result1, "(*foo2).bar1 1st call result different")
	assertEquals(t, dummyResult2, result2, "(*foo2).bar1 2nd call result different")
	assertEquals(t, 0, m.FuncCalledCount(assertEquals), "not mocked function count calculated wrongly")
}

func TestMocker_ShouldMockPrivateMethodOnStructPointerReceiver(t *testing.T) {
	// arrange
	var dummySelf = rand.Intn(100)
	var dummyBar = rand.Intn(100)
	var dummyResult = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// SUT
	var sut = &foo2{
		self: dummySelf,
	}

	// expect
	m.ExpectMethod(&foo2{}, "bar1", 1, func(obj *foo2, bar int) int {
		assertEquals(t, sut, obj, "(*foo2).bar1 call instance different")
		assertEquals(t, dummyBar, bar, "(*foo2).bar1 call parameter different")
		return dummyResult
	})

	// act
	var result = sut.bar1(dummyBar)

	// assert
	assertEquals(t, dummyResult, result, "(*foo2).bar1 call result different")
}

func TestMocker_ShouldMockPrivateMethodOnStructPointerReceiverMultipleTimes(t *testing.T) {
	// arrange
	var dummySelf = rand.Intn(100)
	var dummyBar1 = rand.Intn(100)
	var dummyBar2 = rand.Intn(100)
	var dummyResult1 = rand.Intn(100)
	var dummyResult2 = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// SUT
	var sut = &foo2{
		self: dummySelf,
	}

	// expect
	m.ExpectMethod(&foo2{}, "bar1", 1, func(obj *foo2, bar int) int {
		assertEquals(t, sut, obj, "(*foo2).bar1 call instance different")
		assertEquals(t, dummyBar1, bar, "(*foo2).bar1 call parameter different")
		return dummyResult1
	}).ExpectMethod(&foo2{}, "bar1", 1, func(obj *foo2, bar int) int {
		assertEquals(t, sut, obj, "(*foo2).bar1 call instance different")
		assertEquals(t, dummyBar2, bar, "(*foo2).bar1 call parameter different")
		return dummyResult2
	})

	// act
	var result1 = sut.bar1(dummyBar1)
	var result2 = sut.bar1(dummyBar2)

	// assert
	assertEquals(t, dummyResult1, result1, "(*foo2).bar1 1st call result different")
	assertEquals(t, dummyResult2, result2, "(*foo2).bar1 2nd call result different")
}

func TestMocker_ShouldMockPrivateMethodOnStructPointerReceiverMultipleTimes_UseCounter(t *testing.T) {
	// arrange
	var dummySelf = rand.Intn(100)
	var dummyBar1 = rand.Intn(100)
	var dummyBar2 = rand.Intn(100)
	var dummyResult1 = rand.Intn(100)
	var dummyResult2 = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// SUT
	var sut = &foo2{
		self: dummySelf,
	}

	// expect
	m.ExpectMethod(&foo2{}, "bar1", 2, func(obj *foo2, bar int) int {
		if m.MethodCalledCount(&foo2{}, "bar1") == 1 {
			assertEquals(t, sut, obj, "(*foo2).bar1 call instance different")
			assertEquals(t, dummyBar1, bar, "(*foo2).bar1 call parameter different")
			return dummyResult1
		} else if m.MethodCalledCount(&foo2{}, "bar1") == 2 {
			assertEquals(t, sut, obj, "(*foo2).bar1 call instance different")
			assertEquals(t, dummyBar2, bar, "(*foo2).bar1 call parameter different")
			return dummyResult2
		}
		t.Errorf("(*foo2).bar1 called more than twice")
		return -1
	})

	// act
	var result1 = sut.bar1(dummyBar1)
	var result2 = sut.bar1(dummyBar2)

	// assert
	assertEquals(t, dummyResult1, result1, "(*foo2).bar1 1st call result different")
	assertEquals(t, dummyResult2, result2, "(*foo2).bar1 2nd call result different")
	assertEquals(t, 0, m.MethodCalledCount(&foo2{}, "Bar1"), "not mocked function count calculated wrongly")
}

func (f *foo2) Bar1(bar int) int {
	return f.self * bar
}

func TestMocker_ShouldMockPublicFunctionOnStructPointerReceiver(t *testing.T) {
	// arrange
	var dummySelf = rand.Intn(100)
	var dummyBar = rand.Intn(100)
	var dummyResult = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// SUT
	var sut = &foo2{
		self: dummySelf,
	}

	// expect
	m.ExpectFunc((*foo2).Bar1, 1, func(obj *foo2, bar int) int {
		assertEquals(t, sut, obj, "(*foo2).Bar1 call instance different")
		assertEquals(t, dummyBar, bar, "(*foo2).Bar1 call parameter different")
		return dummyResult
	})

	// act
	var result = sut.Bar1(dummyBar)

	// assert
	assertEquals(t, dummyResult, result, "(*foo2).Bar1 call result different")
}

func TestMocker_ShouldMockPublicFunctionOnStructPointerReceiverMultipleTimes(t *testing.T) {
	// arrange
	var dummySelf = rand.Intn(100)
	var dummyBar1 = rand.Intn(100)
	var dummyBar2 = rand.Intn(100)
	var dummyResult1 = rand.Intn(100)
	var dummyResult2 = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// SUT
	var sut = &foo2{
		self: dummySelf,
	}

	// expect
	m.ExpectFunc((*foo2).Bar1, 1, func(obj *foo2, bar int) int {
		assertEquals(t, sut, obj, "(*foo2).Bar1 call instance different")
		assertEquals(t, dummyBar1, bar, "(*foo2).Bar1 call parameter different")
		return dummyResult1
	}).ExpectFunc((*foo2).Bar1, 1, func(obj *foo2, bar int) int {
		assertEquals(t, sut, obj, "(*foo2).Bar1 call instance different")
		assertEquals(t, dummyBar2, bar, "(*foo2).Bar1 call parameter different")
		return dummyResult2
	})

	// act
	var result1 = sut.Bar1(dummyBar1)
	var result2 = sut.Bar1(dummyBar2)

	// assert
	assertEquals(t, dummyResult1, result1, "(*foo2).Bar1 1st call result different")
	assertEquals(t, dummyResult2, result2, "(*foo2).Bar1 2nd call result different")
}

func TestMocker_ShouldMockPublicFunctionOnStructPointerReceiverMultipleTimes_UseCounter(t *testing.T) {
	// arrange
	var dummySelf = rand.Intn(100)
	var dummyBar1 = rand.Intn(100)
	var dummyBar2 = rand.Intn(100)
	var dummyResult1 = rand.Intn(100)
	var dummyResult2 = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// SUT
	var sut = &foo2{
		self: dummySelf,
	}

	// expect
	m.ExpectFunc((*foo2).Bar1, 2, func(obj *foo2, bar int) int {
		if m.FuncCalledCount((*foo2).Bar1) == 1 {
			assertEquals(t, sut, obj, "(*foo2).Bar1 call instance different")
			assertEquals(t, dummyBar1, bar, "(*foo2).Bar1 call parameter different")
			return dummyResult1
		} else if m.FuncCalledCount((*foo2).Bar1) == 2 {
			assertEquals(t, sut, obj, "(*foo2).Bar1 call instance different")
			assertEquals(t, dummyBar2, bar, "(*foo2).Bar1 call parameter different")
			return dummyResult2
		}
		t.Errorf("(*foo2).Bar1 called more than twice")
		return -1
	})

	// act
	var result1 = sut.Bar1(dummyBar1)
	var result2 = sut.Bar1(dummyBar2)

	// assert
	assertEquals(t, dummyResult1, result1, "(*foo2).Bar1 1st call result different")
	assertEquals(t, dummyResult2, result2, "(*foo2).Bar1 2nd call result different")
	assertEquals(t, 0, m.MethodCalledCount(&foo2{}, "bar1"), "not mocked function count calculated wrongly")
}

func TestMocker_ShouldMockPublicMethodOnStructPointerReceiver(t *testing.T) {
	// arrange
	var dummySelf = rand.Intn(100)
	var dummyBar = rand.Intn(100)
	var dummyResult = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// SUT
	var sut = &foo2{
		self: dummySelf,
	}

	// expect
	m.ExpectMethod(&foo2{}, "Bar1", 1, func(obj *foo2, bar int) int {
		assertEquals(t, sut, obj, "(*foo2).Bar1 call instance different")
		assertEquals(t, dummyBar, bar, "(*foo2).Bar1 call parameter different")
		return dummyResult
	})

	// act
	var result = sut.Bar1(dummyBar)

	// assert
	assertEquals(t, dummyResult, result, "(*foo2).Bar1 call result different")
}

func TestMocker_ShouldMockPublicMethodOnStructPointerReceiverMultipleTimes(t *testing.T) {
	// arrange
	var dummySelf = rand.Intn(100)
	var dummyBar1 = rand.Intn(100)
	var dummyBar2 = rand.Intn(100)
	var dummyResult1 = rand.Intn(100)
	var dummyResult2 = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// SUT
	var sut = &foo2{
		self: dummySelf,
	}

	// expect
	m.ExpectMethod(&foo2{}, "Bar1", 1, func(obj *foo2, bar int) int {
		assertEquals(t, sut, obj, "(*foo2).Bar1 call instance different")
		assertEquals(t, dummyBar1, bar, "(*foo2).Bar1 call parameter different")
		return dummyResult1
	}).ExpectMethod(&foo2{}, "Bar1", 1, func(obj *foo2, bar int) int {
		assertEquals(t, sut, obj, "(*foo2).Bar1 call instance different")
		assertEquals(t, dummyBar2, bar, "(*foo2).Bar1 call parameter different")
		return dummyResult2
	})

	// act
	var result1 = sut.Bar1(dummyBar1)
	var result2 = sut.Bar1(dummyBar2)

	// assert
	assertEquals(t, dummyResult1, result1, "(*foo2).Bar1 1st call result different")
	assertEquals(t, dummyResult2, result2, "(*foo2).Bar1 2nd call result different")
}

func TestMocker_ShouldMockPublicMethodOnStructPointerReceiverMultipleTimes_UseCounter(t *testing.T) {
	// arrange
	var dummySelf = rand.Intn(100)
	var dummyBar1 = rand.Intn(100)
	var dummyBar2 = rand.Intn(100)
	var dummyResult1 = rand.Intn(100)
	var dummyResult2 = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// SUT
	var sut = &foo2{
		self: dummySelf,
	}

	// expect
	m.ExpectMethod(&foo2{}, "Bar1", 2, func(obj *foo2, bar int) int {
		if m.MethodCalledCount(&foo2{}, "Bar1") == 1 {
			assertEquals(t, sut, obj, "(*foo2).Bar1 call instance different")
			assertEquals(t, dummyBar1, bar, "(*foo2).Bar1 call parameter different")
			return dummyResult1
		} else if m.MethodCalledCount(&foo2{}, "Bar1") == 2 {
			assertEquals(t, sut, obj, "(*foo2).Bar1 call instance different")
			assertEquals(t, dummyBar2, bar, "(*foo2).Bar1 call parameter different")
			return dummyResult2
		}
		t.Errorf("(*foo2).Bar1 called more than twice")
		return -1
	})

	// act
	var result1 = sut.Bar1(dummyBar1)
	var result2 = sut.Bar1(dummyBar2)

	// assert
	assertEquals(t, dummyResult1, result1, "(*foo2).Bar1 1st call result different")
	assertEquals(t, dummyResult2, result2, "(*foo2).Bar1 2nd call result different")
	assertEquals(t, 0, m.MethodCalledCount(&foo2{}, "bar1"), "not mocked function count calculated wrongly")
}

func (f foo2) bar2(bar int) int {
	return f.self * bar
}

func TestMocker_ShouldMockPrivateFunctionOnStructValueReceiver(t *testing.T) {
	// arrange
	var dummySelf = rand.Intn(100)
	var dummyBar = rand.Intn(100)
	var dummyResult = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// SUT
	var sut = foo2{
		self: dummySelf,
	}

	// expect
	m.ExpectFunc(foo2.bar2, 1, func(obj foo2, bar int) int {
		assertEquals(t, sut, obj, "foo2.bar2 call instance different")
		assertEquals(t, dummyBar, bar, "foo2.bar2 call parameter different")
		return dummyResult
	})

	// act
	var result = sut.bar2(dummyBar)

	// assert
	assertEquals(t, dummyResult, result, "foo2.bar2 call result different")
}

func TestMocker_ShouldMockPrivateMethodOnStructValueReceiver(t *testing.T) {
	// arrange
	var dummySelf = rand.Intn(100)
	var dummyBar = rand.Intn(100)
	var dummyResult = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// SUT
	var sut = foo2{
		self: dummySelf,
	}

	// expect
	m.ExpectMethod(foo2{}, "bar2", 1, func(obj foo2, bar int) int {
		assertEquals(t, sut, obj, "foo2.bar2 call instance different")
		assertEquals(t, dummyBar, bar, "foo2.bar2 call parameter different")
		return dummyResult
	})

	// act
	var result = sut.bar2(dummyBar)

	// assert
	assertEquals(t, dummyResult, result, "foo2.bar2 call result different")
}

func (f foo2) Bar2(bar int) int {
	return f.self * bar
}

func TestMocker_ShouldMockPublicFunctionOnStructValueReceiver(t *testing.T) {
	// arrange
	var dummySelf = rand.Intn(100)
	var dummyBar = rand.Intn(100)
	var dummyResult = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// SUT
	var sut = foo2{
		self: dummySelf,
	}

	// expect
	m.ExpectFunc(foo2.Bar2, 1, func(obj foo2, bar int) int {
		assertEquals(t, sut, obj, "foo2.Bar2 call instance different")
		assertEquals(t, dummyBar, bar, "foo2.Bar2 call parameter different")
		return dummyResult
	})

	// act
	var result = sut.Bar2(dummyBar)

	// assert
	assertEquals(t, dummyResult, result, "foo2.Bar2 call result different")
}

func TestMocker_ShouldMockPublicMethodOnStructValueReceiver(t *testing.T) {
	// arrange
	var dummySelf = rand.Intn(100)
	var dummyBar = rand.Intn(100)
	var dummyResult = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// SUT
	var sut = foo2{
		self: dummySelf,
	}

	// expect
	m.ExpectMethod(foo2{}, "Bar2", 1, func(obj foo2, bar int) int {
		assertEquals(t, sut, obj, "foo2.Bar2 call instance different")
		assertEquals(t, dummyBar, bar, "foo2.Bar2 call parameter different")
		return dummyResult
	})

	// act
	var result = sut.Bar2(dummyBar)

	// assert
	assertEquals(t, dummyResult, result, "foo2.Bar2 call result different")
}

type Foo3 interface {
	Bar(bar int) int
}

type foo3 struct {
	Foo3
}

func TestMocker_ShouldMockPublicMethodOnInterfacePointer(t *testing.T) {
	// arrange
	var dummyBar = rand.Intn(100)
	var dummyResult = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// SUT
	var sut = &foo3{}

	// expect
	m.ExpectMethod(&foo3{}, "Bar", 1, func(obj *foo3, bar int) int {
		assertEquals(t, sut, obj, "Foo3.Bar call instance different")
		assertEquals(t, dummyBar, bar, "Foo3.Bar call parameter different")
		return dummyResult
	})

	// act
	var result = func(f Foo3, bar int) int {
		return f.Bar(bar)
	}(sut, dummyBar)

	// assert
	assertEquals(t, dummyResult, result, "Foo3.Bar call result different")
}

func foo4(bar ...int) int {
	return len(bar)
}

func TestMocker_ShouldMockPrivateVariadicFunction(t *testing.T) {
	// arrange
	var dummyBar = []int{rand.Intn(100), rand.Intn(100), rand.Intn(100)}
	var dummyResult = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// expect
	m.ExpectFunc(foo4, 1, func(bar ...int) int {
		assertEquals(t, len(dummyBar), len(bar), "foo4 call parameter count different")
		assertEquals(t, dummyBar[0], bar[0], "foo4 call parameter 1 different")
		assertEquals(t, dummyBar[1], bar[1], "foo4 call parameter 2 different")
		assertEquals(t, dummyBar[2], bar[2], "foo4 call parameter 3 different")
		return dummyResult
	})

	// SUT + act
	var result = foo4(dummyBar...)

	// assert
	assertEquals(t, dummyResult, result, "foo4 call result different")
}

func Foo4(bar ...int) int {
	return len(bar)
}

func TestMocker_ShouldMockPublicVariadicFunction(t *testing.T) {
	// arrange
	var dummyBar = []int{rand.Intn(100), rand.Intn(100), rand.Intn(100)}
	var dummyResult = rand.Intn(100)

	// mock
	var m = NewMocker(t)

	// expect
	m.ExpectFunc(Foo4, 1, func(bar ...int) int {
		assertEquals(t, len(dummyBar), len(bar), "Foo4 call parameter count different")
		assertEquals(t, dummyBar[0], bar[0], "Foo4 call parameter 1 different")
		assertEquals(t, dummyBar[1], bar[1], "Foo4 call parameter 2 different")
		assertEquals(t, dummyBar[2], bar[2], "Foo4 call parameter 3 different")
		return dummyResult
	})

	// SUT + act
	var result = Foo4(dummyBar...)

	// assert
	assertEquals(t, dummyResult, result, "Foo4 call result different")
}

type tester struct {
	testing.TB
	t      *testing.T
	errorf func(string, ...interface{})
}

func (t *tester) Errorf(format string, args ...interface{}) {
	t.errorf(format, args...)
}

func (t *tester) Cleanup(f func()) {
	t.t.Cleanup(f)
}

func (t *tester) Helper() {
	t.t.Helper()
}

func TestMocker_ShouldReportTestFailureWhenMockFunctionIsNotCalledButExpected(t *testing.T) {
	// arrange
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
	m.ExpectFunc(foo1, 1, func(int) int {
		return 0
	})
}

func TestMocker_ShouldReportTestFailureWhenMockFunctionIsCalledActualLessThanExpected(t *testing.T) {
	// arrange
	var tester = &tester{t: t}

	// mock
	var m = NewMocker(tester)

	// expect
	tester.errorf = func(format string, args ...interface{}) {
		assertEquals(t, "[%v] Unepxected number of calls: expect %v, actual %v", format, "tester.Errorf called with different message")
		assertEquals(t, 3, len(args), "tester.Errorf called with different number of args")
		assertEquals(t, 2, args[1], "tester.Errorf called with different argument 2")
		assertEquals(t, 1, args[2], "tester.Errorf called with different argument 3")
	}
	m.ExpectFunc(foo1, 2, func(int) int {
		return 0
	})

	// SUT + act
	_ = foo1(0)
}

func TestMocker_ShouldReportTestFailureWhenMockFunctionIsCalledActualMoreThanExpected(t *testing.T) {
	// arrange
	var tester = &tester{t: t}

	// mock
	var m = NewMocker(tester)

	// expect
	tester.errorf = func(format string, args ...interface{}) {
		assertEquals(t, "[%v] Unepxected number of calls: expect %v, actual %v", format, "tester.Errorf called with different message")
		assertEquals(t, 3, len(args), "tester.Errorf called with different number of args")
		assertEquals(t, 1, args[1], "tester.Errorf called with different argument 2")
		assertEquals(t, 2, args[2], "tester.Errorf called with different argument 3")
	}
	m.ExpectFunc(foo1, 1, func(int) int {
		return 0
	})

	// SUT + act
	_ = foo1(0)
	_ = foo1(0)
}

func TestMocker_ShouldReportTestFailureWhenMockFunctionIsCalledButNotExpected(t *testing.T) {
	// arrange
	var tester = &tester{t: t}

	// mock
	var m = NewMocker(tester)

	// expect
	tester.errorf = func(format string, args ...interface{}) {
		assertEquals(t, "[%v] Unepxected number of calls: expect %v, actual %v", format, "tester.Errorf called with different message")
		assertEquals(t, 3, len(args), "tester.Errorf called with different number of args")
		assertEquals(t, 0, args[1], "tester.Errorf called with different argument 1")
		assertEquals(t, 1, args[2], "tester.Errorf called with different argument 2")
	}
	m.ExpectFunc(foo1, 0, func(int) int {
		return 0
	})

	// SUT + act
	_ = foo1(0)
}

func TestMocker_ShouldReportTestFailureWhenMockFunctionWithWrongExpectCount(t *testing.T) {
	// arrange
	var tester = &tester{t: t}

	// mock
	var m = NewMocker(tester)

	// expect
	tester.errorf = func(format string, args ...interface{}) {
		assertEquals(t, "[%v] Expect count must be greather than 0 when already mocked with other expectations", format, "tester.Errorf called with different message")
		assertEquals(t, 1, len(args), "tester.Errorf called with different number of args")
	}
	m.ExpectFunc(foo1, 0, func(int) int {
		return 0
	}).ExpectFunc(foo1, 0, func(int) int {
		return 0
	})
}

func TestMocker_ShouldReportTestFailureWhenMockFunctionPanicsInExecution(t *testing.T) {
	defer func() {
		recover()
	}()

	// arrange
	var tester = &tester{t: t}

	// mock
	var m = NewMocker(tester)

	// expect
	tester.errorf = func(format string, args ...interface{}) {
		assertEquals(t, "[%v] Mocker panicing recovered: %v", format, "tester.Errorf called with different message")
		assertEquals(t, 2, len(args), "tester.Errorf called with different number of args")
		assertEquals(t, "paniced", args[1], "tester.Errorf called with different argument 2")
	}
	m.ExpectFunc(foo1, 1, func(int) int {
		panic(errors.New("paniced"))
	})

	// SUT + act
	_ = foo1(0)
}

func TestMocker_ShouldReportTestFailureWhenMockMethodIsNotCalledButExpected(t *testing.T) {
	// arrange
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
	m.ExpectMethod(&foo2{}, "bar1", 1, func(int) int {
		return 0
	})
}

func TestMocker_ShouldReportTestFailureWhenMockMethodIsCalledActualLessThanExpected(t *testing.T) {
	// arrange
	var tester = &tester{t: t}

	// mock
	var m = NewMocker(tester)

	// expect
	tester.errorf = func(format string, args ...interface{}) {
		assertEquals(t, "[%v] Unepxected number of calls: expect %v, actual %v", format, "tester.Errorf called with different message")
		assertEquals(t, 3, len(args), "tester.Errorf called with different number of args")
		assertEquals(t, 2, args[1], "tester.Errorf called with different argument 2")
		assertEquals(t, 1, args[2], "tester.Errorf called with different argument 3")
	}
	m.ExpectMethod(&foo2{}, "bar1", 2, func(int) int {
		return 0
	})

	// SUT
	var sut = &foo2{}

	// act
	_ = sut.bar1(0)
}

func TestMocker_ShouldReportTestFailureWhenMockMethodIsCalledActualMoreThanExpected(t *testing.T) {
	// arrange
	var tester = &tester{t: t}

	// mock
	var m = NewMocker(tester)

	// expect
	tester.errorf = func(format string, args ...interface{}) {
		assertEquals(t, "[%v] Unepxected number of calls: expect %v, actual %v", format, "tester.Errorf called with different message")
		assertEquals(t, 3, len(args), "tester.Errorf called with different number of args")
		assertEquals(t, 1, args[1], "tester.Errorf called with different argument 2")
		assertEquals(t, 2, args[2], "tester.Errorf called with different argument 3")
	}
	m.ExpectMethod(&foo2{}, "bar1", 1, func(int) int {
		return 0
	})

	// SUT
	var sut = &foo2{}

	// act
	_ = sut.bar1(0)
	_ = sut.bar1(0)
}

func TestMocker_ShouldReportTestFailureWhenMockMethodIsCalledButNotExpected(t *testing.T) {
	// arrange
	var tester = &tester{t: t}

	// mock
	var m = NewMocker(tester)

	// expect
	tester.errorf = func(format string, args ...interface{}) {
		assertEquals(t, "[%v] Unepxected number of calls: expect %v, actual %v", format, "tester.Errorf called with different message")
		assertEquals(t, 3, len(args), "tester.Errorf called with different number of args")
		assertEquals(t, 0, args[1], "tester.Errorf called with different argument 2")
		assertEquals(t, 1, args[2], "tester.Errorf called with different argument 3")
	}
	m.ExpectMethod(&foo2{}, "bar1", 0, func(int) int {
		return 0
	})

	// SUT
	var sut = &foo2{}

	// act
	_ = sut.bar1(0)
}

func TestMocker_ShouldReportTestFailureWhenMockMethodWithWrongExpectCount(t *testing.T) {
	// arrange
	var tester = &tester{t: t}

	// mock
	var m = NewMocker(tester)

	// expect
	tester.errorf = func(format string, args ...interface{}) {
		assertEquals(t, "[%v] Expect count must be greather than 0 when already mocked with other expectations", format, "tester.Errorf called with different message")
		assertEquals(t, 1, len(args), "tester.Errorf called with different number of args")
	}
	m.ExpectMethod(&foo2{}, "bar1", 0, func(int) int {
		return 0
	}).ExpectMethod(&foo2{}, "bar1", 0, func(int) int {
		return 0
	})
}

func TestMocker_ShouldReportTestFailureWhenMockMethodPanicsInExecution(t *testing.T) {
	defer func() {
		recover()
	}()

	// arrange
	var tester = &tester{t: t}

	// mock
	var m = NewMocker(tester)

	// expect
	tester.errorf = func(format string, args ...interface{}) {
		assertEquals(t, "[%v] Mocker panicing recovered: %v", format, "tester.Errorf called with different message")
		assertEquals(t, 2, len(args), "tester.Errorf called with different number of args")
		assertEquals(t, "paniced", args[1], "tester.Errorf called with different argument 2")
	}
	m.ExpectMethod(&foo2{}, "bar1", 1, func(int) int {
		panic("paniced")
	})

	// SUT
	var sut = &foo2{}

	// act
	_ = sut.bar1(0)
}

func TestMock_ShouldReportTestFailureWhenMockPrivateMethodNotFound(t *testing.T) {
	// arrange
	var tester = &tester{t: t}

	// mock
	var m = NewMocker(tester)

	// expect
	tester.errorf = func(format string, args ...interface{}) {
		assertEquals(t, "Method [%v] cannot be located for given target struct [%v]", format, "tester.Errorf called with different message")
		assertEquals(t, 2, len(args), "tester.Errorf called with different number of args")
		assertEquals(t, "nobar", args[0], "tester.Errorf called with different argument 1")
	}
	m.ExpectMethod(&foo2{}, "nobar", 0, func(int) int {
		return 0
	})
}

func TestMock_ShouldReportTestFailureWhenMockPublicMethodNotFound(t *testing.T) {
	// arrange
	var tester = &tester{t: t}

	// mock
	var m = NewMocker(tester)

	// expect
	tester.errorf = func(format string, args ...interface{}) {
		assertEquals(t, "Method [%v] cannot be located for given target struct [%v]", format, "tester.Errorf called with different message")
		assertEquals(t, 2, len(args), "tester.Errorf called with different number of args")
		assertEquals(t, "Nobar", args[0], "tester.Errorf called with different argument 1")
	}
	m.ExpectMethod(&foo2{}, "Nobar", 0, func(int) int {
		return 0
	})
}

func TestMocker_ShouldHandleEntryNotFoundScenarioWhenMakeFunc(t *testing.T) {
	// arrange
	var dummyName = "some name"
	var dummyFuncPtr = uintptr(rand.Intn(100))
	var dummyMockFunc = func(bar int) int { return bar }
	var dummyBar = rand.Intn(100)
	var tester = &tester{t: t}

	// mock
	var m = NewMocker(tester).(*mocker)

	// expect
	tester.errorf = func(format string, args ...interface{}) {
		assertEquals(t, "[%v] Unepxected number of calls: expect %v, actual %v", format, "tester.Errorf called with different message")
		assertEquals(t, 3, len(args), "tester.Errorf called with different number of args")
		assertEquals(t, dummyName, args[0], "tester.Errorf called with different argument 1")
		assertEquals(t, 0, args[1], "tester.Errorf called with different argument 2")
		assertEquals(t, 1, args[2], "tester.Errorf called with different argument 3")
	}

	// SUT
	var sut = m.makeFunc(dummyName, dummyFuncPtr, dummyMockFunc)

	// act
	sut.Call([]reflect.Value{reflect.ValueOf(dummyBar)})

	// assert
}

func assertEquals(t *testing.T, expect interface{}, actual interface{}, message string) {
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
