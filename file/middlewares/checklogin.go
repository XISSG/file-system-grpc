package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/xissg/file-system-grpc/consul"
	pb "github.com/xissg/file-system-grpc/internal"
	"github.com/xissg/file-system-grpc/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
)

func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceName := "account"
		conn, err := consul.NewClientConn(serviceName)
		if err != nil {
			log.Printf("failed to connect user service %v", err)
			c.Abort()
			return
		}

		tokenStr := c.Request.Header.Get("Authorization")
		//初始化grpc客户端
		client := pb.NewUserServiceClient(conn)
		req := &pb.Token{
			TokenStr: tokenStr,
		}

		//进行远程调用
		res, err := client.CheckToken(context.Background(), req)
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				log.Printf("server internal error")
				utils.ErrResponse(c, http.StatusInternalServerError, codes.Internal, "server internal error")
				c.Abort()
				return
			}
			switch st.Code() {
			case codes.Internal:
				log.Printf("failed to parse token")
				utils.ErrResponse(c, http.StatusBadRequest, codes.Internal, "failed to parse token")
				c.Abort()
			case codes.InvalidArgument:
				log.Printf("invalid token string")
				utils.ErrResponse(c, http.StatusBadRequest, codes.InvalidArgument, "invalid token string")
				c.Abort()
			case codes.DeadlineExceeded:
				log.Printf("session has expired")
				utils.ErrResponse(c, http.StatusBadRequest, codes.DeadlineExceeded, "session has expired")
				c.Abort()
			default:
				log.Printf("server internal error")
				utils.ErrResponse(c, http.StatusInternalServerError, codes.Internal, "server internal error")
				c.Abort()
			}

			return
		}

		if res == nil || !res.Success {
			log.Printf("authorization error, please login")
			utils.ErrResponse(c, http.StatusInternalServerError, codes.Internal, "failed to login")
			c.Abort()
			return
		}
		c.Set("userid", res.UserId)
		c.Next()
	}
}
