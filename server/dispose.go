/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-21
* Time: 15:42
 */

package server

import (
	"fmt"
	"go-stress-testing/model"
	"go-stress-testing/server/client"
	"go-stress-testing/server/golink"
	"go-stress-testing/server/statistics"
	"go-stress-testing/server/verify"
	"sync"
	"time"
)

// 注册验证器
func init() {
	model.RegisterVerifyHttp("statusCode", verify.HttpStatusCode)
	model.RegisterVerifyHttp("json", verify.HttpJson)
	model.RegisterVerifyWebSocket("json", verify.WebSocketJson)
}

// 处理函数
func Dispose(concurrency, totalNumber uint64, request *model.Request) {

	// 设置接收数据缓存
	ch := make(chan *model.RequestResults, 1000)
	var (
		// TODO::容易丢数据 或不及时返回
		wg          sync.WaitGroup // 发送数据完成
		wgReceiving sync.WaitGroup // 数据处理完成
	)

	wgReceiving.Add(1)
	go statistics.ReceivingResults(concurrency, ch, &wgReceiving)

	for i := uint64(0); i < concurrency; i++ {
		wg.Add(1)
		switch request.Form {
		case model.FormTypeHttp:
			go golink.Http(i, ch, totalNumber, &wg, request)
		case model.FormTypeWebSocket:

			// 连接以后再启动协程
			ws := client.NewWebSocket(request.Url)
			err := ws.GetConn()
			if err != nil {
				fmt.Println("连接失败:", i, err)

				continue
			}

			go golink.WebSocket(i, ch, totalNumber, &wg, request, ws)
		default:
			// 类型不支持
			wg.Done()
		}
	}

	// 等待所有的数据都发送完成
	wg.Wait()

	// 延时1毫秒 确保数据都处理完成了
	time.Sleep(1 * time.Millisecond)
	close(ch)

	// 数据全部处理完成了
	wgReceiving.Wait()

	return
}
