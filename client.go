package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/linxlib/conv"
	"github.com/linxlib/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpc_startup/proto"
	"io/ioutil"
	"log"
	"time"
)

const (
	address = "localhost:50051"
)

func main() {
	// 注意这里的deploy， 需要和证书公钥生成时的 Common Name 对应
	c, err := credentials.NewClientTLSFromFile("./server.crt", "deploy")
	if err != nil {
		log.Fatalf("credentials.NewClientTLSFromFile err: %v", err)
	}
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(c))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := proto.NewFileServiceClient(conn)

	// 30秒的上下文, 传输大文件适当扩大时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	bs, _ := ioutil.ReadFile("./1.0-window.7z")
	filelen := conv.Int64(len(bs))
	h := sha256.New()
	h.Write(bs)
	myhash := fmt.Sprintf("%x", h.Sum(nil))
	logs.Info("myhash:", myhash)
	start := time.Now()
	r, err := client.Upload(ctx, &proto.FSReq{
		DstDir:   "ehw",
		ProjName: "dsaudg",
		Name:     "dasgf",
		ProjType: 1,
		Hash:     myhash,
		Filelen:  filelen,
		IfReboot: false,
		File:     bs,
	})
	end := time.Now().Sub(start).Seconds()
	kb := filelen / 1024
	logs.Info("time:", end, " file size:", kb, "KB")
	if err != nil {
		log.Fatalf("could not upload: %v", err)
	}
	log.Printf("Upload: %s", r.Message)
}
