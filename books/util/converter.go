package util

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/fatih/structs"
	protostruct "github.com/golang/protobuf/ptypes/struct"
)

func Value(value *protostruct.Value) (interface{}, error) {
	var err error
	if value == nil {
		return nil, nil
	}
	if structValue, ok := value.GetKind().(*protostruct.Value_StructValue); ok {
		result := make(map[string]interface{})
		for k, v := range structValue.StructValue.Fields {
			result[k], err = Value(v)
			if err != nil {
				return nil, err
			}
		}
		return result, err
	}
	if listValue, ok := value.GetKind().(*protostruct.Value_ListValue); ok {
		result := make([]interface{}, len(listValue.ListValue.Values))
		for i, el := range listValue.ListValue.Values {
			result[i], err = Value(el)
			if err != nil {
				return nil, err
			}
		}
		return result, err
	}
	if _, ok := value.GetKind().(*protostruct.Value_NullValue); ok {
		return nil, nil
	}
	if numValue, ok := value.GetKind().(*protostruct.Value_NumberValue); ok {
		return numValue.NumberValue, nil
	}
	if strValue, ok := value.GetKind().(*protostruct.Value_StringValue); ok {
		return strValue.StringValue, nil
	}
	if boolValue, ok := value.GetKind().(*protostruct.Value_BoolValue); ok {
		return boolValue.BoolValue, nil
	}
	return fmt.Errorf(fmt.Sprintf("Cannot convert the value %+v", value)), nil
}

func Entry(entry interface{}) (*protostruct.Value, error) {
	var err error
	if entry == nil {
		return &protostruct.Value{Kind: &protostruct.Value_NullValue{}}, nil
	}
	rt := reflect.TypeOf(entry)
	switch rt.Kind() {
	case reflect.String:
		if realValue, ok := entry.(string); ok {
			return &protostruct.Value{Kind: &protostruct.Value_StringValue{StringValue: realValue}}, nil
		}
		return nil, errors.New("cannot convert string value")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &protostruct.Value{Kind: &protostruct.Value_NumberValue{NumberValue: float64(reflect.ValueOf(entry).Int())}}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return &protostruct.Value{Kind: &protostruct.Value_NumberValue{NumberValue: float64(reflect.ValueOf(entry).Uint())}}, nil
	case reflect.Float32, reflect.Float64:
		return &protostruct.Value{Kind: &protostruct.Value_NumberValue{NumberValue: reflect.ValueOf(entry).Float()}}, nil
	case reflect.Bool:
		if realValue, ok := entry.(bool); ok {
			return &protostruct.Value{Kind: &protostruct.Value_BoolValue{BoolValue: realValue}}, nil
		}
		return nil, errors.New("cannot convert boolean value")
	case reflect.Array, reflect.Slice:
		lstEntry := reflect.ValueOf(entry)

		lstValue := &protostruct.ListValue{Values: make([]*protostruct.Value, lstEntry.Len(), lstEntry.Len())}
		for i := 0; i < lstEntry.Len(); i++ {
			lstValue.Values[i], err = Entry(lstEntry.Index(i).Interface())
			if err != nil {
				return nil, err
			}
		}
		return &protostruct.Value{Kind: &protostruct.Value_ListValue{ListValue: lstValue}}, nil
	case reflect.Struct:
		return Entry(structs.Map(entry))
	case reflect.Map:
		mapEntry := make(map[string]interface{})
		entryValue := reflect.ValueOf(entry)
		for _, k := range entryValue.MapKeys() {
			mapEntry[k.String()] = entryValue.MapIndex(k).Interface()
		}
		structVal, err := MaptoStruct(mapEntry)
		return &protostruct.Value{Kind: &protostruct.Value_StructValue{StructValue: structVal}}, err
	}
	return nil, fmt.Errorf(fmt.Sprintf("Cannot convert [%+v] kind:%s", entry, rt.Kind()))
}

func MaptoStruct(input map[string]interface{}) (*protostruct.Struct, error) {
	var err error
	result := &protostruct.Struct{Fields: make(map[string]*protostruct.Value)}
	for k, v := range input {
		result.Fields[k], err = Entry(v)
		if err != nil {
			return nil, err
		}
	}
	return result, err
}
