package demo

import (
	"fmt"
	"testing"
)

func TestFilterEven(t *testing.T) {
	tests := []struct {
		name    string
		input   []int
		want    []int
		wantErr bool
	}{
		{"正常过滤", []int{1, 2, 3, 4, 5, 6, 8, 10}, []int{2, 4, 6, 8, 10}, false},
		{"没有偶数", []int{1, 3, 5, 7}, []int{}, false},
		{"全是偶数", []int{2, 4, 6}, []int{2, 4, 6}, false},
		{"空列表", []int{}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FilterEven(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilterEven() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !equal(got, tt.want) {
				t.Errorf("FilterEven() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 演示：测试文件里也可以有 main 函数的效果
func ExampleFilterEven() {
	result, _ := FilterEven([]int{1, 2, 3, 4, 5, 6, 8, 10})
	fmt.Println("偶数:", result)
	// Output: 偶数: [2 4 6 8 10]
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
