package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"GoChatRPC/server/service"
	pb "GoChatRPC/proto/chat/v1"
)

func main() {
	// 創建 TCP 監聽器，監聽 50051 端口
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("無法監聽端口: %v", err)
	}

	// 創建 gRPC 服務器
	grpcServer := grpc.NewServer()

	// 註冊聊天服務
	chatServer := service.NewChatServer()
	pb.RegisterChatServiceServer(grpcServer, chatServer)

	// 啟動服務器
	log.Printf("開始監聽 :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("無法啟動服務器: %v", err)
	}
}
