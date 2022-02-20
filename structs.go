// Package checks
// If the error is "this type <%s {%s}> of field is not supported". Add processing of the type of this field.
package checks

import (
	"errors"
	"fmt"
	r "reflect"
)

var (
	ErrType        = errors.New("type not supported")
	ErrListEmpty   = errors.New("list is empty")
	ErrStringEmpty = errors.New("string is empty")
	ErrNumZero     = errors.New("number is zero")
)

// StructureFields recursively traverses the structure and checks that:
// - field values are not equal to an empty string or zero
// - slices and arrays are not empty.
func StructureFields(s interface{}, ignoreRawTypes bool) error {
	v := r.ValueOf(s)
	vType := v.Type()

	for i := 0; i < v.NumField(); i++ {
		if err := typesHandling(i, v, vType, ignoreRawTypes); err != nil {
			return err
		}
	}

	return nil
}

func typesHandling(i int, v r.Value, vType r.Type, ignoreRawTypes bool) error {
	switch v := r.ValueOf(v.Field(i).Interface()); v.Kind() {
	case r.Struct:
		if err := StructureFields(v.Interface(), ignoreRawTypes); err != nil {
			return err
		}

	case r.String:
		if v.String() == "" {
			return fmt.Errorf("field: %s {%s} = \"\"", vType.Field(i).Name, vType.Field(i).Tag)
		}

	case r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
		if v.Int() == 0 {
			return fmt.Errorf("field: %s {%s} = %d", vType.Field(i).Name, vType.Field(i).Tag, v.Int())
		}

	case r.Uint, r.Uint8, r.Uint16, r.Uint32, r.Uint64:
		if v.Uint() == 0 {
			return fmt.Errorf("field: %s {%s} = %d", vType.Field(i).Name, vType.Field(i).Tag, v.Uint())
		}

	case r.Float32, r.Float64:
		if v.Float() == 0 {
			return fmt.Errorf("field: %s {%s} = %f", vType.Field(i).Name, vType.Field(i).Tag, v.Float())
		}

	case r.Slice, r.Array:
		if err := checkFieldsLists(v); err != nil {
			return fmt.Errorf("field: %s {%s}: %w", vType.Field(i).Name, vType.Field(i).Tag, err)
		}

	case r.Bool:
		break

	default:
		if !ignoreRawTypes {
			return fmt.Errorf("this type <%s {%s}>: %w", vType.Field(i).Name, vType.Field(i).Tag, ErrType)
		}
	}

	return nil
}

func checkFieldsLists(rv r.Value) error {
	if rv.Len() == 0 {
		return ErrListEmpty
	}

	for i := 0; i < rv.Len(); i++ {
		if rv.Index(i).Interface() == "" {
			return ErrStringEmpty
		}

		if rv.Index(i).Interface() == 0 {
			return ErrNumZero
		}
	}

	return nil
}
