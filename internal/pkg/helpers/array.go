package helpers

import (
	"reflect"
)

func defaultEqualFn(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func mustArrayValue(val reflect.Value) {
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		panic("type must be array")
	}
}

func mustPointerValue(val reflect.Value) {
	if val.Kind() != reflect.Ptr {
		panic("type must be ptr")
	}
}

// InArray
// 判断一个key是否存在于数组中
func InArray(needle interface{}, haystack interface{}) bool {
	return InArrayFunc(haystack, func(i interface{}) bool {
		return defaultEqualFn(i, needle)
	})
}

func InArrayFunc(haystack interface{}, f func(interface{}) bool) bool {
	val := reflect.ValueOf(haystack)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if f(val.Index(i).Interface()) {
				return true
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if f(val.MapIndex(k).Interface()) {
				return true
			}
		}
	default:
		panic("haystack: haystack type muset be slice, array or map")
	}
	return false
}

// ArrayColumn
// 1. input = []map/struct 返回map/struct 对应的 []map
// 2. input = map 返回　map[key] 对应的 []map
func ArrayColumn(input interface{}, key string) (out []map[string]interface{}) {
	val := reflect.ValueOf(input)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			out = append(out, ArrayColumn(val.Index(i).Interface(), key)...)
		}
		return
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if defaultEqualFn(k.Interface(), key) {
				out = append(out, map[string]interface{}{
					key: val.MapIndex(k).Interface(),
				})
			}
		}
	case reflect.Struct:
		structKeyVal := val.FieldByName(key)
		if structKeyVal.IsValid() {
			out = append(out, map[string]interface{}{key: structKeyVal.Interface()})
		}
	case reflect.Ptr:
		return ArrayColumn(reflect.Indirect(val).Interface(), key)
	}
	return
}

// ArrayIntersect 两个数组取交集
// 取数组的交集，会使用第一个匹配到的交集构造slice结构
// NOTICE:
// 1. 只支持数组
// 2. 如果数组中有类型不一致的会忽略掉
// 3. 如果没有取到交集数组，将会返回nil
func ArrayIntersect(array1 interface{}, array2 interface{}, arrays ...interface{}) interface{} {
	var vals []reflect.Value
	vals = append(vals, reflect.ValueOf(array1))
	vals = append(vals, reflect.ValueOf(array2))
	for _, arr := range arrays {
		vals = append(vals, reflect.ValueOf(arr))
	}

	for _, val := range vals {
		mustArrayValue(val)
	}

	var (
		out       reflect.Value
		firstType reflect.Type
	)

	val0 := vals[0]
	for i := 0; i < val0.Len(); i++ {
		allIn := true
		for j := 1; j < len(vals); j++ {
			if !InArray(val0.Index(i).Interface(), vals[j].Interface()) {
				allIn = false
				break
			}
		}
		if allIn {
			if !out.IsValid() {
				firstType = val0.Index(i).Type()
				out = reflect.MakeSlice(reflect.SliceOf(firstType), 0, val0.Len()-i)
			}
			if firstType == val0.Index(i).Type() {
				out = reflect.Append(out, val0.Index(i))
			}
		}
	}

	if out.IsValid() {
		return out.Interface()
	}
	return nil
}

// ArrayIntersectString 数组取交集　返回[]string
func ArrayIntersectString(array1 interface{}, array2 interface{}, arrays ...interface{}) []string {
	return ArrayIntersect(array1, array2, arrays...).([]string)
}

// ArrayUnique 数组去重
func ArrayUnique(array interface{}) {
	ptrVal := reflect.ValueOf(array)
	mustPointerValue(ptrVal)

	val, preVal := ptrVal, ptrVal
out:
	for {
		switch val.Kind() {
		case reflect.Ptr:
			preVal = val
			fallthrough
		case reflect.Interface:
			val = val.Elem()
		default:
			break out
		}
	}
	mustArrayValue(val)

	for i := 0; i < val.Len(); i++ {
		for j := i + 1; j < val.Len(); {
			if defaultEqualFn(val.Index(i).Interface(), val.Index(j).Interface()) {
				val = reflect.AppendSlice(val.Slice(0, j), val.Slice(j+1, val.Len()))
				continue
			}
			j++
		}
	}
	preVal.Elem().Set(val)
}

// ArrayStringUnique 字符串数组去重 返回[]string
func ArrayStringUnique(l1 []string) []string {
	var l2 []string
	m := make(map[string]struct{})
	for i := range l1 {
		if _, ok := m[l1[i]]; !ok {
			m[l1[i]] = struct{}{}
			l2 = append(l2, l1[i])
		}
	}
	return l2
}

// ArrayUnionString 数组取并集集　返回[]string
func ArrayUnionString(array1, array2 []string) []string {
	m := make(map[string]int)
	for _, v := range array1 {
		m[v]++
	}

	for _, v := range array2 {
		times := m[v]
		if times == 0 {
			array1 = append(array1, v)
		}
	}
	return array1
}

// ArrayUnion 取并集
// NOTICE:
// 1. 只支持数组
// 2. 如果数组中有类型不一致的会忽略掉
func ArrayUnion(array1 interface{}, array2 interface{}) interface{} {

	val1 := reflect.ValueOf(array1)
	val2 := reflect.ValueOf(array2)

	mustArrayValue(val1)
	mustArrayValue(val2)

	var (
		out       reflect.Value
		firstType reflect.Type
	)

	if val1.IsNil() || val1.Len() == 0 {
		return val2.Interface()
	} else if val2.IsNil() || val2.Len() == 0 {
		return val1.Interface()
	}

	firstType = val1.Index(0).Type()
	out = reflect.MakeSlice(reflect.SliceOf(firstType), 0, val1.Len())

	for i := 0; i < val1.Len(); i++ {
		out = reflect.Append(out, val1.Index(i))
	}

	for j := 0; j < val2.Len(); j++ {
		if !InArray(val2.Index(j).Interface(), val1.Interface()) {
			if firstType == val2.Index(j).Type() {
				out = reflect.Append(out, val2.Index(j))
			}
		}
	}
	if out.IsValid() {
		return out.Interface()
	}
	return nil
}

// 求s1在全集s2中的补集
func ArrayComplement(s1, s2 interface{}) interface{} {
	sub := reflect.ValueOf(s1)
	full := reflect.ValueOf(s2)

	mustArrayValue(sub)
	mustArrayValue(full)

	var (
		out       reflect.Value
		firstType reflect.Type
	)

	// 如果sub是空集，那么full都是它的补集
	if sub.IsNil() || sub.Len() == 0 {
		return full.Interface()
	}
	// 如果full是空，或者长度为0,就返回自己
	if full.IsNil() || full.Len() == 0 {
		return full.Interface()
	}

	firstType = full.Index(0).Type()
	out = reflect.MakeSlice(reflect.SliceOf(firstType), 0, full.Cap())
	for i := 0; i < full.Len(); i++ {
		if !InArray(full.Index(i).Interface(), sub.Interface()) {
			out = reflect.Append(out, full.Index(i))
		}
	}

	if out.IsValid() {
		return out.Interface()
	}
	return nil
}
