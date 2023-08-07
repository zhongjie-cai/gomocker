# gomocker

[![Go Report Card](https://goreportcard.com/badge/github.com/zhongjie-cai/mocker)](https://goreportcard.com/report/github.com/zhongjie-cai/mocker)
[![Go Reference](https://pkg.go.dev/badge/github.com/zhongjie-cai/gomocker.svg)](https://pkg.go.dev/github.com/zhongjie-cai/gomocker)

A mocker library for Go based on gomonkey features, allowing developers to mock either functions or struct methods according to unit test needs.

### Scenario 1 - Mock a private function

With the following function `foo` in code:

```go
func foo(bar int) int {
	return bar * 2
}
```

One can mock it with the following code:

```go
// mock
var m = NewMocker(t)

// expect
m.ExpectFunc(foo, 1, func(bar int) int {
    // fill in with your own assertions and return
})
```

### Scenario 2 - Mock a public function

With the following function `Foo` of package `example` in code:

```go
package example

func Foo(bar int) int {
	return bar * 2
}
```

One can mock it with the following code:

```go
// mock
var m = NewMocker(t)

// expect
m.ExpectFunc(example.Foo, 1, func(bar int) int {
    // fill in with your own assertions and return
})
```

### Scenario 3 - Mock a function multiple times

With the following function `foo` in code:

```go
func foo(bar int) int {
	return bar * 2
}
```

One can mock it with the following code when called multiple times:

```go
// mock
var m = NewMocker(t)

// expect
m.ExpectFunc(foo, 2, func(bar int) int { // expect this mockFunc to be executed twice
    // fill in with your own assertions and return for the first two calls
}).ExpectFunc(foo, 1, func(bar int) int { // expect this mockFunc to be executed once
    // fill in with your own assertions and return for the third call
})
```

### Scenario 4 - mock a private struct method

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
// mock
var m = NewMocker(t)

// expect
m.ExpectMethod(&foo{}, "bar", 1, func(obj *foo, val int) int { // obj must be the same pointer type as the targetStruct
    // fill in with your own assertions and return
})
```

Or can mock it with the following code using `ExpectFunc`:

```go
// mock
var m = NewMocker(t)

// expect
m.ExpectFunc((*foo).bar, 1, func(obj *foo, val int) int { // obj must be the same pointer type as the expectFunc's owner struct
    // fill in with your own assertions and return
})
```

Mocking value receivers works similar to pointer receivers, only to make sure the targetStruct or expectFunc points to a struct's value instance, and the `obj` parameter's type should change accordingly.

### Scenario 5 - mock a public struct method

With the following struct `Foo` with method `Bar` of package `example` in code:

```go
package example

type Foo struct {
    self int
}

func (f *Foo) Bar(val int) int {
    return val * self
}
```

One can mock it with the following code using `ExpectMethod`:

```go
// mock
var m = NewMocker(t)

// expect
m.ExpectMethod(&example.Foo{}, "Bar", 1, func(obj *example.Foo, val int) int { // obj must be the same pointer type as the targetStruct
    // fill in with your own assertions and return
})
```

Or can mock it with the following code using `ExpectFunc`:

```go
// mock
var m = NewMocker(t)

// expect
m.ExpectFunc((*example.Foo).Bar, 1, func(obj *example.Foo, val int) int { // obj must be the same pointer type as the expectFunc's owner struct
    // fill in with your own assertions and return
})
```

Mocking value receivers works similar to pointer receivers, only to make sure the targetStruct or expectFunc points to a struct's value instance, and the `obj` parameter's type should change accordingly.

### Scenario 6 - mock a public interface method

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
var m = NewMocker(t)

// expect
m.ExpectMethod(&dummyFoo{}, "Bar", 1, func(obj *dummyFoo, val int) int { // obj must be the same pointer type as the targetStruct
    // fill in with your own assertions and return
})
```
