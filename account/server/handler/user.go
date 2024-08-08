package handler

import (
	"context"
	"github.com/xissg/file-system-grpc/consul"
	pb "github.com/xissg/file-system-grpc/internal"
	"github.com/xissg/file-system-grpc/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

type User struct {
	pb.UnimplementedUserServiceServer
	client pb.DBServiceClient
	conn   *grpc.ClientConn
}

func NewUserClient() *User {
	serviceName := "dbproxy"
	conn, err := consul.NewClientConn(serviceName)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	client := pb.NewDBServiceClient(conn)

	return &User{
		conn:   conn,
		client: client,
	}
}

func (u *User) Register(_ context.Context, account *pb.Account) (*pb.RegisterStatus, error) {
	res := &pb.RegisterStatus{
		Success: false,
	}
	if account == nil || account.Username == "" || account.Password == "" {
		return res, status.New(codes.InvalidArgument, "invalid user account").Err()
	}

	ok, err := u.client.CheckUserExist(context.Background(), &pb.UserName{Username: account.Username})
	if err != nil {
		return res, status.New(codes.Internal, "server get user error").Err()
	}

	if ok.Exist {
		return res, status.New(codes.AlreadyExists, "username already exist").Err()
	}

	addRes, err := u.client.AddUser(context.Background(), &pb.User{Username: account.Username, Password: utils.MD5Crypt(account.Password)})
	if err != nil || !addRes.Success {
		return res, status.New(codes.Internal, "server add user error").Err()
	}

	res.Success = true
	return res, nil
}

func (u *User) Login(_ context.Context, user *pb.Account) (*pb.LoginStatus, error) {
	res := &pb.LoginStatus{
		TokenStr: "",
		Success:  false,
	}

	if user == nil || user.Username == "" {
		return res, status.New(codes.InvalidArgument, "invalid user account").Err()
	}

	ok, err := u.client.CheckUserExist(context.Background(), &pb.UserName{Username: user.Username})
	if err != nil {
		return res, status.New(codes.Internal, "server internal error").Err()
	}

	if !ok.Exist {
		return res, status.New(codes.NotFound, "please register first").Err()
	}
	data, err := u.client.GetUser(context.Background(), &pb.UserName{Username: user.Username})

	if err != nil || data == nil {
		return res, status.New(codes.Internal, "server get user error").Err()
	}

	if data.Password != utils.MD5Crypt(user.Password) {
		return res, status.New(codes.Unauthenticated, "wrong username or password").Err()
	}

	expireTime := time.Hour * 24
	tokenStr, err := utils.Generate(data.Id, expireTime)

	if err != nil {
		return res, status.New(codes.Internal, "server get user error").Err()
	}

	res.Success = true
	res.TokenStr = tokenStr
	return res, nil
}

func (u *User) CheckToken(_ context.Context, token *pb.Token) (*pb.CheckStatus, error) {

	res := &pb.CheckStatus{
		Success: false,
	}
	if token == nil {
		return res, status.New(codes.InvalidArgument, "invalid token string").Err()
	}

	session, err := utils.Parse(token.TokenStr)
	if err != nil || session == nil {
		return res, status.New(codes.Internal, "failed to parse token").Err()
	}

	if session.Expire.Unix() < time.Now().Unix() {
		return res, status.New(codes.DeadlineExceeded, "session has expired").Err()
	}

	res.Success = true
	res.UserId = session.UserID
	return res, nil
}
