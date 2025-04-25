package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"

    pb "GoChatRPC/proto/chat/v1"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    // 連接到服務器
    conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("無法連接到服務器: %v", err)
    }
    defer conn.Close()

    // 創建客戶端
    client := pb.NewChatServiceClient(conn)

    log.Printf("接收者 (user456) 已啟動，等待接收消息...")
    log.Printf("按 Ctrl+C 停止接收\n")

    // 創建接收消息的請求
    receiveReq := &pb.ReceiveRequest{
        UserId: "user456",  // 接收者ID
    }

    // 開始接收消息
    stream, err := client.ReceiveMessages(context.Background(), receiveReq)
    if err != nil {
        log.Fatalf("無法開始接收消息: %v", err)
    }

    // 設置中斷信號處理
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

    // 使用 goroutine 接收消息
    go func() {
        for {
            msg, err := stream.Recv()
            if err != nil {
                log.Printf("接收消息時發生錯誤: %v", err)
                return
            }
            log.Printf("\n收到新消息：\n  - 發送者: %s\n  - 內容: %s\n  - 消息ID: %s\n", 
                msg.Sender, msg.Content, msg.Id)
        }
    }()

    // 等待中斷信號
    <-sigChan
    log.Printf("\n接收者已停止運行")
}
