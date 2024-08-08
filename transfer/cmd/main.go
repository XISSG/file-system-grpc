package main

import "github.com/xissg/file-system-grpc/transfer/handler"

func main() {
	trans := handler.NewTrans()
	ch := make(chan struct{})
	go trans.TransUpload()
	go trans.TransDownload()
	<-ch
}
