package Binary

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"reflect"
)

func Serialize(object interface{}) ([]byte, error) {
	buffer := new(bytes.Buffer)
	objecttype := reflect.ValueOf(object)

	for i := 0; i < objecttype.NumField(); i++ {
		field := objecttype.Field(i)

		switch field.Kind() {
		case reflect.String:
			for i := 0; i < field.Len(); i++ {
				err := binary.Write(buffer, binary.LittleEndian, field.Index(i).Interface())

				if err != nil {
					return nil, fmt.Errorf("Error converting with binary.\n%s", err)
				}
			}
			// binary.Write(buffer, binary.LittleEndian, 0)
			buffer.Write([]byte{0})
		case reflect.Slice:
			for i := 0; i < field.Len(); i++ {
				data, err := Serialize(field.Index(i).Interface())
				if err != nil {
					return nil, fmt.Errorf("Error converting with binary.\n%s", err)
				}
				buffer.Write(data)
			}
			buffer.Write([]byte{0})
		case reflect.Struct:
			data, err := Serialize(field.Interface())
			if err != nil {
				return nil, fmt.Errorf("Error converting with binary.\n%s", err)
			}
			buffer.Write(data)
		default:
			err := binary.Write(buffer, binary.LittleEndian, field.Interface())
			if err != nil {
				return nil, fmt.Errorf("Error converting with binary.\n%s", err)
			}
		}
	}

	fmt.Println(buffer.Bytes())
	fmt.Println(hex.EncodeToString(buffer.Bytes()))
	return buffer.Bytes(), nil
}

func Deserialize(bytes []byte) interface{} {
	return nil
}
