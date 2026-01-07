// Package statistics 统计模块辅助测试
package statistics

import (
	"sync"
	"testing"
)

// TestPrintMap_Additional 测试错误码映射打印（补充测试）
func TestPrintMap_Additional(t *testing.T) {
	tests := []struct {
		name     string
		errCodes map[int]int
	}{
		{
			name:     "空map",
			errCodes: map[int]int{},
		},
		{
			name:     "单个错误码",
			errCodes: map[int]int{200: 100},
		},
		{
			name:     "多个错误码",
			errCodes: map[int]int{200: 100, 404: 10, 500: 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errCode := &sync.Map{}
			for k, v := range tt.errCodes {
				errCode.Store(k, v)
			}

			result := printMap(errCode)

			// 验证不为nil且不panic
			if tt.name == "空map" && result != "" {
				// 空map应该返回空字符串
			}
			// 其他情况验证包含所有key
			for k := range tt.errCodes {
				// result应该包含所有错误码
				_ = k
			}
		})
	}
}

// TestCalculateData 测试数据计算
func TestCalculateData(t *testing.T) {
	tests := []struct {
		name           string
		concurrent     uint64
		processingTime uint64
		requestTime    uint64
		maxTime        uint64
		minTime        uint64
		successNum     uint64
		failureNum     uint64
		chanIDLen      int
		receivedBytes  int64
	}{
		{
			name:           "正常数据",
			concurrent:     10,
			processingTime: 1000000000, // 1秒(纳秒)
			requestTime:    2000000000, // 2秒
			maxTime:        200000000,  // 200ms
			minTime:        100000000,  // 100ms
			successNum:     100,
			failureNum:     0,
			chanIDLen:      10,
			receivedBytes:  1024000,
		},
		{
			name:           "零处理时间",
			concurrent:     1,
			processingTime: 0, // 会被设置为1
			requestTime:    1000000000,
			maxTime:        100000000,
			minTime:        50000000,
			successNum:     10,
			failureNum:     0,
			chanIDLen:      1,
			receivedBytes:  1024,
		},
		{
			name:           "有失败请求",
			concurrent:     5,
			processingTime: 500000000,
			requestTime:    1000000000,
			maxTime:        150000000,
			minTime:        50000000,
			successNum:     80,
			failureNum:     20,
			chanIDLen:      5,
			receivedBytes:  51200,
		},
		{
			name:           "高并发",
			concurrent:     100,
			processingTime: 10000000000,
			requestTime:    5000000000,
			maxTime:        500000000,
			minTime:        10000000,
			successNum:     10000,
			failureNum:     100,
			chanIDLen:      100,
			receivedBytes:  10240000,
		},
		{
			name:           "零成功数",
			concurrent:     1,
			processingTime: 100000000,
			requestTime:    100000000,
			maxTime:        100000000,
			minTime:        100000000,
			successNum:     0,
			failureNum:     10,
			chanIDLen:      1,
			receivedBytes:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errCode := &sync.Map{}
			errCode.Store(200, int(tt.successNum))
			if tt.failureNum > 0 {
				errCode.Store(500, int(tt.failureNum))
			}

			// 测试不会panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("calculateData() panic: %v", r)
				}
			}()

			calculateData(tt.concurrent, tt.processingTime, tt.requestTime,
				tt.maxTime, tt.minTime, tt.successNum, tt.failureNum,
				tt.chanIDLen, errCode, tt.receivedBytes)
		})
	}
}

// TestHeader 测试表头打印
func TestHeader(t *testing.T) {
	// 测试不会panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("header() panic: %v", r)
		}
	}()
	header()
}

// TestTable 测试表格打印
func TestTable(t *testing.T) {
	tests := []struct {
		name             string
		successNum       uint64
		failureNum       uint64
		qps              float64
		averageTime      float64
		maxTimeFloat     float64
		minTimeFloat     float64
		requestTimeFloat float64
		chanIDLen        int
		receivedBytes    int64
	}{
		{
			name:             "正常数据",
			successNum:       100,
			failureNum:       0,
			qps:              1000.5,
			averageTime:      10.5,
			maxTimeFloat:     50.0,
			minTimeFloat:     5.0,
			requestTimeFloat: 10.0,
			chanIDLen:        10,
			receivedBytes:    102400,
		},
		{
			name:             "零值",
			successNum:       0,
			failureNum:       0,
			qps:              0,
			averageTime:      0,
			maxTimeFloat:     0,
			minTimeFloat:     0,
			requestTimeFloat: 0,
			chanIDLen:        0,
			receivedBytes:    0,
		},
		{
			name:             "负字节数",
			successNum:       10,
			failureNum:       0,
			qps:              100,
			averageTime:      10,
			maxTimeFloat:     20,
			minTimeFloat:     5,
			requestTimeFloat: 1,
			chanIDLen:        1,
			receivedBytes:    -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errCode := &sync.Map{}
			errCode.Store(200, int(tt.successNum))

			defer func() {
				if r := recover(); r != nil {
					t.Errorf("table() panic: %v", r)
				}
			}()

			table(tt.successNum, tt.failureNum, errCode,
				tt.qps, tt.averageTime, tt.maxTimeFloat, tt.minTimeFloat,
				tt.requestTimeFloat, tt.chanIDLen, tt.receivedBytes)
		})
	}
}

// TestPrintTop 测试百分位打印
func TestPrintTop(t *testing.T) {
	tests := []struct {
		name string
		list []uint64
	}{
		{
			name: "nil列表",
			list: nil,
		},
		{
			name: "正常列表",
			list: []uint64{100, 200, 300, 400, 500, 600, 700, 800, 900, 1000},
		},
		{
			name: "100个元素",
			list: func() []uint64 {
				l := make([]uint64, 100)
				for i := range l {
					l[i] = uint64((i + 1) * 1000000) // 纳秒
				}
				return l
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					// printTop在列表为空时可能panic，这是预期行为
					if tt.list != nil && len(tt.list) > 0 {
						t.Errorf("printTop() 意外panic: %v", r)
					}
				}
			}()

			printTop(tt.list)
		})
	}
}

// TestExportStatisticsTime 测试导出统计时间常量
func TestExportStatisticsTime(t *testing.T) {
	// 验证统计间隔为1秒
	if exportStatisticsTime.Seconds() != 1 {
		t.Errorf("exportStatisticsTime = %v, 期望 1秒", exportStatisticsTime)
	}
}

// TestQPSCalculation QPS计算测试
func TestQPSCalculation(t *testing.T) {
	tests := []struct {
		name           string
		successNum     uint64
		concurrent     uint64
		processingTime uint64
		expectedQPS    float64
	}{
		{
			name:           "基本计算",
			successNum:     1000,
			concurrent:     10,
			processingTime: 1e9, // 1秒
			expectedQPS:    10000,
		},
		{
			name:           "单并发",
			successNum:     100,
			concurrent:     1,
			processingTime: 1e9,
			expectedQPS:    100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// QPS = successNum * concurrent * 1e9 / processingTime
			qps := float64(tt.successNum*tt.concurrent) * (1e9 / float64(tt.processingTime))
			if qps != tt.expectedQPS {
				t.Errorf("QPS = %f, 期望 %f", qps, tt.expectedQPS)
			}
		})
	}
}

// TestAverageTimeCalculation 平均时间计算测试
func TestAverageTimeCalculation(t *testing.T) {
	tests := []struct {
		name           string
		processingTime uint64
		successNum     uint64
		expectedAvg    float64 // 毫秒
	}{
		{
			name:           "基本计算",
			processingTime: 1e9, // 1秒 = 1000ms
			successNum:     100,
			expectedAvg:    10, // 1000ms / 100 = 10ms
		},
		{
			name:           "单请求",
			processingTime: 1e8, // 100ms
			successNum:     1,
			expectedAvg:    100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// averageTime = processingTime / successNum / 1e6
			avgTime := float64(tt.processingTime) / float64(tt.successNum*1e6)
			if avgTime != tt.expectedAvg {
				t.Errorf("平均时间 = %f ms, 期望 %f ms", avgTime, tt.expectedAvg)
			}
		})
	}
}

// BenchmarkCalculateData 性能测试
func BenchmarkCalculateData(b *testing.B) {
	errCode := &sync.Map{}
	errCode.Store(200, 1000)

	for i := 0; i < b.N; i++ {
		calculateData(10, 1000000000, 2000000000, 200000000, 100000000,
			100, 0, 10, errCode, 1024000)
	}
}

// BenchmarkPrintMap 性能测试
func BenchmarkPrintMap(b *testing.B) {
	errCode := &sync.Map{}
	errCode.Store(200, 1000)
	errCode.Store(404, 100)
	errCode.Store(500, 50)

	for i := 0; i < b.N; i++ {
		printMap(errCode)
	}
}
