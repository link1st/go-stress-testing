// Package verify HTTP验证器测试
package verify

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/link1st/go-stress-testing/model"
)

// mockResponseBody 创建mock响应体
func mockResponseBody(body string) io.ReadCloser {
	return io.NopCloser(bytes.NewBufferString(body))
}

// TestHTTPStatusCode 测试HTTP状态码验证
func TestHTTPStatusCode(t *testing.T) {
	tests := []struct {
		name            string
		requestCode     int
		responseCode    int
		expectedCode    int
		expectedSucceed bool
	}{
		{
			name:            "200成功",
			requestCode:     200,
			responseCode:    200,
			expectedCode:    200,
			expectedSucceed: true,
		},
		{
			name:            "状态码不匹配",
			requestCode:     200,
			responseCode:    404,
			expectedCode:    404,
			expectedSucceed: false,
		},
		{
			name:            "500错误",
			requestCode:     200,
			responseCode:    500,
			expectedCode:    500,
			expectedSucceed: false,
		},
		{
			name:            "201创建成功",
			requestCode:     201,
			responseCode:    201,
			expectedCode:    201,
			expectedSucceed: true,
		},
		{
			name:            "204无内容",
			requestCode:     204,
			responseCode:    204,
			expectedCode:    204,
			expectedSucceed: true,
		},
		{
			name:            "301重定向",
			requestCode:     301,
			responseCode:    301,
			expectedCode:    301,
			expectedSucceed: true,
		},
		{
			name:            "期望200得到201",
			requestCode:     200,
			responseCode:    201,
			expectedCode:    201,
			expectedSucceed: false,
		},
		{
			name:            "400错误请求",
			requestCode:     200,
			responseCode:    400,
			expectedCode:    400,
			expectedSucceed: false,
		},
		{
			name:            "401未授权",
			requestCode:     200,
			responseCode:    401,
			expectedCode:    401,
			expectedSucceed: false,
		},
		{
			name:            "403禁止访问",
			requestCode:     200,
			responseCode:    403,
			expectedCode:    403,
			expectedSucceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &model.Request{
				Code:  tt.requestCode,
				Debug: false,
			}
			response := &http.Response{
				StatusCode: tt.responseCode,
				Body:       mockResponseBody("test body"),
			}
			body := []byte("test body")

			code, isSucceed := HTTPStatusCode(request, response, body)

			if code != tt.expectedCode {
				t.Errorf("HTTPStatusCode() code = %d, 期望 %d", code, tt.expectedCode)
			}
			if isSucceed != tt.expectedSucceed {
				t.Errorf("HTTPStatusCode() isSucceed = %v, 期望 %v", isSucceed, tt.expectedSucceed)
			}
		})
	}
}

// TestHTTPStatusCode_Debug 测试Debug模式
func TestHTTPStatusCode_Debug(t *testing.T) {
	request := &model.Request{
		Code:  200,
		Debug: true,
	}
	response := &http.Response{
		StatusCode: 200,
		Body:       mockResponseBody("debug test"),
	}
	body := []byte("debug test")

	// 不应该panic
	code, isSucceed := HTTPStatusCode(request, response, body)
	if code != 200 || !isSucceed {
		t.Errorf("Debug模式下验证失败")
	}
}

// TestHTTPJson 测试JSON验证
func TestHTTPJson(t *testing.T) {
	tests := []struct {
		name            string
		requestCode     int
		responseCode    int
		body            string
		expectedCode    int
		expectedSucceed bool
	}{
		{
			name:            "JSON code 200成功",
			requestCode:     200,
			responseCode:    200,
			body:            `{"code":200,"msg":"Success","data":{}}`,
			expectedCode:    200,
			expectedSucceed: true,
		},
		{
			name:            "JSON code 不匹配",
			requestCode:     200,
			responseCode:    200,
			body:            `{"code":500,"msg":"Error","data":{}}`,
			expectedCode:    500,
			expectedSucceed: false,
		},
		{
			name:            "HTTP状态码非200",
			requestCode:     200,
			responseCode:    404,
			body:            `{"code":200,"msg":"Success","data":{}}`,
			expectedCode:    404,
			expectedSucceed: false,
		},
		{
			name:            "无效JSON",
			requestCode:     200,
			responseCode:    200,
			body:            `{invalid json}`,
			expectedCode:    model.ParseError,
			expectedSucceed: false,
		},
		{
			name:            "空JSON对象",
			requestCode:     0,
			responseCode:    200,
			body:            `{}`,
			expectedCode:    0,
			expectedSucceed: true,
		},
		{
			name:            "JSON code为0成功",
			requestCode:     0,
			responseCode:    200,
			body:            `{"code":0,"msg":"Success","data":{}}`,
			expectedCode:    0,
			expectedSucceed: true,
		},
		{
			name:            "JSON包含data",
			requestCode:     200,
			responseCode:    200,
			body:            `{"code":200,"msg":"Success","data":{"id":1,"name":"test"}}`,
			expectedCode:    200,
			expectedSucceed: true,
		},
		{
			name:            "空body",
			requestCode:     200,
			responseCode:    200,
			body:            ``,
			expectedCode:    model.ParseError,
			expectedSucceed: false,
		},
		{
			name:            "纯文本响应",
			requestCode:     200,
			responseCode:    200,
			body:            `Hello World`,
			expectedCode:    model.ParseError,
			expectedSucceed: false,
		},
		{
			name:            "JSON数组",
			requestCode:     200,
			responseCode:    200,
			body:            `[1,2,3]`,
			expectedCode:    model.ParseError, // JSON数组无法解析到ResponseJSON结构体
			expectedSucceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &model.Request{
				Code:  tt.requestCode,
				Debug: false,
			}
			response := &http.Response{
				StatusCode: tt.responseCode,
				Body:       mockResponseBody(tt.body),
			}
			body := []byte(tt.body)

			code, isSucceed := HTTPJson(request, response, body)

			if code != tt.expectedCode {
				t.Errorf("HTTPJson() code = %d, 期望 %d", code, tt.expectedCode)
			}
			if isSucceed != tt.expectedSucceed {
				t.Errorf("HTTPJson() isSucceed = %v, 期望 %v", isSucceed, tt.expectedSucceed)
			}
		})
	}
}

// TestHTTPJson_Debug 测试JSON验证Debug模式
func TestHTTPJson_Debug(t *testing.T) {
	request := &model.Request{
		Code:  200,
		Debug: true,
	}
	response := &http.Response{
		StatusCode: 200,
		Body:       mockResponseBody(`{"code":200,"msg":"Success","data":{}}`),
	}
	body := []byte(`{"code":200,"msg":"Success","data":{}}`)

	code, isSucceed := HTTPJson(request, response, body)
	if code != 200 || !isSucceed {
		t.Errorf("Debug模式下JSON验证失败")
	}
}

// TestResponseJSON_Struct 测试ResponseJSON结构体
func TestResponseJSON_Struct(t *testing.T) {
	resp := ResponseJSON{
		Code: 200,
		Msg:  "Success",
		Data: map[string]interface{}{"key": "value"},
	}

	if resp.Code != 200 {
		t.Errorf("ResponseJSON.Code = %d, 期望 200", resp.Code)
	}
	if resp.Msg != "Success" {
		t.Errorf("ResponseJSON.Msg = %q, 期望 'Success'", resp.Msg)
	}
	if resp.Data == nil {
		t.Error("ResponseJSON.Data 不应为 nil")
	}
}

// TestHTTPStatusCode_AllCodes 测试所有常见状态码
func TestHTTPStatusCode_AllCodes(t *testing.T) {
	codes := []int{
		100, 101, // 信息响应
		200, 201, 202, 204, // 成功响应
		301, 302, 304, 307, 308, // 重定向
		400, 401, 403, 404, 405, 408, 429, // 客户端错误
		500, 502, 503, 504, // 服务器错误
	}

	for _, statusCode := range codes {
		t.Run(http.StatusText(statusCode), func(t *testing.T) {
			request := &model.Request{
				Code:  statusCode,
				Debug: false,
			}
			response := &http.Response{
				StatusCode: statusCode,
				Body:       mockResponseBody(""),
			}

			code, isSucceed := HTTPStatusCode(request, response, []byte(""))

			if code != statusCode {
				t.Errorf("code = %d, 期望 %d", code, statusCode)
			}
			if !isSucceed {
				t.Errorf("期望状态码 %d 验证成功", statusCode)
			}
		})
	}
}

// TestHTTPJson_ComplexData 测试复杂JSON数据
func TestHTTPJson_ComplexData(t *testing.T) {
	complexJSON := `{
		"code": 200,
		"msg": "Success",
		"data": {
			"users": [
				{"id": 1, "name": "Alice"},
				{"id": 2, "name": "Bob"}
			],
			"total": 2,
			"page": 1
		}
	}`

	request := &model.Request{Code: 200, Debug: false}
	response := &http.Response{
		StatusCode: 200,
		Body:       mockResponseBody(complexJSON),
	}

	code, isSucceed := HTTPJson(request, response, []byte(complexJSON))

	if code != 200 || !isSucceed {
		t.Errorf("复杂JSON验证失败: code=%d, isSucceed=%v", code, isSucceed)
	}
}

// BenchmarkHTTPStatusCode 性能测试
func BenchmarkHTTPStatusCode(b *testing.B) {
	request := &model.Request{Code: 200, Debug: false}
	response := &http.Response{
		StatusCode: 200,
		Body:       mockResponseBody("test"),
	}
	body := []byte("test")

	for i := 0; i < b.N; i++ {
		HTTPStatusCode(request, response, body)
	}
}

// BenchmarkHTTPJson 性能测试
func BenchmarkHTTPJson(b *testing.B) {
	request := &model.Request{Code: 200, Debug: false}
	body := []byte(`{"code":200,"msg":"Success","data":{}}`)
	response := &http.Response{
		StatusCode: 200,
		Body:       mockResponseBody(string(body)),
	}

	for i := 0; i < b.N; i++ {
		HTTPJson(request, response, body)
	}
}

// BenchmarkHTTPJson_Large 大JSON性能测试
func BenchmarkHTTPJson_Large(b *testing.B) {
	request := &model.Request{Code: 200, Debug: false}
	largeJSON := `{"code":200,"msg":"Success","data":{"items":[`
	for i := 0; i < 100; i++ {
		if i > 0 {
			largeJSON += ","
		}
		largeJSON += `{"id":1,"name":"item","value":"data"}`
	}
	largeJSON += `]}}`
	body := []byte(largeJSON)
	response := &http.Response{
		StatusCode: 200,
		Body:       mockResponseBody(largeJSON),
	}

	for i := 0; i < b.N; i++ {
		HTTPJson(request, response, body)
	}
}
