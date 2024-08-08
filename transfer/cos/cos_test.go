package cos

import (
	pb "github.com/xissg/file-system-grpc/dbproxy/public"
	"testing"
)

func TestClient_IsExist(t *testing.T) {
	client := NewClient()
	client.IsExist("IDE_crack.txt")
	file := &pb.File{
		FileName: "gfx_win_101.5768.exe",
		FileAddr: "./gfx_win_101.5768.exe",
		Checksum: "db7e8abdab55feeaa3a7c00ef4f1f9c846be4b6b",
		FileSize: 911501896,
	}
	//err := client.Upload(file)
	//if err != nil {
	//	panic(err)
	//}
	err := client.Download(file)
	if err != nil {
		panic(err)
	}
}
