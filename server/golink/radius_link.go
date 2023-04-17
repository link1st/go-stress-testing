package golink

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/link1st/go-stress-testing/helper"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"layeh.com/radius/rfc2866"
	"layeh.com/radius/rfc3162"

	"github.com/link1st/go-stress-testing/model"
)

const STRING_AUTH string = "auth"
const STRING_ACCT string = "acct"

var STAGE uint64 = 0
var DefaultClient = &radius.Client{
	Retry:              time.Second,
	MaxPacketErrors:    10,
	InsecureSkipVerify: true,
}

// Grpc grpc 接口请求
func Radius(ctx context.Context, chanID uint64, ch chan<- *model.RequestResults, totalNumber uint64, wg *sync.WaitGroup,
	request *model.Request) {
	defer func() {
		wg.Done()
	}()
	if request.Headers["type"] == "acct" {
		if STAGE == 0 {
			STAGE, _ = strconv.ParseUint(request.Headers["stage"], 10, 64)
		}
		if totalNumber%STAGE != 0 {
			fmt.Println(fmt.Errorf("ERROR:  Radius Account must by of mutiple of stage"))
			return
		}
		for i := uint64(0); i < totalNumber/STAGE; i++ {
			acctRequest(chanID, ch, i+totalNumber/STAGE*0, request, rfc2866.AcctStatusType_Value_Start)
		}
		for i := uint64(0); i < totalNumber/STAGE; i++ {
			acctRequest(chanID, ch, i+totalNumber/STAGE*1, request, rfc2866.AcctStatusType_Value_InterimUpdate)
		}
		if STAGE == 3 {
			for i := uint64(0); i < totalNumber/STAGE; i++ {
				acctRequest(chanID, ch, i+totalNumber/STAGE*2, request, rfc2866.AcctStatusType_Value_Stop)
			}
		}

	} else {
		for i := uint64(0); i < totalNumber; i++ {
			authRequest(chanID, ch, i, request)
		}
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
	rfc2865.CallingStationID_Set(packet, []byte("ac:60:89:70:29:9d"))
	rfc2865.NASIPAddress_Set(packet, []byte(request.Headers["nasip"]))
	rfc2865.NASPort_Set(packet, rfc2865.NASPort(rand.Uint32()))
	rsp, err := DefaultClient.Exchange(context.Background(), packet, host)
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

// radius acctRequest 请求
func acctRequest(chanID uint64, ch chan<- *model.RequestResults, i uint64, request *model.Request, code rfc2866.AcctStatusType) {
	var (
		startTime = time.Now()
		isSucceed = false
		errCode   = int(radius.CodeAccountingResponse)
	)
	secret := request.Headers["secret"]
	username := fmt.Sprintf("acct_%02d_%07d", chanID, i)
	packet := radius.New(radius.CodeAccountingRequest, []byte(secret))
	host := request.URL
	rfc2865.UserName_SetString(packet, username)
	rfc2865.NASIdentifier_Set(packet, []byte(`benchmark`))
	rfc2865.CallingStationID_Set(packet, []byte("ac:60:89:70:29:9d"))
	rfc2865.NASIPAddress_Set(packet, []byte("192.168.10.255"))
	rfc2865.NASPort_Set(packet, rfc2865.NASPort(rand.Uint32()))
	rfc2865.NASPortType_Set(packet, rfc2865.NASPortType_Value_Ethernet)
	rfc2865.ServiceType_Set(packet, rfc2865.ServiceType_Value_FramedUser)
	rfc2865.FramedIPAddress_Set(packet, []byte("192.168.10.254"))
	// rfc2869.NASPortID_Set(packet, []byte("slot=63;subslot=9;port=1;vlanid=255;vlanid2=2438;"))

	rfc2866.AcctStatusType_Set(packet, code)
	rfc2866.AcctSessionID_Set(packet, []byte(username))
	rfc2866.AcctDelayTime_Set(packet, 0)

	if code == rfc2866.AcctStatusType_Value_InterimUpdate {
		rfc2866.AcctSessionTime_Set(packet, 10)
		rfc2866.AcctInputOctets_Set(packet, 10*1024*1024)
		rfc2866.AcctOutputOctets_Set(packet, 10*1024*1024)
		rfc2866.AcctInputPackets_Set(packet, 10*1024)
		rfc2866.AcctOutputPackets_Set(packet, 10*1024)
		rfc3162.FramedIPv6Prefix_Set(packet, &net.IPNet{IP: net.IPv6interfacelocalallnodes, Mask: net.CIDRMask(1, 56)})
	}
	if code == rfc2866.AcctStatusType_Value_Stop {
		rfc2866.AcctSessionTime_Set(packet, 20)
		rfc2866.AcctInputOctets_Set(packet, 20*1024*1024)
		rfc2866.AcctOutputOctets_Set(packet, 20*1024*1024)
		rfc2866.AcctInputPackets_Set(packet, 20*1024)
		rfc2866.AcctOutputPackets_Set(packet, 20*1024)
	}

	_, err := DefaultClient.Exchange(context.Background(), packet, host)
	if err != nil {
		fmt.Println(err)
		errCode = model.RequestErr
	} else {
		isSucceed = true
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
