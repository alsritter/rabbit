package helpers

// InArrayString 字符串 needle 是否在 数组中
func InArrayString(needle string, haystack []string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}

// ArrayDiffString 两个数组的差集
func ArrayDiffString(array1 []string, array2 []string) (diff []string) {
	has1 := make(map[string]bool, len(array1))
	for _, k := range array1 {
		has1[k] = true
	}
	has2 := make(map[string]bool, len(array2))
	for _, k := range array2 {
		has2[k] = true
	}
	for _, k := range array1 {
		if !has2[k] {
			diff = append(diff, k)
		}
	}
	for _, k := range array2 {
		if !has1[k] {
			diff = append(diff, k)
		}
	}
	return
}
