package caststructure

import (
	"testing"
)

func TestDown(t *testing.T) {
	var from impl

	value, err := Down(from, (*testA)(nil), (*TestB)(nil))
	if err != nil {
		t.Fatal(err)
	}

	value.(testA).A()
	if v, ok := value.(testA); !ok {
		t.Fatal("should implement A")
	} else if v.A() != 42 {
		t.Fatal("invalid value")
	}
	if v, ok := value.(TestB); !ok {
		t.Fatal("should implement B")
	} else if v.B() != 42 {
		t.Fatal("invalid value")
	}
	if _, ok := value.(testC); ok {
		t.Fatal("should not implement C")
	}
}

func TestDown_nonImpl(t *testing.T) {
	from := 42
	_, err := Down(from, (*testA)(nil))
	if err == nil {
		t.Fatal("should error")
	}
}

func TestDown_nonPtr(t *testing.T) {
	var from impl
	_, err := Down(from, (testA)(nil))
	if err == nil {
		t.Fatal("should error")
	}
}

func TestCompose(t *testing.T) {
	var a implA
	var b implB

	value, err := Compose(a, (*testA)(nil), b, (*TestB)(nil))
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := value.(testA); !ok {
		t.Fatal("should implement A")
	}
	if _, ok := value.(TestB); !ok {
		t.Fatal("should implement B")
	}
	if _, ok := value.(testC); ok {
		t.Fatal("should not implement C")
	}
}

type testA interface{ A() int }
type TestB interface{ B() int } // Purposefully exported to test that case
type testC interface{ C() int }

type impl struct{}

func (impl) A() int { return 42 }
func (impl) B() int { return 42 }
func (impl) C() int { return 42 }

type implA struct{}

func (implA) A() int { return 42 }

type implB struct{}

func (implB) B() int { return 42 }
