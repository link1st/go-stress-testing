/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-15
* Time: 21:03
 */

package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"strings"
	"time"
)

type GrpcSocket struct {
	conn    *grpc.ClientConn
	address string
}

func NewGrpcSocket(address string) (s *GrpcSocket) {
	var (
		newAddr string
	)
	arr := strings.Split(address, "//")
	if len(arr) >= 2 {
		newAddr = arr[1]
	}
	s = &GrpcSocket{
		address: newAddr,
	}
	return
}

func (g *GrpcSocket) getAddress() (address string) {
	return g.address
}

// 关闭
func (g *GrpcSocket) Close() (err error) {
	if g == nil {
		return
	}
	if g.conn == nil {
		return
	}
	g.conn.Close()
	return
}

// Link 建立连接
func (g *GrpcSocket) Link() (err error) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	conn, err := grpc.DialContext(ctx, g.address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("GetConn: 连接失败 address:%s %w", g.address, err)
	}
	g.conn = conn
	return
}

func (g *GrpcSocket) GetConn() (conn *grpc.ClientConn) {
	return g.conn
}
