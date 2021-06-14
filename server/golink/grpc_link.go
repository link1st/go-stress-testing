// Package golink 连接
package golink

import (
	"context"
	"go-stress-testing/server/statistics"
	"sync"
	"time"

	"go-stress-testing/helper"
	pb "go-stress-testing/proto"

	"go-stress-testing/model"
	"go-stress-testing/server/client"
)

// Grpc grpc 接口请求
func Grpc(chanID uint64, ch chan<- *model.RequestResults, totalNumber uint64, wg *sync.WaitGroup,
	request *model.Request, ws *client.GrpcSocket) {
	defer func() {
		wg.Done()
	}()
	defer func() {
		_ = ws.Close()
	}()
	for i := uint64(0); i < totalNumber; i++ {
		grpcRequest(chanID, ch, i, request, ws)
	}
	return
}

// grpcRequest 请求
func grpcRequest(chanID uint64, ch chan<- *model.RequestResults, i uint64, request *model.Request,
	ws *client.GrpcSocket) {
	var (
		startTime = time.Now()
		isSucceed = false
		errCode   = model.HTTPOk
	)
	// 需要发送的数据
	conn := ws.GetConn()
	if conn == nil {
		errCode = model.RequestErr
	} else {
		// TODO::请求接口示例
		c := pb.NewApiServerClient(conn)
		var (
			ctx = context.Background()
			req = &pb.Request{
				UserName: request.Body,
			}
		)
		rsp, err := c.HelloWorld(ctx, req)
		// fmt.Printf("rsp:%+v", rsp)
		if err != nil {
			errCode = model.RequestErr
		} else {
			// 200 为成功
			if rsp.Code != 200 {
				errCode = model.RequestErr
			} else {
				isSucceed = true
			}
		}
	}
	requestTime := uint64(helper.DiffNano(startTime))
	statistics.RequestTimeList = append(statistics.RequestTimeList, requestTime)
	requestResults := &model.RequestResults{
		Time:      requestTime,
		IsSucceed: isSucceed,
		ErrCode:   errCode,
	}
	requestResults.SetID(chanID, i)
	ch <- requestResults
}
