package Reflection

import (
	"reflect"
)

// Returns the fields of the object as an array.
func GetArrayOfFields(object interface{}) []interface{} {
	value := reflect.ValueOf(object).Elem()
	// numfields := value.NumField()
	// fmt.Printf("Value type: %s, number of fields: %d\n", value.Type().Name(), numfields)

	fields := make([]interface{}, value.NumField())

	for i := 0; i < value.NumField(); i++ {
		fieldtype := value.Field(i).Type()

		fields[i] = reflect.New(fieldtype).Interface()
	}

	return fields
}

// This expects the input in fields to be pointers
// TODO: Update to handle conversion between pointers and doublecheck the types?
func CreateResultFromFields(fields []interface{}, result interface{}) interface{} {
	value := reflect.ValueOf(result).Elem()

	for i := 0; i < value.NumField(); i++ {
		fieldvalue := reflect.ValueOf(fields[i]).Elem()

		// if (fieldvalue.Kind() == reflect.Ptr && value.Field(i).Kind() == reflect.Ptr) ||
		// 	(fieldvalue.Kind() == reflect.Interface && value.Field(i).Kind() == reflect.Interface) {
		// 	value.Field(i).Set(fieldvalue)
		// }

		value.Field(i).Set(fieldvalue)
	}

	return result
}
