package dto_utils

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func ConvertSlice[S any, D any](source []S) []D {
	var dest []D
	for _, srcItem := range source {
		var destItem D
		jsonData, _ := json.Marshal(srcItem)
		jsonUnmarshalError := json.Unmarshal(jsonData, &destItem)
		if jsonUnmarshalError != nil {
			panic(jsonUnmarshalError)
		}
		//err := mapstructure.Decode(srcItem, &destItem)
		//if err != nil {
		//	panic(err)
		//}
		dest = append(dest, destItem)
	}
	return dest
}

// DtoConvertErrorHandled marshals source 'S' into JSON, then unmarshals it into destination 'D'
func DtoConvertErrorHandled[S any, D any](source S, dest *D) {
	// Marshal source into JSON
	jsonData, err := json.Marshal(source)
	if err != nil {
		panic(err)
	}

	// Unmarshal into destination
	err = json.Unmarshal(jsonData, dest)
	if err != nil {
		panic(err)
	}
}

// DtoConvertErrorHandled marshals source 'S' into JSON, then unmarshals it into destination 'D'
func DtoConvertErrorHandledReturnError[S any, D any](source S, dest *D) error {
	// Marshal source into JSON
	jsonData, err := json.Marshal(source)
	if err != nil {
		return fmt.Errorf("error marshaling source: %w", err)
	}

	// Unmarshal into destination
	err = json.Unmarshal(jsonData, dest)
	if err != nil {
		return fmt.Errorf("error unmarshaling into destination: %w", err)
	}

	return nil
}

func NullAwareMapDtoConvertor(src interface{}, dest interface{}) {
	srcValue := reflect.ValueOf(src)
	destValue := reflect.ValueOf(dest).Elem()

	if srcValue.Kind() != reflect.Struct || destValue.Kind() != reflect.Struct {
		panic("src and dest must be structs")
	}

	for i := 0; i < srcValue.NumField(); i++ {
		srcField := srcValue.Field(i)
		destField := destValue.FieldByName(srcValue.Type().Field(i).Name)

		if !destField.IsValid() {
			continue
		}

		if srcField.Kind() == reflect.Ptr {
			if srcField.IsNil() {
				// If srcField is nil, do not update destField
				continue
			}
		}

		if srcField.CanSet() {
			if destField.Kind() == reflect.Ptr {
				// Create a new pointer if srcField is non-nil
				if !srcField.IsNil() {
					destField.Set(reflect.New(srcField.Type().Elem()))
					destField.Elem().Set(srcField.Elem())
				}
			} else {
				destField.Set(srcField)
			}
		}
	}
}
