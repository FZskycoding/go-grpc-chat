package main

import (
    "context"
    "log"
    "time"

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

    log.Printf("發送者 (user123) 已啟動，準備發送消息...")

    // 創建並發送消息
    sendReq := &pb.SendMessageRequest{
        Receiver: "user456",
        Content:  "你好！這是一條來自sender的消息消息",
    }

    resp, err := client.SendMessage(context.Background(), sendReq)
    if err != nil {
        log.Fatalf("發送消息失敗: %v", err)
    }

    if resp.Success {
        log.Printf("\n消息發送成功：\n  - 消息ID: %s\n  - 發送時間: %s\n", 
            resp.Message.Id, 
            time.Unix(resp.Message.Timestamp, 0).Format("2006-01-02 15:04:05"))
    } else {
        log.Printf("消息發送失敗")
    }

    // 等待一下，讓用戶看到結果
    time.Sleep(time.Second)
}
