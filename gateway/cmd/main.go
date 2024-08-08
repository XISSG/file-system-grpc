package main

import "github.com/xissg/file-system-grpc/gateway/router"

func main() {
	r := router.Router()
	_ = r.Run(":8080")
}
