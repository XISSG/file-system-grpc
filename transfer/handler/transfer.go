package handler

import (
	"context"
	"github.com/xissg/file-system-grpc/consul"
	pb "github.com/xissg/file-system-grpc/internal"
	queue "github.com/xissg/file-system-grpc/rabbitmq"
	"github.com/xissg/file-system-grpc/transfer/cos"
	"github.com/xissg/file-system-grpc/utils"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	exchangeUpload   = "file-upload-exchange"
	exchangeDownload = "file-download-exchange"
)

type Trans struct {
	q   *queue.Client
	cos *cos.Client

	client pb.DBServiceClient
	conn   *grpc.ClientConn
}

func NewTrans() *Trans {
	q, err := queue.NewClient()
	if err != nil {
		log.Printf("failed to connect rabbitmq")
	}

	cosClient := cos.NewClient()
	serviceName := "dbproxy"
	consulClient := consul.NewClient("")
	serviceAddr, err := consulClient.GetService(serviceName)

	//和对应服务建立连接
	conn, err := grpc.Dial(serviceAddr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	//初始化grpc客户端
	client := pb.NewDBServiceClient(conn)
	return &Trans{
		q:      q,
		cos:    cosClient,
		conn:   conn,
		client: client,
	}
}

func (t *Trans) TransUpload() {
	defer func() { t.q.Close() }()
	t.q.Consume("fanout", exchangeUpload, "", t.upload)
}

func (t *Trans) TransDownload() {
	defer func() { t.q.Close() }()
	t.q.Consume("fanout", exchangeDownload, "", t.download)
}

func (t *Trans) upload(file *pb.File) error {
	err := t.cos.Upload(file)
	if err != nil {
		return err
	}
	res, err := t.client.UpdateFileStatus(context.Background(), &pb.UpdateRequest{Status: utils.SuccessBackup})
	if err != nil || !res.Success {
		log.Printf("failed to upload file to cos %v", err)
		return err
	}
	log.Printf("backing up file to cos success")
	return nil
}

func (t *Trans) download(file *pb.File) error {
	err := t.cos.Download(file)
	if err != nil {
		log.Printf("failed to download file from cos to local path %v", err)
		return err
	}
	log.Printf("download file from cos to local path success")
	return nil
}
