// Package golink 连接
package golink

import (
	"net/http"
	"sync"

	"go-stress-testing/model"
	"go-stress-testing/server/client"
)

// HTTP 请求
func HTTP(chanID uint64, ch chan<- *model.RequestResults, totalNumber uint64, wg *sync.WaitGroup,
	request *model.Request) {
	defer func() {
		wg.Done()
	}()
	// fmt.Printf("启动协程 编号:%05d \n", chanID)
	for i := uint64(0); i < totalNumber; i++ {
		list := getRequestList(request)
		isSucceed, errCode, requestTime, contentLength := sendList(list)
		requestResults := &model.RequestResults{
			Time:          requestTime,
			IsSucceed:     isSucceed,
			ErrCode:       errCode,
			ReceivedBytes: contentLength,
		}
		requestResults.SetID(chanID, i)
		ch <- requestResults
	}

	return
}

// sendList 多个接口分步压测
func sendList(requestList []*model.Request) (isSucceed bool, errCode int, requestTime uint64, contentLength int64) {
	errCode = model.HTTPOk
	for _, request := range requestList {
		succeed, code, u, length := send(request)
		isSucceed = succeed
		errCode = code
		requestTime = requestTime + u
		contentLength = contentLength + length
		if succeed == false {
			break
		}
	}
	return
}

// send 发送一次请求
func send(request *model.Request) (bool, int, uint64, int64) {
	var (
		// startTime = time.Now()
		isSucceed     = false
		errCode       = model.HTTPOk
		contentLength = int64(0)
		err           error
		resp          *http.Response
		requestTime   uint64
	)
	newRequest := getRequest(request)

	resp, requestTime, err = client.HTTPRequest(newRequest)

	if err != nil {
		errCode = model.RequestErr // 请求错误
	} else {
		contentLength = resp.ContentLength
		// 验证请求是否成功
		errCode, isSucceed = newRequest.GetVerifyHTTP()(newRequest, resp)
	}
	return isSucceed, errCode, requestTime, contentLength
}
