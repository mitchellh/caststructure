// Package caststructure can downcast a value to a subset of the interfaces
// that it implements. This is useful in situations where the set of interfaces
// you want a type to implement is dynamic. In cases where the set is small or
// fixed, it is likely a better idea to define a new type instead.
package caststructure

import (
	"fmt"
	"reflect"
)

// Must is a helper that wraps a call to a function returning (interface{}, error)
// and panics if the error is non-nil. This can be used around Interface to
// force a successful return value or panic.
func Must(v interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}

	return v
}

// Down downcasts the value "from" to a type that implements the list
// of interfaces "to". The "from" value must already implement these interfaces.
// This is typically used to downcast a value that implements other interfaces
// to one that only implements the given list.
//
// An error will be returned if from doesn't implement any one of the to values,
// or creating the structure fails (due to overlapping interface implements,
// unexported types being used, etc.).
func Down(from interface{}, to ...interface{}) (interface{}, error) {
	fromVal := reflect.ValueOf(from)
	fromTyp := fromVal.Type()

	// Go through all our to interfaces and build up our struct field list.
	// We also do some validation here.
	fields := make([]reflect.StructField, len(to))
	for i, typPtr := range to {
		typ := reflect.TypeOf(typPtr)
		if typ == nil {
			return nil, fmt.Errorf("to type must be a pointer, got nil value")
		}
		if typ.Kind() != reflect.Ptr {
			return nil, fmt.Errorf("to type must be a pointer, got %s", typ.String())
		}
		typ = typ.Elem()

		// Ensure our from value maps to this
		if !fromTyp.AssignableTo(typ) {
			return nil, fmt.Errorf("from value is not assignable to destination type %s", typ.String())
		}

		// Build the field
		fields[i] = reflect.StructField{
			Name:      "A_" + typ.Name(), // Force capitalization so its exported
			Type:      typ,
			Anonymous: true, // Anonymous so it is embedded
		}
	}

	// Make our structure
	resultVal := reflect.New(reflect.StructOf(fields)).Elem()
	for i := 0; i < len(to); i++ {
		resultVal.Field(i).Set(fromVal)
	}

	return resultVal.Interface(), nil
}

// Compose composes multiple types into a single type. For example if you
// have a value that implements A and another that implements B, this can
// compose those two values into a single value that implements both A and B.
//
// The arguments to this must be a list of pairs in (value, type) order.
// For example:
//
//     Compose(a, (*A)(nil), b, (*B)(nil))
//
func Compose(pairs ...interface{}) (interface{}, error) {
	if len(pairs)%2 != 0 {
		return nil, fmt.Errorf("Compose requires an even number of arguments since they're pairs")
	}

	fields := make([]reflect.StructField, 0, len(pairs)/2)
	values := make([]reflect.Value, 0, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		pairVal := reflect.ValueOf(pairs[i])
		pairTyp := pairs[i+1]

		typ := reflect.TypeOf(pairTyp)
		if typ == nil {
			return nil, fmt.Errorf("to type must be a pointer, got nil value")
		}
		if typ.Kind() != reflect.Ptr {
			return nil, fmt.Errorf("to type must be a pointer, got %s", typ.String())
		}
		typ = typ.Elem()

		// Ensure our from value maps to this
		if !pairVal.Type().AssignableTo(typ) {
			return nil, fmt.Errorf("value is not assignable to destination type %s", typ.String())
		}

		// Build the field
		fields = append(fields, reflect.StructField{
			Name:      "A_" + typ.Name(), // Force capitalization so its exported
			Type:      typ,
			Anonymous: true, // Anonymous so it is embedded
		})

		// Build the value
		values = append(values, pairVal)
	}

	// Make our structure
	resultVal := reflect.New(reflect.StructOf(fields)).Elem()
	for i := 0; i < len(fields); i++ {
		resultVal.Field(i).Set(values[i])
	}

	return resultVal.Interface(), nil
}
