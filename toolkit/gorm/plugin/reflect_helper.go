package plugin

import "reflect"

// WalkStruct 遍历 ReflectValue 中的所有 struct
func WalkStruct(rv reflect.Value, fn func(v reflect.Value)) {
	if !rv.IsValid() {
		return
	}
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Struct:
		fn(rv)
	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			WalkStruct(rv.Index(i), fn)
		}
	}
}
