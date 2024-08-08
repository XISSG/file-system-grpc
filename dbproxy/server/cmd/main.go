package main

import (
	"github.com/xissg/file-system-grpc/consul"
	"github.com/xissg/file-system-grpc/dbproxy/server/handler"
	pb "github.com/xissg/file-system-grpc/internal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"log"
	"time"
)

func main() {
	serviceID := "dbproxy"
	serviceName := "dbproxy"
	ip := "localhost"
	port := 10000
	lis, err := consul.NewServerConn(serviceID, serviceName, ip, port)
	if err != nil {
		log.Printf("failed to start listen %v", err)
	}
	//初始化grpc服务，注册方法
	s := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionIdle: 5 * time.Minute}))
	db := handler.NewDB()
	pb.RegisterDBServiceServer(s, db)
	if err := s.Serve(lis); err != nil {
		log.Printf("failed to serve %v", err)
	}
}
