package types

import "strings"

type Values = map[string]any

func Get[T any](v Values, def T, path string) T {
	if v == nil || len(path) == 0 {
		return def
	}
	pathArr := strings.Split(path, ".")
	value, ok := v[pathArr[0]]
	if !ok {
		return def
	}
	if len(pathArr) == 1 {
		typedValue, ok := value.(T)
		if ok {
			return typedValue
		} else {
			return def
		}
	}

	subMap, ok := value.(map[string]any)
	if !ok {
		return def
	}
	return Get[T](subMap, def, strings.Join(pathArr[1:], "."))
}

func Exists(v Values, path string) bool {
	if v == nil || len(path) == 0 {
		return false
	}
	pathArr := strings.Split(path, ".")
	value, ok := v[pathArr[0]]
	if !ok {
		return false
	}
	if len(pathArr) == 1 {
		return true
	}

	subMap, ok := value.(map[string]any)
	if !ok {
		return false
	}
	return Exists(subMap, strings.Join(pathArr[1:], "."))
}

func MergeValues(v1, v2 Values) Values {
	result := make(Values, len(v1))
	for k, v := range v1 {
		result[k] = v
	}
	for k, v := range v2 {
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := result[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					result[k] = MergeValues(bv, v)
					continue
				}
			}
		}
		result[k] = v
	}
	return result
}
