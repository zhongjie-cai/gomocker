# gomocker

[![Test](https://github.com/zhongjie-cai/gomocker/actions/workflows/ci.yaml/badge.svg)](https://github.com/zhongjie-cai/gomocker/actions/workflows/ci.yaml)
![Coverage](https://img.shields.io/badge/Coverage-100.0%25-brightgreen)
[![Go Report Card](https://goreportcard.com/badge/github.com/zhongjie-cai/gomocker)](https://goreportcard.com/report/github.com/zhongjie-cai/gomocker)
[![Go Reference](https://pkg.go.dev/badge/github.com/zhongjie-cai/gomocker.svg)](https://pkg.go.dev/github.com/zhongjie-cai/gomocker)

A mocker library for Go inspired by gomonkey features, allowing developers to mock functions or struct methods for unit tests.

**Important Note: must set the build flag `-gcflags=all=-l` so as to make this library properly functional.**

- [gomocker](#gomocker)
    - [Scenario 1 - mock a function (either private or public, as long as accessible)](#scenario-1---mock-a-function-either-private-or-public-as-long-as-accessible)
    - [Scenario 2 - mock a struct method (either private or public, as long as accessible)](#scenario-2---mock-a-struct-method-either-private-or-public-as-long-as-accessible)
    - [Scenario 3 - mock a public interface method](#scenario-3---mock-a-public-interface-method)
    - [Scenario 4 - mock a function / method with side effects](#scenario-4---mock-a-function--method-with-side-effects)
    - [Scenario 5 - mock a function / method to be not called](#scenario-5---mock-a-function--method-to-be-not-called)
    - [Scenario 6 - bypass parameter matching](#scenario-6---bypass-parameter-matching)
    - [Scenario 7 - customize parameter matching](#scenario-7---customize-parameter-matching)

### Scenario 1 - mock a function (either private or public, as long as accessible)

With the following function `foo` in code:

```go
func foo(bar int) int {
	return bar * 2
}
```

One can mock it with the following code:

```go
// mock
var m = gomocker.NewMocker(t)

// expect
m.Mock(
    foo
).Expects(
    // place your expected parameters here
).Returns(
    // place your anticipated returns here
).Once(
    // or choose Twice, Times method instead, this function must be called to complete a Mock or Stub
)
```

### Scenario 2 - mock a struct method (either private or public, as long as accessible)

With the following struct `foo` with method `bar` in code:

```go
type foo struct {
    self int
}

func (f *foo) bar(val int) int {
    return val * self
}
```

One can mock it with the following code using `ExpectMethod`:

```go
// arrange
var f = &foo{}

// mock
var m = gomocker.NewMocker(t)

// expect
m.Mock(
    (*foo).bar
).Expects(
    // the first parameter should be the exact instance of the struct foo that initiates the method call
    //   e.g. `f` in this example
    // followed by other expected parameters here
).Returns(
    // place your anticipated returns here
).Once(
    // or choose Twice, Times method instead, this function must be called to complete a Mock or Stub
)
```

### Scenario 3 - mock a public interface method

With the following interface `Foo` with method `Bar` of package `example` in code:

```go
package example

type Foo interface {
    Bar(val int) int
}
```

And the calling function `doSomething` in package `main` in code:

```go
package main

func doSomething(f example.Foo) int {
    return f.Bar(123)
}
```

One can mock it with the following code using `ExpectMethod` when testing `doSomething`:

```go
// arrange
type dummyFoo struct {
    example.Foo
}
var f = &dummyFoo{}

// mock
var m = gomocker.NewMocker(t)

// expect
m.Mock(
    (*dummyFoo).Bar
).Expects(
    // the first parameter should be the exact instance of the struct dummyFoo that initiates the method call
    //   e.g. `f` in this example
    // followed by other expected parameters here
).Returns(
    // place your anticipated returns here
).Once(
    // or choose Twice, Times method instead, this function must be called to complete a Mock or Stub
)
```

### Scenario 4 - mock a function / method with side effects

```go
// mock
var m = gomocker.NewMocker(t)

// expect
m.Mock(
    foo
).Expects(
    // place your expected parameters here
).Returns(
    // place your anticipated returns here
).SideEffect(
    func(index int) {
        // place your side effect code logic here
        // the given parameter `index` means the number of calls (including the current call)
        //   to the mocked or stubbed function happened so far
	}
).Once(
    // or choose Twice, Times method instead, this function must be called to complete a Mock or Stub
)
```

### Scenario 5 - mock a function / method to be not called

```go
// mock
var m = gomocker.NewMocker(t)

// expect
m.Mock(foo).NotCalled(
    // this completes a mock, and if `foo` is called, the test shall fail.
    //   note that a function or method cannot be mocked or stubbed again if it is set to NotCalled
)
```

### Scenario 6 - bypass parameter matching

```go
// arrange
var foo = func(int) {}

// mock
var m = gomocker.NewMocker(t)

// expect
m.Mock(foo).Expects(
    gomocker.Anything(), // this allows bypassing the value check for a particular parameter
).Returns()
```

### Scenario 7 - customize parameter matching

```go
// arrange
var foo = func(int) {}

// mock
var m = gomocker.NewMocker(t)

// expect
m.Mock(foo).Expects(
    gomocker.Matches(func(value interface{}) bool) {
        // this allows customization of the check for a particular parameter
        //   the original parameter is wrapped into an interface and is given as `value` here
        //   returning false would cause the corresponding test to fail
    },
).Returns()
```
