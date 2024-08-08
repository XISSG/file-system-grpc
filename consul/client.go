package consul

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClientConn(serviceName string) (*grpc.ClientConn, error) {
	consulClient := newClient("")
	ip, err := consulClient.getService(serviceName)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.Dial(ip, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
