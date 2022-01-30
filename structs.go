// Package checks
// If the error is "this type <%s {%s}> of field is not supported". Add processing of the type of this field.
package checks

import (
	"errors"
	"fmt"
	"reflect"
)

var ErrNotHandled = errors.New("this type is not handled")

// StructureFields recursively traverses the structure and checks that:
// - field values are not equal to an empty string or zero
// - slices and arrays are not empty.
func StructureFields(s interface{}) error {
	v := reflect.ValueOf(s)
	vType := v.Type()

	for i := 0; i < v.NumField(); i++ {
		if err := typesHandling(i, v, vType); err != nil {
			return err
		}
	}

	return nil
}

func typesHandling(i int, v reflect.Value, vType reflect.Type) error {
	switch v := reflect.ValueOf(v.Field(i).Interface()); v.Kind() {
	case reflect.Struct:
		if err := StructureFields(v.Interface()); err != nil {
			return err
		}

	case reflect.String:
		if v.String() == "" {
			return fmt.Errorf("field: %s {%s} = \"\"", vType.Field(i).Name, vType.Field(i).Tag)
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() == 0 {
			return fmt.Errorf("field: %s {%s} = %d", vType.Field(i).Name, vType.Field(i).Tag, v.Int())
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v.Uint() == 0 {
			return fmt.Errorf("field: %s {%s} = %d", vType.Field(i).Name, vType.Field(i).Tag, v.Uint())
		}

	case reflect.Slice, reflect.Array:
		if err := checkFieldsLists(v); err != nil {
			return fmt.Errorf("field: %s {%s}: %w", vType.Field(i).Name, vType.Field(i).Tag, err)
		}

	case reflect.Bool:
		break

	default:
		return fmt.Errorf("this type <%s {%s}>: %w",
			vType.Field(i).Name, vType.Field(i).Tag, ErrNotHandled)
	}

	return nil
}

func checkFieldsLists(rv reflect.Value) error {
	if rv.Len() == 0 {
		return errors.New("list is empty")
	}

	for i := 0; i < rv.Len(); i++ {
		if rv.Index(i).Interface() == "" {
			return fmt.Errorf("%s", "\"\"")
		}

		if rv.Index(i).Interface() == 0 {
			return fmt.Errorf("%v", rv.Index(i).Interface())
		}
	}

	return nil
}
