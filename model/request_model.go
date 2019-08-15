/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-15
* Time: 18:19
 */

package model

// 请求结果
type RequestResults struct {
	Time      uint64 // 请求时间 纳秒
	IsSucceed bool   // 是否请求成功
	ErrCode   uint32 // 错误码
}
