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

func MergeValues(values ...Values) Values {
	if len(values) == 0 {
		return nil
	}
	result := make(Values, len(values[0]))
	for k, v := range values[0] {
		result[k] = v
	}
	for i := 1; i < len(values); i++ {
		for k, v := range values[i] {
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
	}
	return result
}
