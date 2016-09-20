package varutil

import (
	"errors"
	"fmt"
	"reflect"
)

// SetField set a value of struct
func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		invalidTypeError := errors.New("Provided value type didn't match obj field type")
		return invalidTypeError
	}

	structFieldValue.Set(val)
	return nil
}

// FillStruct load structure data from map
func FillStruct(s interface{}, m map[string]interface{}) error {
	for k, v := range m {
		err := SetField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// LoadStruct load structure data from map. Use name from tag or field name
func LoadStruct(obj interface{}, m map[string]interface{}, tagname string, ignoreUndefined bool) error {
	structValue := reflect.ValueOf(obj).Elem()
	for i := 0; i < structValue.NumField(); i++ {
		valueField := structValue.Field(i)
		structField := structValue.Type().Field(i)

		dest := structField.Tag.Get(tagname)
		if dest == "" {
			dest = structField.Name
		}
		if !valueField.IsValid() {
			return fmt.Errorf("No such field: %s in obj", structField.Name)
		}
		if !valueField.CanSet() {
			return fmt.Errorf("Cannot set %s field value", structField.Name)
		}
		importedValue, ok := m[dest]
		if !ok && !ignoreUndefined {
			return fmt.Errorf("input map[%v] is undefined", dest)
		}
		if importedValue == nil {
			if ignoreUndefined {
				continue
			}
			return fmt.Errorf("set value (named %v) can not be nil", dest)
		}
		newValue := reflect.ValueOf(importedValue)
		if structField.Type != newValue.Type() {
			invalidTypeError := errors.New("Provided value type didn't match obj field type")
			return invalidTypeError
		}
		valueField.Set(newValue)
	}
	return nil
}
