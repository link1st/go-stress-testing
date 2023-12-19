// Package verify 校验
package verify

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/link1st/go-stress-testing/model"
)

// HTTPStatusCode 通过 HTTP 状态码判断是否请求成功
func HTTPStatusCode(request *model.Request, response *http.Response, body []byte) (code int, isSucceed bool) {
	code = response.StatusCode
	if code == request.Code {
		isSucceed = true
	}
	// 开启调试模式
	if request.GetDebug() {
		fmt.Printf("请求结果 httpCode:%d body:%s \n", response.StatusCode, string(body))
	}
	return
}

/***************************  返回值为json  ********************************/

// ResponseJSON 返回数据结构体
type ResponseJSON struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// HTTPJson  通过返回的Body 判断
// 返回示例: {"code":200,"msg":"Success","data":{}}
// code 默认将http code作为返回码，http code 为200时 取body中的返回code
func HTTPJson(request *model.Request, response *http.Response, body []byte) (code int, isSucceed bool) {
	code = response.StatusCode
	if code == http.StatusOK {
		responseJSON := &ResponseJSON{}
		if err := json.Unmarshal(body, responseJSON); err != nil {
			code = model.ParseError
			fmt.Printf("请求结果 json.Unmarshal err:%v", err)
		} else {
			code = responseJSON.Code
			// body 中code返回200为返回数据成功
			if responseJSON.Code == request.Code {
				isSucceed = true
			}
		}
		// 开启调试模式
		if request.GetDebug() {
			fmt.Printf("请求结果 httpCode:%d body:%s  \n", response.StatusCode, string(body))
		}
	}
	return
}
