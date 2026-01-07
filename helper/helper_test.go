// Package helper 帮助函数测试
package helper

import (
	"testing"
	"time"
)

// TestDiffNano 测试时间差计算
func TestDiffNano(t *testing.T) {
	tests := []struct {
		name      string
		sleepTime time.Duration
		minDiff   int64
	}{
		{
			name:      "测试100毫秒延迟",
			sleepTime: 100 * time.Millisecond,
			minDiff:   100 * 1e6, // 100ms = 100,000,000 纳秒
		},
		{
			name:      "测试10毫秒延迟",
			sleepTime: 10 * time.Millisecond,
			minDiff:   10 * 1e6,
		},
		{
			name:      "测试零延迟",
			sleepTime: 0,
			minDiff:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startTime := time.Now()
			time.Sleep(tt.sleepTime)
			diff := DiffNano(startTime)

			if diff < tt.minDiff {
				t.Errorf("DiffNano() = %d, 期望 >= %d", diff, tt.minDiff)
			}
		})
	}
}

// TestDiffNano_Positive 测试返回值始终为正
func TestDiffNano_Positive(t *testing.T) {
	startTime := time.Now()
	diff := DiffNano(startTime)
	if diff < 0 {
		t.Errorf("DiffNano() = %d, 期望 >= 0", diff)
	}
}

// TestInArrayStr 测试字符串是否在数组内
func TestInArrayStr(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		arr      []string
		expected bool
	}{
		{
			name:     "字符串在数组中",
			str:      "hello",
			arr:      []string{"hello", "world", "test"},
			expected: true,
		},
		{
			name:     "字符串不在数组中",
			str:      "foo",
			arr:      []string{"hello", "world", "test"},
			expected: false,
		},
		{
			name:     "空数组",
			str:      "hello",
			arr:      []string{},
			expected: false,
		},
		{
			name:     "空字符串在数组中",
			str:      "",
			arr:      []string{"", "hello"},
			expected: true,
		},
		{
			name:     "空字符串不在数组中",
			str:      "",
			arr:      []string{"hello", "world"},
			expected: false,
		},
		{
			name:     "数组只有一个元素且匹配",
			str:      "only",
			arr:      []string{"only"},
			expected: true,
		},
		{
			name:     "数组只有一个元素且不匹配",
			str:      "other",
			arr:      []string{"only"},
			expected: false,
		},
		{
			name:     "区分大小写-大写",
			str:      "Hello",
			arr:      []string{"hello", "world"},
			expected: false,
		},
		{
			name:     "包含特殊字符",
			str:      "hello-world_123",
			arr:      []string{"test", "hello-world_123", "foo"},
			expected: true,
		},
		{
			name:     "包含空格",
			str:      "hello world",
			arr:      []string{"hello", "world", "hello world"},
			expected: true,
		},
		{
			name:     "nil数组",
			str:      "hello",
			arr:      nil,
			expected: false,
		},
		{
			name:     "HTTP方法GET",
			str:      "GET",
			arr:      []string{"GET", "POST", "PUT", "DELETE"},
			expected: true,
		},
		{
			name:     "HTTP方法PATCH不在列表",
			str:      "PATCH",
			arr:      []string{"GET", "POST", "PUT", "DELETE"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InArrayStr(tt.str, tt.arr)
			if result != tt.expected {
				t.Errorf("InArrayStr(%q, %v) = %v, 期望 %v", tt.str, tt.arr, result, tt.expected)
			}
		})
	}
}

// TestInArrayStr_FirstMatch 测试匹配第一个元素
func TestInArrayStr_FirstMatch(t *testing.T) {
	arr := []string{"first", "second", "third"}
	if !InArrayStr("first", arr) {
		t.Error("InArrayStr 应该能找到数组第一个元素")
	}
}

// TestInArrayStr_LastMatch 测试匹配最后一个元素
func TestInArrayStr_LastMatch(t *testing.T) {
	arr := []string{"first", "second", "third"}
	if !InArrayStr("third", arr) {
		t.Error("InArrayStr 应该能找到数组最后一个元素")
	}
}

// BenchmarkInArrayStr 性能测试
func BenchmarkInArrayStr(b *testing.B) {
	arr := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	for i := 0; i < b.N; i++ {
		InArrayStr("DELETE", arr)
	}
}

// BenchmarkDiffNano 性能测试
func BenchmarkDiffNano(b *testing.B) {
	startTime := time.Now()
	for i := 0; i < b.N; i++ {
		DiffNano(startTime)
	}
}
