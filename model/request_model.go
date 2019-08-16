/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-15
* Time: 18:19
 */

package model

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	HttpOk         = 200 // 请求成功
	RequestTimeout = 506 // 请求超时
	RequestErr     = 509 // 请求错误
	ParseError     = 510 // 解析错误
)

var (
	verifyMap = map[string]VerifyHttp{
		"http.statusCode": HttpStatusCode,
		"http.json":       HttpJson,
	}
)

// 验证器
type Verify interface {
	GetCode() int    // 有一个方法，返回code为200为成功
	GetResult() bool // 返回是否成功
}

// 200为成功
type VerifyHttp func(request *Request, response *http.Response) (code int, isSucceed bool)

// 请求结果
type Request struct {
	Url        string     // Url
	Form       string     // http/webSocket/tcp
	Method     string     // 方法 get/post/put
	Verify     string     // 验证的方法
	VerifyHttp VerifyHttp // 验证的方法
	Timeout    uint32     // 请求超时时间 秒
	Debug      bool       // 是否开启Debug模式
}

func (r *Request) GetDebug() bool {

	return r.Debug
}

func (r *Request) IsParameterLegal() (err error) {

	r.Form = "http"
	// statusCode json
	r.Verify = "json"

	key := fmt.Sprintf("%s.%s", r.Form, r.Verify)
	value, ok := verifyMap[key]
	if !ok {

		return errors.New("验证器不存在:" + key)
	}

	r.VerifyHttp = value

	return
}

// 请求结果
type RequestResults struct {
	Id        string // 消息Id
	Time      uint64 // 请求时间 纳秒
	IsSucceed bool   // 是否请求成功
	ErrCode   int    // 错误码
}

func (r *RequestResults) SetId(chanId uint64, number uint64) {
	id := fmt.Sprintf("%d_%d", chanId, number)

	r.Id = id
}
