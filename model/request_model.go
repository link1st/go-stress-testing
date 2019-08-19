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
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	HttpOk         = 200 // 请求成功
	RequestTimeout = 506 // 请求超时
	RequestErr     = 509 // 请求错误
	ParseError     = 510 // 解析错误

	FormTypeHttp      = "http"
	FormTypeWebSocket = "webSocket"
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
	Url        string            // Url
	Form       string            // http/webSocket/tcp
	Method     string            // 方法 get/post/put
	Headers    map[string]string // Headers
	Body       io.Reader         // body
	Verify     string            // 验证的方法
	VerifyHttp VerifyHttp        // 验证的方法
	Timeout    time.Duration     // 请求超时时间
	Debug      bool              // 是否开启Debug模式

	// 连接以后初始化事件
	// 循环事件 切片 时间 动作
}

func NewRequest(url string, verify string, timeout time.Duration, debug bool, path string) (request *Request, err error) {

	var (
		method  = "GET"
		headers = make(map[string]string)
		body    io.Reader
	)

	if path != "" {
		curl, err := ParseTheFile(path)
		if err != nil {

			return nil, err
		}

		if url == "" {
			url = curl.GetUrl()
		}

		method = curl.GetMethod()
		headers = curl.GetHeaders()
		body = curl.GetBody()
	}

	form := ""
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		form = FormTypeHttp
	} else if strings.HasPrefix(url, "ws://") || strings.HasPrefix(url, "wss://") {
		form = FormTypeWebSocket
	}

	if form == "" {
		err = errors.New(fmt.Sprintf("url:%s 不合法,必须是完整http、webSocket连接", url))

		return
	}

	// verify
	if verify == "" {
		verify = "statusCode"
	}

	key := fmt.Sprintf("%s.%s", form, verify)
	verifyHttp, ok := verifyMap[key]
	if !ok {
		err = errors.New("验证器不存在:" + key)

		return
	}

	if timeout == 0 {
		timeout = 3 * time.Second
	}

	request = &Request{
		Url:        url,
		Form:       form,
		Method:     strings.ToUpper(method),
		Headers:    headers,
		Body:       body,
		Verify:     verify,
		VerifyHttp: verifyHttp,
		Timeout:    timeout,
		Debug:      debug,
	}

	return

}

// 打印
func (r *Request) Print() {
	if r == nil {

		return
	}

	result := fmt.Sprintf("request:\n url:%s \n form:%s \n method:%s \n headers:%v \n", r.Url, r.Form, r.Method, r.Headers)
	result = fmt.Sprintf("%s verify:%s \n timeout:%s \n debu:%v \n", result, r.Verify, r.Timeout, r.Debug)
	fmt.Println(result)

	return
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
