package golink

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/link1st/go-stress-testing/helper"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"

	"github.com/link1st/go-stress-testing/model"
)

const STRING_AUTH string = "auth"
const STRING_ACCT string = "acct"


// Grpc grpc 接口请求
func Radius(ctx context.Context, chanID uint64, ch chan<- *model.RequestResults, totalNumber uint64, wg *sync.WaitGroup,
	request *model.Request) {
	defer func() {
		wg.Done()
	}()
	for i := uint64(0); i < totalNumber; i++ {
		authRequest(chanID, ch, i, request)
	}
	return
}

// radius authRequest 请求
func authRequest(chanID uint64, ch chan<- *model.RequestResults, i uint64, request *model.Request) {
	var (
		startTime = time.Now()
		isSucceed = false
		errCode   = int(radius.CodeAccessAccept)
	)
	username := request.Headers["username"]
    password := request.Headers["password"]
    secret := request.Headers["secret"]
    if username == "" {
        username = "tim"
    }
    if password == "" {
        password = "12345678"
    }
	packet := radius.New(radius.CodeAccessRequest, []byte(secret))
	host := request.URL
	rfc2865.UserName_SetString(packet, username)
	rfc2865.UserPassword_SetString(packet, password)
	rfc2865.NASPortType_Set(packet, rfc2865.NASPortType_Value_Ethernet)
	rfc2865.ServiceType_Set(packet, rfc2865.ServiceType_Value_FramedUser)
	rfc2865.NASIdentifier_Set(packet, []byte(`benchmark`))
	rsp, err := radius.Exchange(context.Background(), packet, host)
	if err != nil {
		errCode = model.RequestErr
	} else {
		if rsp.Code != radius.CodeAccessAccept {
			errCode = int(rsp.Code)
		} else {
			isSucceed = true
		}
	}
	requestTime := uint64(helper.DiffNano(startTime))
	requestResults := &model.RequestResults{
		Time:      requestTime,
		IsSucceed: isSucceed,
		ErrCode:   errCode,
	}
	requestResults.SetID(chanID, i)
	ch <- requestResults
}
