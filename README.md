# caststructure [![Godoc](https://godoc.org/github.com/mitchellh/caststructure?status.svg)](https://godoc.org/github.com/mitchellh/caststructure)

caststructure is a Go library that safely casts a value to a type
that implements a dynamic list of interfaces.

For example, you get a value that implements interfaces `A`, `B`, `C` and you
want to return a value that only implements `A` and `B`. If this is the only
scenario, you can manually create a [new interface type](https://play.golang.org/p/Sgn7MhsyrXt).
But if the list of interfaces you want to implement is dynamic, you would
have to declare all combinations of the interface. This library does this
dynamically at runtime, without panics.

## Installation

Standard `go get`:

```
$ go get github.com/mitchellh/caststructure
```

## Usage & Example

For usage and examples see the [Godoc](http://godoc.org/github.com/mitchellh/caststructure).

A quick code example is shown below:

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
newValue := caststructure.Must(caststructure.Interface(value, (*A)(nil), (*B)(nil)))
_, ok := newValue.(A) // ok == true
_, ok := newValue.(B) // ok == true
_, ok := newValue.(C) // ok == false
```
