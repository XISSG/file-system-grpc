package main

import (
	"github.com/xissg/file-system-grpc/account/server/handler"
	"github.com/xissg/file-system-grpc/consul"
	pb "github.com/xissg/file-system-grpc/internal"
	"google.golang.org/grpc"
	"log"
)

func main() {
	serviceID := "account"
	serviceName := "account"
	ip := "localhost"
	port := 10001
	lis, err := consul.NewServerConn(serviceID, serviceName, ip, port)
	if err != nil {
		log.Printf("failed to start listen %v", err)
	}
	//初始化grpc服务，注册方法
	s := grpc.NewServer()
	user := handler.NewUserClient()
	pb.RegisterUserServiceServer(s, user)
	if err := s.Serve(lis); err != nil {
		log.Printf("failed to serve %v", err)
	}
}
