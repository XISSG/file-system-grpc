package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xissg/file-system-grpc/consul"
	pb "github.com/xissg/file-system-grpc/internal"
	queue "github.com/xissg/file-system-grpc/rabbitmq"
	"github.com/xissg/file-system-grpc/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
)

type FileClient struct {
	client pb.DBServiceClient
	conn   *grpc.ClientConn

	queue *queue.Client
}

func NewFileClient() *FileClient {
	serviceName := "dbproxy"
	conn, err := consul.NewClientConn(serviceName)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	//初始化grpc客户端
	client := pb.NewDBServiceClient(conn)

	//初始化rabbitmq客户端
	q, err := queue.NewClient()
	if err != nil {
		log.Printf("failed to connect rabbitmq %v", err)
	}

	return &FileClient{
		client: client,
		conn:   conn,
		queue:  q,
	}
}

const (
	exchangeUpload   = "file-upload-exchange"
	exchangeDownload = "file-download-exchange"
)

func (file *FileClient) Upload(c *gin.Context) {
	defer func() {
		file.conn.Close()
	}()

	userID, ok := c.Get("userid")
	if !ok {
		utils.ErrResponse(c, http.StatusBadRequest, codes.Unauthenticated, "please login first")
		return
	}
	//从请求中读取文件
	f, err := c.FormFile("file")
	if err != nil {
		utils.ErrResponse(c, http.StatusBadRequest, codes.InvalidArgument, "failed to get file")
		return
	}

	//创建保存文件的文件夹
	dst := path.Join("tmp", f.Filename)
	err = os.MkdirAll("tmp", 0754)
	if err != nil {
		utils.ErrResponse(c, http.StatusBadRequest, codes.Internal, "failed to save file")
		return
	}

	//TODO:优化文件保存逻辑，安全校验等
	//保存文件
	err = c.SaveUploadedFile(f, dst)
	if err != nil {
		utils.ErrResponse(c, http.StatusBadRequest, codes.Internal, "failed to save file")
		return
	}

	//生成文件元信息
	checksum, err := utils.GetSH1Checksum(dst)

	if err != nil {
		utils.ErrResponse(c, http.StatusBadRequest, codes.Internal, "failed to save file")
		return
	}

	//保存用户的文件元信息
	id, ok := userID.(int64)
	if !ok {
		utils.ErrResponse(c, http.StatusInternalServerError, codes.Internal, "server public error")
		return
	}

	meta := &pb.File{
		Id:       utils.Snowflake(),
		UserId:   id,
		FileName: f.Filename,
		FileSize: f.Size,
		Checksum: checksum,
		FileAddr: dst,
		Status:   utils.WaitingBackup,
	}

	res, err := file.client.AddFile(context.Background(), meta)
	if err != nil || !res.Success {
		utils.ErrResponse(c, http.StatusBadRequest, codes.Internal, "failed to save file")
		return
	}

	//发送到消息队列中，然后传输到腾讯云
	data, err := json.Marshal(meta)
	if err != nil {
		log.Printf("failed to marshal data , %v", err)
	}

	err = file.queue.Publish("fanout", exchangeUpload, "", data)
	if err != nil {
		log.Printf("failed to send msg to queue, %v", err)
	}

	utils.Response(c, http.StatusOK, codes.OK, "upload file success", nil)
}

// 获取用户上传的文件
func (file *FileClient) GetFile(c *gin.Context) {
	userID, ok := c.Get("userid")
	if !ok {
		utils.ErrResponse(c, http.StatusInternalServerError, codes.Internal, "server public error")
		return
	}

	id := userID.(int64)
	idObj := &pb.UserID{
		UserId: id,
	}

	files, err := file.client.GetFileByUserID(context.Background(), idObj)
	if err != nil {
		utils.ErrResponse(c, http.StatusInternalServerError, codes.Internal, "failed to load files you uploaded")
		return
	}

	c.JSON(http.StatusOK, files)
}

func (file *FileClient) Download(c *gin.Context) {
	checksum := c.Query("file")

	checksumObj := &pb.Checksum{
		Checksum: checksum,
	}
	f, err := file.client.GetFileByChecksum(context.Background(), checksumObj)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename*=UTF-8''"+f.FileName)
	c.Header("Content-Type", "application/octet-stream")

	if err != nil {
		utils.ErrResponse(c, http.StatusBadRequest, codes.Internal, "failed to download file")
		return
	}

	if f == nil {
		utils.ErrResponse(c, http.StatusBadRequest, codes.Internal, "the file you download is not exist")
		return
	}

	//如果本地没有该文件，则向消息队列中发一个消息，让transfer服务服务将文件下载到本地
	ok := utils.CheckFileExists(f.FileAddr)

	if !ok && f.Status == utils.SuccessBackup {
		data, err := json.Marshal(f)
		if err != nil {
			log.Printf("failed to marshal data , %v", err)
		}
		err = file.queue.Publish("fanout", exchangeDownload, "", data)
		if err != nil {
			log.Printf("failed to send data into queue for uploading file to cos")
		}
		utils.ErrResponse(c, http.StatusBadRequest, codes.Aborted, "getting file from cloud please waiting for a minute")
		return
	}

	c.File(f.FileAddr)
}

func (file *FileClient) Delete(c *gin.Context) {
	checksum := c.Query("file")

	checksumObj := &pb.Checksum{
		Checksum: checksum,
	}

	res, err := file.client.DeleteFile(context.Background(), checksumObj)

	if err != nil || !res.Success {
		utils.ErrResponse(c, http.StatusInternalServerError, codes.Internal, "failed to delete file")
		return
	}

	utils.Response(c, http.StatusOK, codes.OK, "delete file success", nil)
}

// 分块上传
func (file *FileClient) MultipartUpload(c *gin.Context) {
	//获取文件校验和，分块序列号，分块数据，文件名，分片总数
	userID, _ := c.Get("userid")

	f1, err := c.FormFile("file")
	if err != nil {
		utils.ErrResponse(c, http.StatusBadRequest, codes.InvalidArgument, "no file uploaded")
		return
	}
	fileName := c.PostForm("fileName")
	checkSum := c.PostForm("checksum")
	chunkNum := c.PostForm("chunk")
	total := c.PostForm("total")

	//临时存储到对应文件临时目录中</tmp/校验和/chunk>
	tmpPath := filepath.Join("tmp", checkSum, chunkNum)
	err = c.SaveUploadedFile(f1, tmpPath)
	if err != nil {
		utils.ErrResponse(c, http.StatusBadRequest, codes.Internal, "failed to save chunk "+chunkNum)
		return
	}

	tmpDir := filepath.Dir(tmpPath)
	dstPath := filepath.Join("tmp", fileName)

	totalCount, err := strconv.Atoi(total)
	receiveCount := countFiles(tmpDir)

	if totalCount == receiveCount {

		err := combineFile(checkSum, dstPath, tmpDir)
		if err != nil {
			utils.ErrResponse(c, http.StatusBadRequest, codes.Internal, "failed to combine upload file")
			return
		}

		//存储文件信息
		fileInfo, err := os.Stat(dstPath)
		id, ok := userID.(int64)
		if !ok {
			utils.ErrResponse(c, http.StatusBadRequest, codes.Internal, "invalid user account")
			return
		}
		f := &pb.File{
			FileName: fileName,
			FileSize: fileInfo.Size(),
			Checksum: checkSum,
			FileAddr: dstPath,
			UserId:   id,
		}
		res, err := file.client.AddFile(context.Background(), f)
		if err != nil || res == nil || !res.Success {
			utils.ErrResponse(c, http.StatusInternalServerError, codes.Internal, "server public error")
			return
		}
		utils.Response(c, http.StatusOK, codes.OK, "success receive file", nil)
		return
	}

	msg := fmt.Sprintf("success receive number chunck %v,already receive chunk count: %v,total chunk count: %v", chunkNum, receiveCount, totalCount)
	utils.Response(c, http.StatusOK, codes.OK, msg, nil)
}

func countFiles(dir string) int {
	count := 0

	err := filepath.Walk(dir, func(_ string, info os.FileInfo, _ error) error {
		if !info.IsDir() { // 检查是否为文件
			count++
		}
		return nil
	})

	if err != nil {
		return 0
	}

	return count
}

func combineFile(checksum string, dstPath, tmpDir string) error {

	//读出指定目录下的文件，并排序
	files, err := os.ReadDir(tmpDir)
	if err != nil {
		return err
	}
	var blocks []int
	for _, file := range files {
		if !file.IsDir() {
			name, err := strconv.Atoi(file.Name())
			if err != nil {
				continue
			}
			blocks = append(blocks, name)
		}
	}

	sort.Ints(blocks)

	f, err := os.Create(dstPath)
	defer f.Close()

	//合并文件
	for _, fileName := range blocks {
		name := strconv.Itoa(fileName)
		filePath := filepath.Join(tmpDir, name)
		content, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}
		_, err = f.Write(content)
		if err != nil {
			return err
		}
	}

	//校验和比较
	hash, err := utils.GetSH1Checksum(dstPath)
	if hash != checksum {
		return errors.New("file has damaged")
	}

	err = os.RemoveAll(tmpDir)
	if err != nil {
		log.Println("failed to remove temp file " + tmpDir)
	}
	return nil
}
