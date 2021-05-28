/**
* Created by GoLand.
* User: link1st
* Date: 2021/2/3
* Time: 23:44
 */

package main

import (
	"context"
	"fmt"
	pb "go-stress-testing/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	// port 监听端口
	port = ":8099"
)

// server is used to implement helloWorld.GreeterServer.
type server struct {
	pb.UnimplementedApiServerServer
}

// HelloWorld hello world 接口
func (s *server) HelloWorld(_ context.Context, req *pb.Request) (rsp *pb.Response, err error) {
	rsp = &pb.Response{
		Code: 200,
		Msg:  "success",
		Data: fmt.Sprintf("hello %s !", req.UserName),
	}
	return
}

// main 主函数
func main() {
	fmt.Println("trpc server 启动中...")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterApiServerServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
