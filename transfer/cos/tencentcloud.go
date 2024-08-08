package cos

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	pb "github.com/xissg/file-system-grpc/internal"
	"net/http"
	"net/url"
)

const rawURL = "https://xissg-buket-1319206652.cos.ap-chengdu.myqcloud.com"
const secretId = "AKIDiO6mpkXGsQVi7UksudTBc4cBj9Z7siaX"
const secretKey = "AKbm0InO0eN4lL3kWRds9UI6pQuKh2Gr"

type Client struct {
	client *cos.Client
}

func NewClient() *Client {
	u, _ := url.Parse(rawURL)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretId,
			SecretKey: secretKey,
		},
	})
	return &Client{
		client,
	}
}

func (c *Client) IsExist(key string) {
	ok, err := c.client.Object.IsExist(context.Background(), key)
	if err != nil {
		fmt.Printf("err")
		return
	}
	fmt.Println(ok)
}

func (c *Client) Upload(file *pb.File) error {
	_, _, err := c.client.Object.Upload(context.Background(), file.Checksum, file.FileAddr, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Download(file *pb.File) error {
	_, err := c.client.Object.GetToFile(context.Background(), file.Checksum, file.FileAddr, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Delete(file *pb.File) error {
	_, err := c.client.Object.Delete(context.Background(), file.Checksum, nil)
	if err != nil {
		return err
	}
	return nil
}
