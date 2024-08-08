package test

import (
	"context"
	"github.com/xissg/file-system-grpc/consul"
	"github.com/xissg/file-system-grpc/internal"
	"testing"
)

func TestRegister(t *testing.T) {
	serviceName := "account"
	conn, err := consul.NewClientConn(serviceName)
	if err != nil {
		panic(err)
	}
	client := internal.NewUserServiceClient(conn)

	account := &internal.Account{
		Username: "xissg",
		Password: "0709",
	}
	_, err = client.Register(context.Background(), account)
	if err != nil {
		panic(err)
	}
}
