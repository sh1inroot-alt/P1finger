package sliceopt

import (
	"strconv"
)

// 切片去重 使用泛型来去除切片中的重复元素 go > 1.18
func SliceRmDuplication[T comparable](old []T) []T {
	r := make([]T, 0, len(old))
	checkMap := make(map[T]struct{})

	for _, item := range old {
		if _, ok := checkMap[item]; !ok {
			checkMap[item] = struct{}{}
			r = append(r, item)
		}
	}
	return r
}

func FilterUniqueElements(original, filter []string) []string {
	// 使用map来存储过滤集合中的元素，以快速查找
	filterSet := make(map[string]struct{})
	for _, item := range filter {
		filterSet[item] = struct{}{}
	}

	// 存储去重后的结果
	var unique []string

	// 遍历原始切片，检查元素是否在过滤集合中
	for _, item := range original {
		if _, exists := filterSet[item]; !exists {
			unique = append(unique, item)
		}
	}
	return unique
}

func StringSliceToIntSlice(stringSlice []string) ([]int, error) {
	intSlice := make([]int, len(stringSlice))
	for i, str := range stringSlice {
		intValue, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		intSlice[i] = intValue
	}
	return intSlice, nil
}

// 定义泛型函数，T是一个类型参数
func ElementExists[T comparable](list []T, element T) bool {
	for _, item := range list {
		if item == element {
			return true
		}
	}
	return false
}
