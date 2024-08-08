package consul

import (
	"fmt"
	"log"
	"net"
)

func NewServerConn(serviceID string, serviceName string, ip string, port int) (net.Listener, error) {
	client := newClient(serviceID)
	if client == nil {
		panic("failed to connect consul")
	}

	err := client.registerService(serviceName, ip, port)
	if err != nil {
		log.Printf("failed to register service %v", serviceName)
		return nil, err
	}

	addr := fmt.Sprintf("%v:%v", ip, port)
	//监听端口
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("failed to listen: %v", err)
	}
	log.Printf("server listening %v", addr)
	return lis, nil
}

func Close() {}
