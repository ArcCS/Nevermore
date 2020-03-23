package utils

import (
	"log"
	"reflect"
)


func PopulateStruct(popStruct interface{}, popValues interface{}){
	/// Determine the type of struct so that we can populate appropriately
	erv := reflect.ValueOf(popStruct).Elem()

	var subType reflect.Type
	var subStruct, valStruct interface{}

	for i := 0; i < erv.NumField(); i++ {
		if erv.Field(i).Kind() == reflect.Struct {
			subType = erv.Type().Field(i).Type
			subStruct = reflect.New(subType).Elem().Interface()
			valStruct = GetField(popValues, erv.Type().Field(i).Name)
			PopulateStruct(&subStruct, &valStruct)
			SetField(popStruct, erv.Type().Field(i).Name, subStruct)
		}else{
			SetField(popStruct, erv.Type().Field(i).Name, GetField(popValues, erv.Type().Field(i).Name))
		}
	}
}

// setField sets field of v with given name to given value.
func SetField(v interface{}, name string, value interface{}){
	// v must be a pointer to a struct
	rv := reflect.ValueOf(v).Elem()
	rv.FieldByName(name).Set(reflect.ValueOf(value))

}

// Get field gets a dynamic value name
func GetField(v interface{}, name string) interface{} {
	// v must be a pointer to a struct
	rv := reflect.ValueOf(v).Elem()

	// Lookup field by name
	fv := rv.FieldByName(name)
	if !fv.IsValid() {
		log.Printf("Panic on getField failed on %s field %s", "test", name)
		return nil
	}

	return fv.Interface()
}