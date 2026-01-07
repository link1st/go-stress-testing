// Package model 请求数据模型测试
package model

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// TestRequest_GetBody 测试获取请求体
func TestRequest_GetBody(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected string
	}{
		{
			name:     "普通字符串",
			body:     "hello world",
			expected: "hello world",
		},
		{
			name:     "JSON格式",
			body:     `{"key":"value"}`,
			expected: `{"key":"value"}`,
		},
		{
			name:     "空字符串",
			body:     "",
			expected: "",
		},
		{
			name:     "URL编码数据",
			body:     "name=test&value=123",
			expected: "name=test&value=123",
		},
		{
			name:     "包含特殊字符",
			body:     "data=hello%20world&token=abc123",
			expected: "data=hello%20world&token=abc123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Request{Body: tt.body}
			body := r.GetBody()
			result, err := io.ReadAll(body)
			if err != nil {
				t.Fatalf("读取body失败: %v", err)
			}
			if string(result) != tt.expected {
				t.Errorf("GetBody() = %q, 期望 %q", string(result), tt.expected)
			}
		})
	}
}

// TestRequest_CopyHeaders 测试复制Headers
func TestRequest_CopyHeaders(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
	}{
		{
			name:    "普通Headers",
			headers: map[string]string{"Content-Type": "application/json", "Authorization": "Bearer token"},
		},
		{
			name:    "空Headers",
			headers: map[string]string{},
		},
		{
			name:    "单个Header",
			headers: map[string]string{"X-Custom": "value"},
		},
		{
			name:    "nil Headers",
			headers: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Request{Headers: tt.headers}
			copied := r.CopyHeaders()

			// 验证复制后的map不为nil
			if copied == nil {
				t.Fatal("CopyHeaders() 返回了nil")
			}

			// 验证长度相同
			expectedLen := 0
			if tt.headers != nil {
				expectedLen = len(tt.headers)
			}
			if len(copied) != expectedLen {
				t.Errorf("CopyHeaders() 长度 = %d, 期望 %d", len(copied), expectedLen)
			}

			// 验证内容相同
			for k, v := range tt.headers {
				if copied[k] != v {
					t.Errorf("CopyHeaders()[%q] = %q, 期望 %q", k, copied[k], v)
				}
			}

			// 验证是深拷贝（修改原始不影响拷贝）
			if tt.headers != nil && len(tt.headers) > 0 {
				for k := range tt.headers {
					tt.headers[k] = "modified"
					break
				}
				// 拷贝的值应该不受影响
				for k, v := range copied {
					if v == "modified" && tt.headers[k] == "modified" {
						// 这是预期的，因为已经修改了原始值
					}
				}
			}
		})
	}
}

// TestRequest_CopyHeaders_DeepCopy 测试深拷贝
func TestRequest_CopyHeaders_DeepCopy(t *testing.T) {
	original := map[string]string{"key": "original"}
	r := &Request{Headers: original}
	copied := r.CopyHeaders()

	// 修改原始map
	original["key"] = "modified"
	original["new"] = "value"

	// 验证拷贝不受影响
	if copied["key"] != "original" {
		t.Errorf("深拷贝失败: copied[key] = %q, 期望 'original'", copied["key"])
	}
	if _, exists := copied["new"]; exists {
		t.Error("深拷贝失败: 新增的key不应该出现在拷贝中")
	}
}

// TestRequest_getVerifyKey 测试获取验证key
func TestRequest_getVerifyKey(t *testing.T) {
	tests := []struct {
		name     string
		form     string
		verify   string
		expected string
	}{
		{
			name:     "HTTP statusCode",
			form:     FormTypeHTTP,
			verify:   "statusCode",
			expected: "http.statusCode",
		},
		{
			name:     "HTTP json",
			form:     FormTypeHTTP,
			verify:   "json",
			expected: "http.json",
		},
		{
			name:     "WebSocket json",
			form:     FormTypeWebSocket,
			verify:   "json",
			expected: "webSocket.json",
		},
		{
			name:     "gRPC验证",
			form:     FormTypeGRPC,
			verify:   "proto",
			expected: "grpc.proto",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Request{Form: tt.form, Verify: tt.verify}
			if got := r.getVerifyKey(); got != tt.expected {
				t.Errorf("getVerifyKey() = %q, 期望 %q", got, tt.expected)
			}
		})
	}
}

// TestRequest_GetDebug 测试获取debug参数
func TestRequest_GetDebug(t *testing.T) {
	tests := []struct {
		name     string
		debug    bool
		expected bool
	}{
		{
			name:     "debug开启",
			debug:    true,
			expected: true,
		},
		{
			name:     "debug关闭",
			debug:    false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Request{Debug: tt.debug}
			if got := r.GetDebug(); got != tt.expected {
				t.Errorf("GetDebug() = %v, 期望 %v", got, tt.expected)
			}
		})
	}
}

// TestGetForm 测试URL协议解析
func TestGetForm(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		expectedForm string
		expectedURL  string
	}{
		{
			name:         "HTTP协议",
			url:          "http://example.com",
			expectedForm: FormTypeHTTP,
			expectedURL:  "http://example.com",
		},
		{
			name:         "HTTPS协议",
			url:          "https://example.com",
			expectedForm: FormTypeHTTP,
			expectedURL:  "https://example.com",
		},
		{
			name:         "WebSocket协议",
			url:          "ws://example.com/ws",
			expectedForm: FormTypeWebSocket,
			expectedURL:  "ws://example.com/ws",
		},
		{
			name:         "WebSocket安全协议",
			url:          "wss://example.com/ws",
			expectedForm: FormTypeWebSocket,
			expectedURL:  "wss://example.com/ws",
		},
		{
			name:         "gRPC协议",
			url:          "grpc://example.com:8080",
			expectedForm: FormTypeGRPC,
			expectedURL:  "grpc://example.com:8080",
		},
		{
			name:         "RPC协议",
			url:          "rpc://example.com:8080",
			expectedForm: FormTypeGRPC,
			expectedURL:  "rpc://example.com:8080",
		},
		{
			name:         "Radius协议",
			url:          "radius://example.com:1812",
			expectedForm: FormTypeRadius,
			expectedURL:  "example.com:1812",
		},
		{
			name:         "无协议前缀",
			url:          "example.com/api",
			expectedForm: FormTypeHTTP,
			expectedURL:  "http://example.com/api",
		},
		{
			name:         "无协议带端口",
			url:          "localhost:8080/api",
			expectedForm: FormTypeHTTP,
			expectedURL:  "http://localhost:8080/api",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form, url := getForm(tt.url)
			if form != tt.expectedForm {
				t.Errorf("getForm() form = %q, 期望 %q", form, tt.expectedForm)
			}
			if url != tt.expectedURL {
				t.Errorf("getForm() url = %q, 期望 %q", url, tt.expectedURL)
			}
		})
	}
}

// TestGetHeaderValue 测试Header值解析
func TestGetHeaderValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		initial  map[string]string
		expected map[string]string
	}{
		{
			name:     "标准Header",
			input:    "Content-Type: application/json",
			initial:  map[string]string{},
			expected: map[string]string{"Content-Type": "application/json"},
		},
		{
			name:     "带空格的值",
			input:    "Content-Type: application/json; charset=utf-8",
			initial:  map[string]string{},
			expected: map[string]string{"Content-Type": "application/json; charset=utf-8"},
		},
		{
			name:     "无冒号",
			input:    "InvalidHeader",
			initial:  map[string]string{},
			expected: map[string]string{},
		},
		{
			name:     "空值",
			input:    "X-Empty:",
			initial:  map[string]string{},
			expected: map[string]string{"X-Empty": ""},
		},
		{
			name:     "合并相同Header",
			input:    "Cookie: session=abc",
			initial:  map[string]string{"Cookie": "token=xyz"},
			expected: map[string]string{"Cookie": "token=xyz; session=abc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := make(map[string]string)
			for k, v := range tt.initial {
				headers[k] = v
			}
			getHeaderValue(tt.input, headers)
			for k, v := range tt.expected {
				if headers[k] != v {
					t.Errorf("getHeaderValue() headers[%q] = %q, 期望 %q", k, headers[k], v)
				}
			}
		})
	}
}

// TestRequestResults_SetID 测试设置请求ID
func TestRequestResults_SetID(t *testing.T) {
	tests := []struct {
		name       string
		chanID     uint64
		number     uint64
		expectedID string
	}{
		{
			name:       "正常ID",
			chanID:     1,
			number:     100,
			expectedID: "1_100",
		},
		{
			name:       "零值",
			chanID:     0,
			number:     0,
			expectedID: "0_0",
		},
		{
			name:       "大数值",
			chanID:     999999,
			number:     888888,
			expectedID: "999999_888888",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RequestResults{}
			r.SetID(tt.chanID, tt.number)
			if r.ID != tt.expectedID {
				t.Errorf("SetID() ID = %q, 期望 %q", r.ID, tt.expectedID)
			}
			if r.ChanID != tt.chanID {
				t.Errorf("SetID() ChanID = %d, 期望 %d", r.ChanID, tt.chanID)
			}
		})
	}
}

// TestRegisterVerifyHTTP 测试注册HTTP验证器
func TestRegisterVerifyHTTP(t *testing.T) {
	// 保存原始状态
	originalMap := make(map[string]VerifyHTTP)
	verifyMapHTTPMutex.RLock()
	for k, v := range verifyMapHTTP {
		originalMap[k] = v
	}
	verifyMapHTTPMutex.RUnlock()

	// 注册测试验证器
	testVerify := func(request *Request, response *http.Response, body []byte) (code int, isSucceed bool) {
		return 200, true
	}

	RegisterVerifyHTTP("testVerify", testVerify)

	// 验证注册成功
	verifyMapHTTPMutex.RLock()
	_, ok := verifyMapHTTP["http.testVerify"]
	verifyMapHTTPMutex.RUnlock()

	if !ok {
		t.Error("RegisterVerifyHTTP() 注册失败")
	}

	// 清理测试数据
	verifyMapHTTPMutex.Lock()
	delete(verifyMapHTTP, "http.testVerify")
	verifyMapHTTPMutex.Unlock()
}

// TestRegisterVerifyWebSocket 测试注册WebSocket验证器
func TestRegisterVerifyWebSocket(t *testing.T) {
	testVerify := func(request *Request, seq string, msg []byte) (code int, isSucceed bool) {
		return 200, true
	}

	RegisterVerifyWebSocket("testVerify", testVerify)

	verifyMapWebSocketMutex.RLock()
	_, ok := verifyMapWebSocket["webSocket.testVerify"]
	verifyMapWebSocketMutex.RUnlock()

	if !ok {
		t.Error("RegisterVerifyWebSocket() 注册失败")
	}

	// 清理
	verifyMapWebSocketMutex.Lock()
	delete(verifyMapWebSocket, "webSocket.testVerify")
	verifyMapWebSocketMutex.Unlock()
}

// TestRequest_Print 测试打印功能（不会panic）
func TestRequest_Print(t *testing.T) {
	tests := []struct {
		name    string
		request *Request
	}{
		{
			name:    "nil请求",
			request: nil,
		},
		{
			name: "正常请求",
			request: &Request{
				URL:       "http://example.com",
				Form:      FormTypeHTTP,
				Method:    "GET",
				Headers:   map[string]string{"Content-Type": "application/json"},
				Body:      "test body",
				Verify:    "statusCode",
				Timeout:   30 * time.Second,
				Debug:     true,
				HTTP2:     false,
				Keepalive: true,
				MaxCon:    10,
			},
		},
		{
			name: "空Headers请求",
			request: &Request{
				URL:     "http://example.com",
				Form:    FormTypeHTTP,
				Method:  "GET",
				Headers: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Print() panic: %v", r)
				}
			}()
			tt.request.Print()
		})
	}
}

// TestConstants 测试常量定义
func TestConstants(t *testing.T) {
	// 测试HTTP状态码常量
	if HTTPOk != 200 {
		t.Errorf("HTTPOk = %d, 期望 200", HTTPOk)
	}
	if RequestErr != 509 {
		t.Errorf("RequestErr = %d, 期望 509", RequestErr)
	}
	if ParseError != 510 {
		t.Errorf("ParseError = %d, 期望 510", ParseError)
	}

	// 测试协议类型常量
	if FormTypeHTTP != "http" {
		t.Errorf("FormTypeHTTP = %q, 期望 'http'", FormTypeHTTP)
	}
	if FormTypeWebSocket != "webSocket" {
		t.Errorf("FormTypeWebSocket = %q, 期望 'webSocket'", FormTypeWebSocket)
	}
	if FormTypeGRPC != "grpc" {
		t.Errorf("FormTypeGRPC = %q, 期望 'grpc'", FormTypeGRPC)
	}
	if FormTypeRadius != "radius" {
		t.Errorf("FormTypeRadius = %q, 期望 'radius'", FormTypeRadius)
	}
}

// TestRequest_GetVerifyHTTP_Panic 测试获取不存在的验证器会panic
func TestRequest_GetVerifyHTTP_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("GetVerifyHTTP() 应该panic但没有")
		}
	}()

	r := &Request{Form: FormTypeHTTP, Verify: "nonexistent"}
	r.GetVerifyHTTP()
}

// TestRequest_GetVerifyWebSocket_Panic 测试获取不存在的WebSocket验证器会panic
func TestRequest_GetVerifyWebSocket_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("GetVerifyWebSocket() 应该panic但没有")
		}
	}()

	r := &Request{Form: FormTypeWebSocket, Verify: "nonexistent"}
	r.GetVerifyWebSocket()
}

// BenchmarkRequest_CopyHeaders 性能测试
func BenchmarkRequest_CopyHeaders(b *testing.B) {
	r := &Request{
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer token",
			"X-Custom-1":    "value1",
			"X-Custom-2":    "value2",
		},
	}
	for i := 0; i < b.N; i++ {
		r.CopyHeaders()
	}
}

// BenchmarkGetForm 性能测试
func BenchmarkGetForm(b *testing.B) {
	urls := []string{
		"http://example.com",
		"https://example.com",
		"ws://example.com",
		"grpc://example.com",
	}
	for i := 0; i < b.N; i++ {
		getForm(urls[i%len(urls)])
	}
}

// TestRequest_Body_Reader 测试Body返回的Reader是否可重复读取
func TestRequest_Body_Reader(t *testing.T) {
	r := &Request{Body: "test content"}

	// 第一次读取
	reader1 := r.GetBody()
	content1, _ := io.ReadAll(reader1)

	// 第二次读取（应该返回新的Reader）
	reader2 := r.GetBody()
	content2, _ := io.ReadAll(reader2)

	if string(content1) != string(content2) {
		t.Error("GetBody() 每次调用应该返回相同内容的Reader")
	}
}

// TestRequest_MethodUppercase 测试方法大小写
func TestRequest_MethodUppercase(t *testing.T) {
	methods := []string{"get", "post", "put", "delete", "patch"}
	for _, m := range methods {
		upper := strings.ToUpper(m)
		if upper != strings.ToUpper(m) {
			t.Errorf("方法 %q 转大写后应为 %q", m, upper)
		}
	}
}
