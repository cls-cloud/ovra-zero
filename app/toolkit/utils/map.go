package utils

import (
	"reflect"
	"time"
)

// StructToMap 将结构体转为 map，并支持白名单或黑名单字段过滤
func StructToMap[T any](input T, includeFields []string, excludeFields []string) map[string]interface{} {
	includeSet := make(map[string]struct{})
	excludeSet := make(map[string]struct{})

	for _, f := range includeFields {
		includeSet[f] = struct{}{}
	}
	for _, f := range excludeFields {
		excludeSet[f] = struct{}{}
	}

	result := make(map[string]interface{})
	val := reflect.ValueOf(input)
	typ := reflect.TypeOf(input)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" { // 跳过未导出字段
			continue
		}

		fieldName := field.Name

		// 黑名单优先
		if _, skip := excludeSet[fieldName]; skip {
			continue
		}
		// 如果有白名单，只保留白名单字段
		if len(includeSet) > 0 {
			if _, keep := includeSet[fieldName]; !keep {
				continue
			}
		}

		result[fieldName] = val.Field(i).Interface()
	}

	return result
}

// StructToMapOmit [T any]
//
//	@Description: 将结构体转为 map
//	@param input 输入
//	@param includeFields 需要保留的字段
//	@param excludeFields 需要移除的字段
//	@param omitEmpty true 忽略零值数据
//	@return map[string]interface{} 返回的map
func StructToMapOmit[T any](input T, includeFields []string, excludeFields []string, omitEmpty bool) map[string]interface{} {
	includeSet := make(map[string]struct{})
	excludeSet := make(map[string]struct{})

	for _, f := range includeFields {
		includeSet[f] = struct{}{}
	}
	for _, f := range excludeFields {
		excludeSet[f] = struct{}{}
	}

	result := make(map[string]interface{})
	val := reflect.ValueOf(input)
	typ := reflect.TypeOf(input)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" { // 跳过未导出字段
			continue
		}

		fieldName := field.Name

		// 黑名单优先
		if _, skip := excludeSet[fieldName]; skip {
			continue
		}
		// 如果设置了白名单，只保留白名单字段
		if len(includeSet) > 0 {
			if _, keep := includeSet[fieldName]; !keep {
				continue
			}
		}

		value := val.Field(i)

		// 如果启用了 omitEmpty，跳过零值字段
		if omitEmpty && isEmptyValue(value) {
			continue
		}

		result[fieldName] = value.Interface()
	}

	return result
}

func isEmptyValue(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface:
		if v.IsNil() {
			return true
		}
		return isEmptyValue(v.Elem())

	case reflect.String, reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0

	case reflect.Float32, reflect.Float64:
		return v.Float() == 0

	case reflect.Struct:
		if v.Type().AssignableTo(reflect.TypeOf(time.Time{})) {
			t := v.Interface().(time.Time)
			if t.IsZero() {
				return true
			}
			// 明确判断 MySQL 异常时间字符串
			if t.Format(time.DateTime) == "0000-00-00 00:00:00" {
				return true
			}
		}
		return false
	default:
		return false
	}
}
