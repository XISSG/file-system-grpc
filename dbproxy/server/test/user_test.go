package test

import (
	"context"
	"fmt"
	"github.com/xissg/file-system-grpc/consul"
	"github.com/xissg/file-system-grpc/internal"
	"testing"
)

func TestAddUser(t *testing.T) {
	serviceName := "dbproxy"
	conn, err := consul.NewClientConn(serviceName)
	if err != nil {
		panic(err)
	}
	client := internal.NewDBServiceClient(conn)
	res, err := client.AddUser(context.Background(), &internal.User{Id: 1, Username: "xissg", Password: "0709"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", res)
}

func TestGetUser(t *testing.T) {
	serviceName := "dbproxy"
	conn, err := consul.NewClientConn(serviceName)
	if err != nil {
		panic(err)
	}
	client := internal.NewDBServiceClient(conn)
	res, err := client.GetUser(context.Background(), &internal.UserName{Username: "xissg"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", res)
}
func TestCheckUserExist(t *testing.T) {
	serviceName := "dbproxy"
	conn, err := consul.NewClientConn(serviceName)
	if err != nil {
		panic(err)
	}
	client := internal.NewDBServiceClient(conn)
	username := &internal.UserName{
		Username: "xissg",
	}
	res, err := client.CheckUserExist(context.Background(), username)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", res)
}

func TestAddFile(t *testing.T) {
	serviceName := "dbproxy"
	conn, err := consul.NewClientConn(serviceName)
	if err != nil {
		panic(err)
	}
	client := internal.NewDBServiceClient(conn)
	res, err := client.AddFile(context.Background(), &internal.File{Id: 1, UserId: 1821071640117645312, Checksum: "1", FileName: "test", FileAddr: "tmp/test.txt", FileSize: 1024, Status: 0})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", res)
}

func TestUpdateFileStatus(t *testing.T) {

	serviceName := "dbproxy"
	conn, err := consul.NewClientConn(serviceName)
	if err != nil {
		panic(err)
	}
	client := internal.NewDBServiceClient(conn)
	res, err := client.UpdateFileStatus(context.Background(), &internal.UpdateRequest{Status: 1})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", res)

}
func TestGetFileByUserID(t *testing.T) {

	serviceName := "dbproxy"
	conn, err := consul.NewClientConn(serviceName)
	if err != nil {
		panic(err)
	}
	client := internal.NewDBServiceClient(conn)
	res, err := client.GetFileByUserID(context.Background(), &internal.UserID{UserId: 1821071640117645312})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", res)
}
func TestGetFileByChecksum(t *testing.T) {

	serviceName := "dbproxy"
	conn, err := consul.NewClientConn(serviceName)
	if err != nil {
		panic(err)
	}
	client := internal.NewDBServiceClient(conn)
	res, err := client.GetFileByChecksum(context.Background(), &internal.Checksum{Checksum: "1"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", res)
}
func TestDeleteFile(t *testing.T) {

	serviceName := "dbproxy"
	conn, err := consul.NewClientConn(serviceName)
	if err != nil {
		panic(err)
	}
	client := internal.NewDBServiceClient(conn)
	if err != nil {
		panic(err)
	}

	res, err := client.DeleteFile(context.Background(), &internal.Checksum{Checksum: "1"})
	fmt.Printf("%v", res)
}
