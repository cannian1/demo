package main

import (
	"context"
	"demo/grpc-demo/add_server/pb"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type AddService struct {
	pb.UnimplementedAddServer
}

// Sum 对两个数字求和
func (AddService) Sum(ctx context.Context, req *pb.SumRequest) (*pb.SumResponse, error) {
	return &pb.SumResponse{
		V: req.A + req.B,
	}, nil
}

// Concat 方法拼接两个字符串
func (AddService) Concat(ctx context.Context, req *pb.ConcatRequest) (*pb.ConcatResponse, error) {
	return &pb.ConcatResponse{V: req.A + req.B}, nil
}

func main() {
	// 监听本地的8972端口
	lis, err := net.Listen("tcp", ":8972")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()                  // 创建gRPC服务器
	pb.RegisterAddServer(s, &AddService{}) // 在gRPC服务端注册服务
	// 启动服务
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
