// Package model CURL模型综合测试
package model

import (
	"os"
	"path/filepath"
	"testing"
)

// TestParseTheFile_Errors 测试ParseTheFile错误处理
func TestParseTheFile_Errors(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expectError bool
	}{
		{
			name:        "空路径",
			path:        "",
			expectError: true,
		},
		{
			name:        "不存在的文件",
			path:        "/nonexistent/path/file.txt",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseTheFile(tt.path)
			if tt.expectError && err == nil {
				t.Error("期望返回错误，但没有")
			}
			if !tt.expectError && err != nil {
				t.Errorf("不期望返回错误: %v", err)
			}
		})
	}
}

// TestCURL_GetURL 测试获取URL
func TestCURL_GetURL(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string][]string
		expected string
	}{
		{
			name:     "curl key",
			data:     map[string][]string{"curl": {"https://example.com"}},
			expected: "https://example.com",
		},
		{
			name:     "--url key",
			data:     map[string][]string{"--url": {"https://example.com/api"}},
			expected: "https://example.com/api",
		},
		{
			name:     "--location key",
			data:     map[string][]string{"--location": {"https://example.com/redirect"}},
			expected: "https://example.com/redirect",
		},
		{
			name:     "空数据",
			data:     map[string][]string{},
			expected: "",
		},
		{
			name:     "多个URL取第一个",
			data:     map[string][]string{"curl": {"https://first.com", "https://second.com"}},
			expected: "https://first.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CURL{Data: tt.data}
			url := c.GetURL()
			if url != tt.expected {
				t.Errorf("GetURL() = %q, 期望 %q", url, tt.expected)
			}
		})
	}
}

// TestCURL_GetMethod 测试获取HTTP方法
func TestCURL_GetMethod(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string][]string
		expected string
	}{
		{
			name:     "-X GET",
			data:     map[string][]string{"-X": {"GET"}},
			expected: "GET",
		},
		{
			name:     "-X POST",
			data:     map[string][]string{"-X": {"POST"}},
			expected: "POST",
		},
		{
			name:     "-X PUT",
			data:     map[string][]string{"-X": {"PUT"}},
			expected: "PUT",
		},
		{
			name:     "-X DELETE",
			data:     map[string][]string{"-X": {"DELETE"}},
			expected: "DELETE",
		},
		{
			name:     "--request POST",
			data:     map[string][]string{"--request": {"POST"}},
			expected: "POST",
		},
		{
			name:     "小写方法转大写",
			data:     map[string][]string{"-X": {"post"}},
			expected: "POST",
		},
		{
			name:     "无方法默认GET",
			data:     map[string][]string{},
			expected: "GET",
		},
		{
			name:     "无效方法返回默认",
			data:     map[string][]string{"-X": {"INVALID"}},
			expected: "GET",
		},
		{
			name:     "有body默认POST",
			data:     map[string][]string{"--data": {"test=data"}},
			expected: "POST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CURL{Data: tt.data}
			method := c.GetMethod()
			if method != tt.expected {
				t.Errorf("GetMethod() = %q, 期望 %q", method, tt.expected)
			}
		})
	}
}

// TestCURL_GetHeaders 测试获取Headers
func TestCURL_GetHeaders(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string][]string
		expected map[string]string
	}{
		{
			name:     "单个Header",
			data:     map[string][]string{"-H": {"Content-Type: application/json"}},
			expected: map[string]string{"Content-Type": "application/json"},
		},
		{
			name: "多个Header",
			data: map[string][]string{"-H": {
				"Content-Type: application/json",
				"Authorization: Bearer token",
			}},
			expected: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer token",
			},
		},
		{
			name:     "--header key",
			data:     map[string][]string{"--header": {"X-Custom: value"}},
			expected: map[string]string{"X-Custom": "value"},
		},
		{
			name:     "空Headers",
			data:     map[string][]string{},
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CURL{Data: tt.data}
			headers := c.GetHeaders()

			if len(headers) != len(tt.expected) {
				t.Errorf("Headers长度 = %d, 期望 %d", len(headers), len(tt.expected))
			}

			for k, v := range tt.expected {
				if headers[k] != v {
					t.Errorf("Headers[%q] = %q, 期望 %q", k, headers[k], v)
				}
			}
		})
	}
}

// TestCURL_GetBody 测试获取Body
func TestCURL_GetBody(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string][]string
		expected string
	}{
		{
			name:     "--data",
			data:     map[string][]string{"--data": {"name=test&value=123"}},
			expected: "name=test&value=123",
		},
		{
			name:     "-d",
			data:     map[string][]string{"-d": {`{"key":"value"}`}},
			expected: `{"key":"value"}`,
		},
		{
			name:     "--data-raw",
			data:     map[string][]string{"--data-raw": {"raw data"}},
			expected: "raw data",
		},
		{
			name:     "--data-urlencode",
			data:     map[string][]string{"--data-urlencode": {"encoded=data"}},
			expected: "encoded=data",
		},
		{
			name:     "--data-binary",
			data:     map[string][]string{"--data-binary": {"binary data"}},
			expected: "binary data",
		},
		{
			name:     "空body",
			data:     map[string][]string{},
			expected: "",
		},
		{
			name:     "--form",
			data:     map[string][]string{"--form": {"field1=value1", "field2=value2"}},
			expected: "field1=value1&field2=value2",
		},
		{
			name:     "-F",
			data:     map[string][]string{"-F": {"file=@/path/to/file"}},
			expected: "file=@/path/to/file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CURL{Data: tt.data}
			body := c.GetBody()
			if body != tt.expected {
				t.Errorf("GetBody() = %q, 期望 %q", body, tt.expected)
			}
		})
	}
}

// TestCURL_GetHeadersStr 测试获取Headers字符串
func TestCURL_GetHeadersStr(t *testing.T) {
	c := &CURL{Data: map[string][]string{"-H": {"Content-Type: application/json"}}}
	str := c.GetHeadersStr()

	if str == "" {
		t.Error("GetHeadersStr() 返回空字符串")
	}
	// 应该是JSON格式
	if str[0] != '{' || str[len(str)-1] != '}' {
		t.Errorf("GetHeadersStr() 应该返回JSON格式: %s", str)
	}
}

// TestCURL_String 测试String方法
func TestCURL_String(t *testing.T) {
	c := &CURL{Data: map[string][]string{"curl": {"https://example.com"}}}
	str := c.String()

	if str == "" {
		t.Error("String() 返回空字符串")
	}
}

// TestCURL_getDataValue 测试获取数据值
func TestCURL_getDataValue(t *testing.T) {
	c := &CURL{Data: map[string][]string{
		"-H":       {"header1", "header2"},
		"--header": {"header3"},
	}}

	// 测试找到第一个key
	value := c.getDataValue([]string{"-H", "--header"})
	if len(value) != 2 || value[0] != "header1" {
		t.Errorf("getDataValue() = %v, 期望从-H获取", value)
	}

	// 测试找不到任何key
	value = c.getDataValue([]string{"-X", "--request"})
	if len(value) != 0 {
		t.Errorf("getDataValue() = %v, 期望空切片", value)
	}
}

// TestArgsTrim 测试参数修剪
func TestArgsTrim(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "空数组",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "包含换行",
			input:    []string{"arg1\n", "arg2"},
			expected: []string{"arg1", "arg2"},
		},
		{
			name:     "只有换行",
			input:    []string{"\n"},
			expected: []string{""}, // TrimSpace后变成空字符串，但仍会被添加
		},
		{
			name:     "-XPOST分割",
			input:    []string{"-XPOST"},
			expected: []string{"-X", "POST"},
		},
		{
			name:     "-XGET分割",
			input:    []string{"-XGET"},
			expected: []string{"-X", "GET"},
		},
		{
			name:     "带空格",
			input:    []string{"  arg1  ", "arg2"},
			expected: []string{"arg1", "arg2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := argsTrim(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("argsTrim() 长度 = %d, 期望 %d", len(result), len(tt.expected))
				return
			}

			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("argsTrim()[%d] = %q, 期望 %q", i, v, tt.expected[i])
				}
			}
		})
	}
}

// TestRemoveSpaces 测试移除空格
func TestRemoveSpaces(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "前后空格",
			input:    "  test  ",
			expected: "test",
		},
		{
			name:     "包含反斜杠",
			input:    "\\test\\",
			expected: "test",
		},
		{
			name:     "包含换行",
			input:    "\ntest\n",
			expected: "test",
		},
		{
			name:     "混合字符",
			input:    " \\ \n test \n \\ ",
			expected: "test",
		},
		{
			name:     "空字符串",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeSpaces(tt.input)
			if result != tt.expected {
				t.Errorf("removeSpaces(%q) = %q, 期望 %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsKey 测试是否为key
func TestIsKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "短选项",
			input:    "-H",
			expected: true,
		},
		{
			name:     "长选项",
			input:    "--header",
			expected: true,
		},
		{
			name:     "curl关键字",
			input:    "curl",
			expected: true,
		},
		{
			name:     "普通值",
			input:    "value",
			expected: false,
		},
		{
			name:     "URL",
			input:    "https://example.com",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isKey(tt.input)
			if result != tt.expected {
				t.Errorf("isKey(%q) = %v, 期望 %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsURL 测试是否为URL
func TestIsURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "HTTP URL",
			input:    "http://example.com",
			expected: true,
		},
		{
			name:     "HTTPS URL",
			input:    "https://example.com",
			expected: true,
		},
		{
			name:     "普通字符串",
			input:    "example.com",
			expected: false,
		},
		{
			name:     "WS URL",
			input:    "ws://example.com",
			expected: false,
		},
		{
			name:     "空字符串",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isURL(tt.input)
			if result != tt.expected {
				t.Errorf("isURL(%q) = %v, 期望 %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestParseTheFile_ValidFile 测试解析有效文件
func TestParseTheFile_ValidFile(t *testing.T) {
	// 创建临时测试文件
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.curl.txt")

	content := `curl 'https://example.com/api' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer token' \
  --data '{"key":"value"}'`

	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	curl, err := ParseTheFile(testFile)
	if err != nil {
		t.Fatalf("ParseTheFile() 返回错误: %v", err)
	}

	// 验证URL
	url := curl.GetURL()
	if url != "https://example.com/api" {
		t.Errorf("URL = %q, 期望 'https://example.com/api'", url)
	}

	// 验证Headers
	headers := curl.GetHeaders()
	if headers["Content-Type"] != "application/json" {
		t.Errorf("Content-Type = %q, 期望 'application/json'", headers["Content-Type"])
	}

	// 验证Body
	body := curl.GetBody()
	if body != `{"key":"value"}` {
		t.Errorf("Body = %q, 期望 '{\"key\":\"value\"}'", body)
	}
}

// TestParseTheFile_GET 测试解析GET请求
func TestParseTheFile_GET(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "get.curl.txt")

	content := `curl 'https://example.com/api?page=1' \
  -H 'Accept: application/json'`

	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	curl, err := ParseTheFile(testFile)
	if err != nil {
		t.Fatalf("ParseTheFile() 返回错误: %v", err)
	}

	method := curl.GetMethod()
	if method != "GET" {
		t.Errorf("Method = %q, 期望 'GET'", method)
	}
}

// BenchmarkParseTheFile 性能测试
func BenchmarkCURL_GetURL(b *testing.B) {
	c := &CURL{Data: map[string][]string{"curl": {"https://example.com"}}}
	for i := 0; i < b.N; i++ {
		c.GetURL()
	}
}

// BenchmarkCURL_GetHeaders 性能测试
func BenchmarkCURL_GetHeaders(b *testing.B) {
	c := &CURL{Data: map[string][]string{"-H": {
		"Content-Type: application/json",
		"Authorization: Bearer token",
		"X-Custom-Header: value",
	}}}
	for i := 0; i < b.N; i++ {
		c.GetHeaders()
	}
}

// BenchmarkCURL_GetMethod 性能测试
func BenchmarkCURL_GetMethod(b *testing.B) {
	c := &CURL{Data: map[string][]string{"-X": {"POST"}}}
	for i := 0; i < b.N; i++ {
		c.GetMethod()
	}
}
