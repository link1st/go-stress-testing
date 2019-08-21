/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-21
* Time: 15:43
 */

package golink

import (
	"fmt"
	"go-stress-testing/heper"
	"go-stress-testing/model"
	"go-stress-testing/server/client"
	"sync"
	"time"
)

// web socket go link
func WebSocket(chanId uint64, ch chan<- *model.RequestResults, totalNumber uint64, wg *sync.WaitGroup, request *model.Request, ws *client.WebSocket) {

	defer func() {
		wg.Done()
	}()

	// fmt.Printf("启动协程 编号:%05d \n", chanId)

	defer func() {
		ws.Close()
	}()

	// 初始化请求
	for i := uint64(0); i < totalNumber; i++ {

		var (
			startTime = time.Now()
			isSucceed = false
			errCode   = model.HttpOk
		)

		// 需要发送的数据
		seq := fmt.Sprintf("%d_%d", chanId, i)
		err := ws.Write([]byte(`{"seq":"` + seq + `","cmd":"ping","data":{}}`))
		if err != nil {
			errCode = model.RequestErr // 请求错误
		} else {
			msg, err := ws.Read()
			if err != nil {
				errCode = model.ParseError
				fmt.Println("读取数据 失败~")
			} else {
				// fmt.Println(msg)
				errCode, isSucceed = request.VerifyWebSocket(request, seq, msg)
			}
		}

		requestTime := uint64(heper.DiffNano(startTime))

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
