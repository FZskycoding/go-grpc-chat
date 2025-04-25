package main

import (
    "bufio"
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    pb "GoChatRPC/proto/chat/v1"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    // 檢查命令行參數
    if len(os.Args) < 2 {
        log.Fatalf("使用方式: %s <使用者ID>", os.Args[0])
    }
    userId := os.Args[1]

    // 連接到服務器
    conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("無法連接到服務器: %v", err)
    }
    defer conn.Close()

    // 創建客戶端
    client := pb.NewChatServiceClient(conn)

    // 建立聊天串流
    stream, err := client.ChatStream(context.Background())
    if err != nil {
        log.Fatalf("無法建立聊天串流: %v", err)
    }

    // 發送第一條消息以識別用戶
    firstMsg := &pb.ChatMessage{
        UserId:    userId,
        Content:   "",
        Timestamp: time.Now().Unix(),
    }
    if err := stream.Send(firstMsg); err != nil {
        log.Fatalf("無法發送首次連接消息: %v", err)
    }

    // 設置信號處理
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

    // 使用 goroutine 接收消息
    go func() {
        for {
            msg, err := stream.Recv()
            if err != nil {
                log.Printf("接收消息錯誤: %v", err)
                return
            }
            // 如果不是自己發送的消息才顯示
            if msg.UserId != userId {
                fmt.Printf("\n%s: %s\n> ", msg.UserId, msg.Content)
            }
        }
    }()

    // 讀取用戶輸入並發送
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Print("> ")
    
    // 用於通知主程序退出的通道
    done := make(chan bool)
    
    // 處理用戶輸入
    go func() {
        for scanner.Scan() {
            // 建立並發送聊天消息
            msg := &pb.ChatMessage{
                UserId:    userId,
                Content:   scanner.Text(),
                Timestamp: time.Now().Unix(),
            }
            if err := stream.Send(msg); err != nil {
                log.Printf("發送消息錯誤: %v", err)
                break
            }
            fmt.Print("> ")
        }
        done <- true
    }()

    // 等待退出信號或用戶輸入結束
    select {
    case <-sigChan:
        fmt.Println("\n收到退出信號，正在關閉連接...")
    case <-done:
        fmt.Println("\n輸入結束，正在關閉連接...")
    }
}
