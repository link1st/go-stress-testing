// Package client grpc 客户端
package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"
)

// GrpcSocket grpc
type GrpcSocket struct {
	conn    *grpc.ClientConn
	address string
}

// NewGrpcSocket new
func NewGrpcSocket(address string) (s *GrpcSocket) {
	var newAddr string
	arr := strings.Split(address, "//")
	if len(arr) >= 2 {
		newAddr = arr[1]
	}
	s = &GrpcSocket{
		address: newAddr,
	}
	return
}

// getAddress 获取地址
func (g *GrpcSocket) getAddress() (address string) {
	return g.address
}

// Close 关闭
func (g *GrpcSocket) Close() (err error) {
	if g == nil {
		return
	}
	if g.conn == nil {
		return
	}
	return g.conn.Close()
}

// Link 建立连接
func (g *GrpcSocket) Link() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, g.address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("getConn: 连接失败 address:%s %w", g.address, err)
	}
	g.conn = conn
	return
}

// GetConn 获取连接
func (g *GrpcSocket) GetConn() (conn *grpc.ClientConn) {
	return g.conn
}
