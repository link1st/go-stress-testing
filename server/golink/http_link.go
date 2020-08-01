/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-21
* Time: 15:43
 */

package golink

import (
	"sync"

	"go-stress-testing/model"
	"go-stress-testing/server/client"
)

// http go link
func Http(chanId uint64, ch chan<- *model.RequestResults, totalNumber uint64, wg *sync.WaitGroup, request *model.Request) {

	defer func() {
		wg.Done()
	}()

	// fmt.Printf("启动协程 编号:%05d \n", chanId)
	for i := uint64(0); i < totalNumber; i++ {

		list := getRequestList(request)

		isSucceed, errCode, requestTime := sendList(list)

		requestResults := &model.RequestResults{
			Time:      requestTime,
			IsSucceed: isSucceed,
			ErrCode:   errCode,
		}

		requestResults.SetId(chanId, i)

		ch <- requestResults
	}

	return
}
// 多个接口分步压测
func sendList(requestList []*model.Request) (isSucceed bool, errCode int, requestTime uint64) {

	errCode = model.HttpOk
	for _, request := range requestList {
		succeed, code, u := send(request)
		isSucceed = succeed
		errCode = code
		requestTime = requestTime + u
		if succeed == false {

			break
		}
	}

	return
}

// send 发送一次请求
func send(request *model.Request) (bool, int, uint64) {
	var (
		// startTime = time.Now()
		isSucceed = false
		errCode   = model.HttpOk
	)

	newRequest := getRequest(request)
	// newRequest := request

	resp, requestTime, err := client.HttpRequest(newRequest.Method, newRequest.Url, newRequest.GetBody(), newRequest.Headers, newRequest.Timeout)
	// requestTime := uint64(heper.DiffNano(startTime))
	if err != nil {
		errCode = model.RequestErr // 请求错误
	} else {
		// 验证请求是否成功
		errCode, isSucceed = newRequest.GetVerifyHttp()(newRequest, resp)
	}
	return isSucceed, errCode, requestTime
}
