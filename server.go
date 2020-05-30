package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/linxlib/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpc_startup/proto"
	"log"
	"math"
	"net"
)

const (
	port = ":50051"
)

type server struct {
}

func (s *server) verifyFile(file []byte, hash string, length int64) bool {
	h := sha256.New()
	h.Write(file)
	myHash := fmt.Sprintf("%x", h.Sum(nil))
	logs.Info("hash:", hash, " myHash:", myHash, " len:", length, " myLen:", len(file))
	return hash == myHash
}

func (s *server) Upload(ctx context.Context, in *proto.FSReq) (*proto.FSResp, error) {
	if !s.verifyFile(in.File, in.Hash, in.Filelen) {
		return &proto.FSResp{
			Status:  false,
			Message: "数据包哈希校验失败，请重新部署",
		}, nil
	}
	return &proto.FSResp{
		Status:  true,
		Message: "received",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logs.Fatalf("failed to listen: %v", err)
	}
	c, err := credentials.NewServerTLSFromFile("./server.crt", "./server.key")
	if err != nil {
		log.Fatalf("credentials.NewServerTLSFromFile err: %v", err)
	}
	//由于要发送较大的压缩包，默认为 4M。
	//如果需要向客户端发送大文件则增加一条grpc.MaxSendMsgSize()
	s := grpc.NewServer(
		grpc.Creds(c),
		grpc.MaxRecvMsgSize(math.MaxInt64))
	proto.RegisterFileServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
