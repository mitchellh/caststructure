# caststructure [![Godoc](https://godoc.org/github.com/mitchellh/caststructure?status.svg)](https://godoc.org/github.com/mitchellh/caststructure)

caststructure is a Go library that provides functions for downcasting types,
composing values dynamically, and more. See the examples below for more details.

## Installation

Standard `go get`:

```
$ go get github.com/mitchellh/caststructure
```

## Usage & Example

For usage and examples see the [Godoc](http://godoc.org/github.com/mitchellh/caststructure).

A quick code example is shown below that shows **downcasting**:

```go
// Three interface types.
type A interface { A() int }
type B interface { B() int }
type C interface { C() int }

// Impl implements A, B, AND C.
type Impl struct {}
func (Impl) A() int { return 42 }
func (Impl) B() int { return 42 }
func (Impl) C() int { return 42 }

// We have a value of type Impl, so it implements all interfaces.
var value Impl

// But we only want value to implement A and B, not C.
newValue := caststructure.Must(caststructure.Down(value, (*A)(nil), (*B)(nil)))
_, ok := newValue.(A) // ok == true
_, ok := newValue.(B) // ok == true
_, ok := newValue.(C) // ok == false
```

Here is an example that shows **composing**:

```go
// Three interface types.
type A interface { A() int }
type B interface { B() int }
type C interface { C() int }

// Implementation for A and B, respectively
type ImplA struct {}
func (ImplA) A() int { return 42 }

type ImplB struct {}
func (Impl) B() int { return 42 }

// We have a value that implements A and B, separately.
var valueA implA
var valueB implB

// But we want a value that implements BOTH A and B.
newValue := caststructure.Must(caststructure.Compose(valueA, (*A)(nil), valueB, (*B)(nil)))
_, ok := newValue.(A) // ok == true
_, ok := newValue.(B) // ok == true
```

## But... Why?

In general, what this library achieves can be done manually through
explicit type declarations. For example, composing `A` and `B` can be
done [explicitly like this](https://play.golang.org/p/It0NHvZt_-w). But in
cases where the set of interfaces is large and the combinations dynamic,
creating explicit types becomes a burden.

Some Go code allows optionally implementing interfaces to change behavior.
For example, code might check if `io.Closer` is implemented and call it. If
it is not implemented, it isn't an error, it just isn't called. Taking this
further, some Go code has numerous opt-in interfaces to change behavior. This
library was born out of the need to work with such APIs.

For downcasing, you might get a value that implements interfaces `A`, `B`, `C` and you
want to return a value that only implements `A` and `B`. If this is the only
scenario, you can manually create a [new interface type](https://play.golang.org/p/Sgn7MhsyrXt).
But if the list of interfaces you want to implement is dynamic, you would
have to declare all combinations of the interface. This library does this
dynamically at runtime, without panics.

For composing, you may have two separate values that implement two separate
interfaces `A` and `B` respectively. For example, you may have an `io.Reader`
and an `io.Writer` but you want an `io.ReadWriter` using the two together.
This library can construct that ReadWriter for you.
