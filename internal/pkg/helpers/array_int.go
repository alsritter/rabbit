package helpers

// ArrayDiffInt32 两个数组的差集
func ArrayDiffInt32(array1 []int32, array2 []int32) (diff []int32) {
	has1 := make(map[int32]bool, len(array1))
	for _, k := range array1 {
		has1[k] = true
	}
	has2 := make(map[int32]bool, len(array2))
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

// rmDuplicateIntArray 数组去重
func RmDuplicateIntArray(intArr []int32) []int32 {
	intMap := make(map[int32]bool)
	newIntArr := make([]int32, 0)
	for _, value := range intArr {
		if _, ok := intMap[value]; !ok {
			intMap[value] = true
			newIntArr = append(newIntArr, value)
		}
	}
	return newIntArr
}
