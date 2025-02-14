package tg

import (
	"reflect"
	"strconv"
	"strings"
)

func defaults[T any](val T) T {
	return defaultsInternal(val).(T)
}

func defaultsInternal(val any) any {
	reflectVal := reflect.ValueOf(val)
	for reflectVal.Kind() == reflect.Ptr {
		if reflectVal.IsNil() {
			return val
		}
		reflectVal = reflectVal.Elem()
	}

	switch reflectVal.Kind() {
	case reflect.Struct:
		reflectType := reflectVal.Type()
		for i := 0; i < reflectType.NumField(); i++ {
			field := reflectVal.Field(i)

			switch field.Kind() {
			case reflect.Struct, reflect.Array, reflect.Slice, reflect.Map:
				field.Set(reflect.Indirect(reflect.ValueOf(defaultsInternal(field.Addr().Interface()))))
			case reflect.Ptr, reflect.Interface:
				if !field.IsNil() {
					field.Set(reflect.ValueOf(defaultsInternal(field.Interface())))
				}
			default:
				defaultTag := reflectType.Field(i).Tag.Get("default")
				if defaultTag != "" && field.IsZero() {
					setDefaultValue(field, defaultTag)
				}
			}
		}

	case reflect.Array, reflect.Slice:
		for i := 0; i < reflectVal.Len(); i++ {
			value := reflectVal.Index(i)
			switch value.Kind() {
			case reflect.Struct, reflect.Array, reflect.Slice, reflect.Map:
				value.Set(reflect.Indirect(reflect.ValueOf(defaultsInternal(value.Addr().Interface()))))
			case reflect.Ptr, reflect.Interface:
				value.Set(reflect.ValueOf(defaultsInternal(value.Interface())))
			}
		}

	case reflect.Map:
		for _, key := range reflectVal.MapKeys() {
			value := reflectVal.MapIndex(key)

			switch value.Kind() {
			case /*reflect.Struct, you can not take addr of map[string]struct */ reflect.Array, reflect.Slice, reflect.Map:
				reflectVal.
					SetMapIndex(key, reflect.Indirect(reflect.ValueOf(defaultsInternal(value.Addr().Interface()))))
			case reflect.Ptr, reflect.Interface:
				reflectVal.
					SetMapIndex(key, reflect.ValueOf(defaultsInternal(value.Interface())))
			}
		}
	}

	return val
}

func setDefaultValue(field reflect.Value, defaultValue string) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(defaultValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, err := strconv.ParseInt(defaultValue, 10, 64)
		if err != nil {
			panic("bad default: " + err.Error())
		}
		field.SetInt(value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, err := strconv.ParseUint(defaultValue, 10, 64)
		if err != nil {
			panic("bad default: " + err.Error())
		}
		field.SetUint(value)
	case reflect.Float32, reflect.Float64:
		value, err := strconv.ParseFloat(defaultValue, 64)
		if err != nil {
			panic("bad default: " + err.Error())
		}
		field.SetFloat(value)
	case reflect.Bool:
		field.SetBool(strings.ToLower(defaultValue) == "true")
	}
}
