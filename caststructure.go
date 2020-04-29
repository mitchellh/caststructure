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

// Interface converts the value "from" to a type that implements the list
// of interfaces "to". The "from" value must already implement these interfaces.
// This is typically used to downcast a value that implements other interfaces
// to one that only implements the given list.
//
// An error will be returned if from doesn't implement any one of the to values,
// or creating the structure fails (due to overlapping interface implements,
// unexported types being used, etc.).
func Interface(from interface{}, to ...interface{}) (interface{}, error) {
	fromVal := reflect.ValueOf(from)
	fromTyp := fromVal.Type()

	// Go through all our to interfaces and build up our struct field list.
	// We also do some validation here.
	fields := make([]reflect.StructField, len(to))
	for i, typPtr := range to {
		typ := reflect.TypeOf(typPtr)
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
