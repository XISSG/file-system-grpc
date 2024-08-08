package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/xissg/file-system-grpc/consul"
	pb "github.com/xissg/file-system-grpc/internal"
	"github.com/xissg/file-system-grpc/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
)

type UserClient struct {
	client pb.UserServiceClient
	conn   *grpc.ClientConn
}

func NewUserClient() *UserClient {
	serviceName := "account"
	conn, err := consul.NewClientConn(serviceName)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	//初始化grpc客户端
	client := pb.NewUserServiceClient(conn)

	return &UserClient{
		client: client,
		conn:   conn,
	}

}

func (u *UserClient) Register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	account := &pb.Account{
		Username: username,
		Password: password,
	}
	_, err := u.client.Register(context.Background(), account)

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Printf("server internal error")
			utils.ErrResponse(c, http.StatusInternalServerError, codes.Internal, "server internal error")
			return
		}

		switch st.Code() {
		case codes.InvalidArgument:
			log.Printf("invalid user account")
			utils.ErrResponse(c, http.StatusBadRequest, codes.InvalidArgument, "invalid username or password")
		case codes.Internal:
			log.Printf("server internal error")
			utils.ErrResponse(c, http.StatusBadRequest, codes.Internal, "server internal error")
		case codes.AlreadyExists:
			log.Printf("account already exist")
			utils.ErrResponse(c, http.StatusBadRequest, codes.AlreadyExists, "account already exist")
		default:
			log.Printf("failed to register")
			utils.ErrResponse(c, http.StatusInternalServerError, codes.Internal, "failed to register")
		}

		return
	}

	c.HTML(http.StatusFound, "login.html", nil)
}

func (u *UserClient) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	account := &pb.Account{
		Username: username,
		Password: password,
	}

	res, err := u.client.Login(context.Background(), account)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Printf("server internal error")
			utils.ErrResponse(c, http.StatusInternalServerError, codes.Internal, "server internal error")
			return
		}

		switch st.Code() {
		case codes.InvalidArgument:
			log.Printf("invalid user account")
			utils.ErrResponse(c, http.StatusBadRequest, codes.InvalidArgument, "invalid username or password")
		case codes.Internal:
			log.Printf("server internal error")
			utils.ErrResponse(c, http.StatusBadRequest, codes.Internal, "server internal error")
		case codes.NotFound:
			log.Printf("user account not found")
			utils.ErrResponse(c, http.StatusBadRequest, codes.NotFound, "you have not registered")
		case codes.Unauthenticated:
			log.Printf("account already exist")
			utils.ErrResponse(c, http.StatusBadRequest, codes.Unauthenticated, "wrong username or password")
		default:
			log.Printf("failed to register")
			utils.ErrResponse(c, http.StatusInternalServerError, codes.Internal, "failed to register")
		}

		return
	}

	if res == nil || !res.Success {
		log.Printf("server internal error")
		utils.ErrResponse(c, http.StatusInternalServerError, codes.Internal, "server internal error")
		return
	}

	log.Printf("login success")
	utils.Response(c, http.StatusBadRequest, codes.OK, "login success", res.TokenStr)
}
