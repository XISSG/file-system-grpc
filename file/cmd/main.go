package main

import "github.com/xissg/file-system-grpc/file/router"

func main() {
	r := router.Router()
	_ = r.Run(":9090")
}
