/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-21
* Time: 15:43
 */

package golink

import (
	"go-stress-testing/heper"
	"go-stress-testing/model"
	"go-stress-testing/server/client"
	"sync"
	"time"
)

// http go link
func Http(chanId uint64, ch chan<- *model.RequestResults, totalNumber uint64, wg *sync.WaitGroup, request *model.Request) {

	defer func() {
		wg.Done()
	}()

	// fmt.Printf("启动协程 编号:%05d \n", chanId)
	for i := uint64(0); i < totalNumber; i++ {

		var (
			startTime = time.Now()
			isSucceed = false
			errCode   = model.HttpOk
		)

		resp, err := client.HttpRequest(request.Method, request.Url, request.GetBody(), request.Headers, request.Timeout)
		requestTime := uint64(heper.DiffNano(startTime))
		// resp, err := server.HttpGetResp(request.Url)
		if err != nil {
			errCode = model.RequestErr // 请求错误
		} else {
			// 验证请求是否成功
			errCode, isSucceed = request.VerifyHttp(request, resp)
		}

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
