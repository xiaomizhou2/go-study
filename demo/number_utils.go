package demo

import (
	"fmt"
)

// FilterEven 从数字列表中找出所有偶数
// 对比 Java: public static List<Integer> filterEven(List<Integer> numbers) throws IllegalArgumentException
func FilterEven(numbers []int) ([]int, error) {
	if len(numbers) == 0 {
		return nil, fmt.Errorf("数字列表不能为空")
	}

	evens := make([]int, 0)
	for _, num := range numbers {
		if num%2 == 0 {
			evens = append(evens, num)
		}
	}

	return evens, nil
}
