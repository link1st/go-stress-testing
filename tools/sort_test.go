// Package tools 工具包测试
package tools

import (
	"sort"
	"testing"
)

// TestUint64List_Len 测试长度方法
func TestUint64List_Len(t *testing.T) {
	tests := []struct {
		name     string
		list     Uint64List
		expected int
	}{
		{
			name:     "空列表",
			list:     Uint64List{},
			expected: 0,
		},
		{
			name:     "单元素列表",
			list:     Uint64List{100},
			expected: 1,
		},
		{
			name:     "多元素列表",
			list:     Uint64List{1, 2, 3, 4, 5},
			expected: 5,
		},
		{
			name:     "nil列表",
			list:     nil,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.list.Len(); got != tt.expected {
				t.Errorf("Uint64List.Len() = %v, 期望 %v", got, tt.expected)
			}
		})
	}
}

// TestUint64List_Swap 测试交换方法
func TestUint64List_Swap(t *testing.T) {
	tests := []struct {
		name     string
		list     Uint64List
		i, j     int
		expected Uint64List
	}{
		{
			name:     "交换首尾元素",
			list:     Uint64List{1, 2, 3},
			i:        0,
			j:        2,
			expected: Uint64List{3, 2, 1},
		},
		{
			name:     "交换相邻元素",
			list:     Uint64List{1, 2, 3},
			i:        0,
			j:        1,
			expected: Uint64List{2, 1, 3},
		},
		{
			name:     "交换相同位置",
			list:     Uint64List{1, 2, 3},
			i:        1,
			j:        1,
			expected: Uint64List{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.list.Swap(tt.i, tt.j)
			for k, v := range tt.list {
				if v != tt.expected[k] {
					t.Errorf("Swap后索引%d: got %v, 期望 %v", k, v, tt.expected[k])
				}
			}
		})
	}
}

// TestUint64List_Less 测试比较方法
func TestUint64List_Less(t *testing.T) {
	tests := []struct {
		name     string
		list     Uint64List
		i, j     int
		expected bool
	}{
		{
			name:     "前小后大",
			list:     Uint64List{1, 2, 3},
			i:        0,
			j:        1,
			expected: true,
		},
		{
			name:     "前大后小",
			list:     Uint64List{3, 2, 1},
			i:        0,
			j:        1,
			expected: false,
		},
		{
			name:     "相等元素",
			list:     Uint64List{5, 5, 5},
			i:        0,
			j:        1,
			expected: false,
		},
		{
			name:     "零值比较",
			list:     Uint64List{0, 1},
			i:        0,
			j:        1,
			expected: true,
		},
		{
			name:     "大数比较",
			list:     Uint64List{18446744073709551615, 18446744073709551614},
			i:        1,
			j:        0,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.list.Less(tt.i, tt.j); got != tt.expected {
				t.Errorf("Uint64List.Less(%d, %d) = %v, 期望 %v", tt.i, tt.j, got, tt.expected)
			}
		})
	}
}

// TestUint64List_Sort 测试排序功能
func TestUint64List_Sort(t *testing.T) {
	tests := []struct {
		name     string
		list     Uint64List
		expected Uint64List
	}{
		{
			name:     "已排序列表",
			list:     Uint64List{1, 2, 3, 4, 5},
			expected: Uint64List{1, 2, 3, 4, 5},
		},
		{
			name:     "逆序列表",
			list:     Uint64List{5, 4, 3, 2, 1},
			expected: Uint64List{1, 2, 3, 4, 5},
		},
		{
			name:     "随机顺序",
			list:     Uint64List{3, 1, 4, 1, 5, 9, 2, 6},
			expected: Uint64List{1, 1, 2, 3, 4, 5, 6, 9},
		},
		{
			name:     "单元素",
			list:     Uint64List{42},
			expected: Uint64List{42},
		},
		{
			name:     "空列表",
			list:     Uint64List{},
			expected: Uint64List{},
		},
		{
			name:     "重复元素",
			list:     Uint64List{5, 5, 5, 5},
			expected: Uint64List{5, 5, 5, 5},
		},
		{
			name:     "包含零",
			list:     Uint64List{3, 0, 2, 0, 1},
			expected: Uint64List{0, 0, 1, 2, 3},
		},
		{
			name:     "纳秒时间模拟",
			list:     Uint64List{1000000, 500000, 2000000, 100000},
			expected: Uint64List{100000, 500000, 1000000, 2000000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sort.Sort(tt.list)
			if len(tt.list) != len(tt.expected) {
				t.Errorf("排序后长度不匹配: got %d, 期望 %d", len(tt.list), len(tt.expected))
				return
			}
			for i, v := range tt.list {
				if v != tt.expected[i] {
					t.Errorf("排序后索引%d: got %v, 期望 %v", i, v, tt.expected[i])
				}
			}
		})
	}
}

// TestUint64List_SortInterface 确认实现了sort.Interface
func TestUint64List_SortInterface(t *testing.T) {
	var _ sort.Interface = Uint64List{}
}

// TestUint64List_Percentile 测试百分位计算场景
func TestUint64List_Percentile(t *testing.T) {
	// 模拟响应时间数据（纳秒）
	list := Uint64List{
		100000000,  // 100ms
		150000000,  // 150ms
		200000000,  // 200ms
		120000000,  // 120ms
		180000000,  // 180ms
		250000000,  // 250ms
		90000000,   // 90ms
		300000000,  // 300ms
		110000000,  // 110ms
		160000000,  // 160ms
	}

	sort.Sort(list)

	// 验证排序后的顺序
	for i := 1; i < len(list); i++ {
		if list[i-1] > list[i] {
			t.Errorf("排序错误: list[%d]=%d > list[%d]=%d", i-1, list[i-1], i, list[i])
		}
	}

	// 验证第一个元素是最小值
	if list[0] != 90000000 {
		t.Errorf("最小值错误: got %d, 期望 90000000", list[0])
	}

	// 验证最后一个元素是最大值
	if list[len(list)-1] != 300000000 {
		t.Errorf("最大值错误: got %d, 期望 300000000", list[len(list)-1])
	}
}

// BenchmarkUint64List_Sort 排序性能测试
func BenchmarkUint64List_Sort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		list := Uint64List{
			1000000, 500000, 2000000, 100000, 1500000,
			800000, 300000, 900000, 1200000, 700000,
		}
		sort.Sort(list)
	}
}

// BenchmarkUint64List_Sort_Large 大数据量排序性能测试
func BenchmarkUint64List_Sort_Large(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		list := make(Uint64List, 10000)
		for j := range list {
			list[j] = uint64(10000 - j)
		}
		b.StartTimer()
		sort.Sort(list)
	}
}
