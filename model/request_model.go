/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-15
* Time: 18:19
 */

package model

import "net/http"

// 验证器
type Verify interface {
	GetCode() int    // 有一个方法，返回code为200为成功
	GetResult() bool // 返回是否成功
}

// 200为成功
type VerifyHttp func(response *http.Response) (code int)

// 请求结果
type Request struct {
	Url        string // Url
	Method     string // 方法 get/post/put
	VerifyHttp VerifyHttp
}

// 请求结果
type RequestResults struct {
	Time      uint64 // 请求时间 纳秒
	IsSucceed bool   // 是否请求成功
	ErrCode   uint32 // 错误码
}
