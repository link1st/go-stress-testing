/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-15
* Time: 18:19
 */

package model

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	HttpOk         = 200
	RequestTimeout = 506
	RequestErr     = 509
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
	Method     string     // 方法 get/post/put
	VerifyHttp VerifyHttp // 验证的方法
	Timeout    uint32     // 请求超时时间 秒
	Debug      bool       // 是否开启Debug模式
}

func (r *Request) GetDebug() bool {

	return r.Debug
}

func (r *Request) IsParameterLegal() (err error) {

	r.VerifyHttp = HttpCode

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

/***************************  校验信息  ********************************/

func HttpCode(request *Request, response *http.Response) (code int, isSucceed bool) {

	defer response.Body.Close()
	code = response.StatusCode
	if code == http.StatusOK {
		isSucceed = true
	}

	if request.GetDebug() {
		body, err := ioutil.ReadAll(response.Body)
		fmt.Printf("请求结果 httpCode:%d body:%s err:%v", response.StatusCode, string(body), err)
	}

	return
}
