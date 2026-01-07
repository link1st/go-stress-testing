// Package verify WebSocket验证器测试
package verify

import (
	"testing"

	"github.com/link1st/go-stress-testing/model"
)

// TestWebSocketJSON 测试WebSocket JSON验证
func TestWebSocketJSON(t *testing.T) {
	tests := []struct {
		name            string
		seq             string
		msg             string
		expectedCode    int
		expectedSucceed bool
	}{
		{
			name:            "正常响应",
			seq:             "1566276523281-585638",
			msg:             `{"seq":"1566276523281-585638","cmd":"heartbeat","response":{"code":200,"codeMsg":"Success","data":null}}`,
			expectedCode:    200,
			expectedSucceed: true,
		},
		{
			name:            "seq不匹配",
			seq:             "1566276523281-585638",
			msg:             `{"seq":"different-seq","cmd":"heartbeat","response":{"code":200,"codeMsg":"Success","data":null}}`,
			expectedCode:    model.ParseError,
			expectedSucceed: false,
		},
		{
			name:            "code非200",
			seq:             "test-seq",
			msg:             `{"seq":"test-seq","cmd":"error","response":{"code":500,"codeMsg":"Error","data":null}}`,
			expectedCode:    500,
			expectedSucceed: false,
		},
		{
			name:            "无效JSON",
			seq:             "test-seq",
			msg:             `{invalid json}`,
			expectedCode:    model.ParseError,
			expectedSucceed: false,
		},
		{
			name:            "空消息",
			seq:             "test-seq",
			msg:             ``,
			expectedCode:    model.ParseError,
			expectedSucceed: false,
		},
		{
			name:            "空seq匹配",
			seq:             "",
			msg:             `{"seq":"","cmd":"test","response":{"code":200,"codeMsg":"Success","data":null}}`,
			expectedCode:    200,
			expectedSucceed: true,
		},
		{
			name:            "code为0",
			seq:             "test-seq",
			msg:             `{"seq":"test-seq","cmd":"test","response":{"code":0,"codeMsg":"Success","data":null}}`,
			expectedCode:    0,
			expectedSucceed: false,
		},
		{
			name:            "包含复杂data",
			seq:             "complex-seq",
			msg:             `{"seq":"complex-seq","cmd":"data","response":{"code":200,"codeMsg":"Success","data":{"users":[{"id":1}]}}}`,
			expectedCode:    200,
			expectedSucceed: true,
		},
		{
			name:            "401未授权",
			seq:             "auth-seq",
			msg:             `{"seq":"auth-seq","cmd":"auth","response":{"code":401,"codeMsg":"Unauthorized","data":null}}`,
			expectedCode:    401,
			expectedSucceed: false,
		},
		{
			name:            "403禁止",
			seq:             "forbidden-seq",
			msg:             `{"seq":"forbidden-seq","cmd":"access","response":{"code":403,"codeMsg":"Forbidden","data":null}}`,
			expectedCode:    403,
			expectedSucceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &model.Request{Debug: false}
			msg := []byte(tt.msg)

			code, isSucceed := WebSocketJSON(request, tt.seq, msg)

			if code != tt.expectedCode {
				t.Errorf("WebSocketJSON() code = %d, 期望 %d", code, tt.expectedCode)
			}
			if isSucceed != tt.expectedSucceed {
				t.Errorf("WebSocketJSON() isSucceed = %v, 期望 %v", isSucceed, tt.expectedSucceed)
			}
		})
	}
}

// TestWebSocketJSON_Debug 测试Debug模式
func TestWebSocketJSON_Debug(t *testing.T) {
	request := &model.Request{Debug: true}
	seq := "debug-seq"
	msg := []byte(`{"seq":"debug-seq","cmd":"test","response":{"code":200,"codeMsg":"Success","data":null}}`)

	code, isSucceed := WebSocketJSON(request, seq, msg)

	if code != 200 || !isSucceed {
		t.Errorf("Debug模式验证失败: code=%d, isSucceed=%v", code, isSucceed)
	}
}

// TestWebSocketResponseJSON_Struct 测试响应结构体
func TestWebSocketResponseJSON_Struct(t *testing.T) {
	resp := WebSocketResponseJSON{
		Seq: "test-seq",
		Cmd: "heartbeat",
	}
	resp.Response.Code = 200
	resp.Response.CodeMsg = "Success"
	resp.Response.Data = nil

	if resp.Seq != "test-seq" {
		t.Errorf("Seq = %q, 期望 'test-seq'", resp.Seq)
	}
	if resp.Cmd != "heartbeat" {
		t.Errorf("Cmd = %q, 期望 'heartbeat'", resp.Cmd)
	}
	if resp.Response.Code != 200 {
		t.Errorf("Response.Code = %d, 期望 200", resp.Response.Code)
	}
}

// TestWebSocketJSON_DifferentCmds 测试不同命令
func TestWebSocketJSON_DifferentCmds(t *testing.T) {
	cmds := []string{"heartbeat", "login", "logout", "message", "ping", "pong"}

	for _, cmd := range cmds {
		t.Run(cmd, func(t *testing.T) {
			request := &model.Request{Debug: false}
			seq := "test-" + cmd
			msg := []byte(`{"seq":"test-` + cmd + `","cmd":"` + cmd + `","response":{"code":200,"codeMsg":"Success","data":null}}`)

			code, isSucceed := WebSocketJSON(request, seq, msg)

			if code != 200 || !isSucceed {
				t.Errorf("命令 %s 验证失败", cmd)
			}
		})
	}
}

// TestWebSocketJSON_LargeMessage 测试大消息
func TestWebSocketJSON_LargeMessage(t *testing.T) {
	request := &model.Request{Debug: false}
	seq := "large-seq"

	// 构建大消息
	largeData := `{"seq":"large-seq","cmd":"data","response":{"code":200,"codeMsg":"Success","data":{"items":[`
	for i := 0; i < 100; i++ {
		if i > 0 {
			largeData += ","
		}
		largeData += `{"id":1,"name":"item","content":"some long content here"}`
	}
	largeData += `]}}}`

	code, isSucceed := WebSocketJSON(request, seq, []byte(largeData))

	if code != 200 || !isSucceed {
		t.Errorf("大消息验证失败: code=%d, isSucceed=%v", code, isSucceed)
	}
}

// TestWebSocketJSON_SpecialCharacters 测试特殊字符
func TestWebSocketJSON_SpecialCharacters(t *testing.T) {
	tests := []struct {
		name string
		seq  string
		msg  string
	}{
		{
			name: "包含中文",
			seq:  "chinese-seq",
			msg:  `{"seq":"chinese-seq","cmd":"msg","response":{"code":200,"codeMsg":"成功","data":"你好世界"}}`,
		},
		{
			name: "包含表情",
			seq:  "emoji-seq",
			msg:  `{"seq":"emoji-seq","cmd":"msg","response":{"code":200,"codeMsg":"Success","data":"Hello! 👋"}}`,
		},
		{
			name: "包含转义字符",
			seq:  "escape-seq",
			msg:  `{"seq":"escape-seq","cmd":"msg","response":{"code":200,"codeMsg":"Success","data":"line1\nline2\ttab"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &model.Request{Debug: false}
			code, isSucceed := WebSocketJSON(request, tt.seq, []byte(tt.msg))

			if code != 200 || !isSucceed {
				t.Errorf("特殊字符验证失败: code=%d, isSucceed=%v", code, isSucceed)
			}
		})
	}
}

// TestWebSocketJSON_MalformedJSON 测试各种格式错误的JSON
func TestWebSocketJSON_MalformedJSON(t *testing.T) {
	malformedJSONs := []string{
		`{`,
		`}`,
		`[]`,
		`null`,
		`"string"`,
		`123`,
		`{"seq":}`,
		`{"seq":"test","response":{"code":}}`,
		`{"seq":"test","response":{"code":"not_a_number"}}`,
	}

	for _, malformed := range malformedJSONs {
		t.Run(malformed, func(t *testing.T) {
			request := &model.Request{Debug: false}
			code, isSucceed := WebSocketJSON(request, "test", []byte(malformed))

			// 格式错误应该返回ParseError或失败
			if isSucceed {
				t.Errorf("格式错误的JSON不应该验证成功: %s", malformed)
			}
			_ = code // code可能是ParseError或0
		})
	}
}

// BenchmarkWebSocketJSON 性能测试
func BenchmarkWebSocketJSON(b *testing.B) {
	request := &model.Request{Debug: false}
	seq := "bench-seq"
	msg := []byte(`{"seq":"bench-seq","cmd":"heartbeat","response":{"code":200,"codeMsg":"Success","data":null}}`)

	for i := 0; i < b.N; i++ {
		WebSocketJSON(request, seq, msg)
	}
}

// BenchmarkWebSocketJSON_Large 大消息性能测试
func BenchmarkWebSocketJSON_Large(b *testing.B) {
	request := &model.Request{Debug: false}
	seq := "large-seq"

	largeData := `{"seq":"large-seq","cmd":"data","response":{"code":200,"codeMsg":"Success","data":{"items":[`
	for i := 0; i < 100; i++ {
		if i > 0 {
			largeData += ","
		}
		largeData += `{"id":1,"name":"item"}`
	}
	largeData += `]}}}`
	msg := []byte(largeData)

	for i := 0; i < b.N; i++ {
		WebSocketJSON(request, seq, msg)
	}
}
