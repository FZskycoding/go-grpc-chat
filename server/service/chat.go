package service

import (
    "fmt"
    "log"
    "sync"
    "time"

    pb "GoChatRPC/proto/chat/v1"
)

// ChatServer 實現 gRPC 服務器
type ChatServer struct {
    pb.UnimplementedChatServiceServer
    chatStreams map[string]pb.ChatService_ChatStreamServer // 儲存用戶串流連接
    mu         sync.RWMutex                               // 保護 map 的併發訪問
}

// NewChatServer 創建一個新的聊天服務器實例
func NewChatServer() *ChatServer {
    return &ChatServer{
        chatStreams: make(map[string]pb.ChatService_ChatStreamServer),
    }
}

// ChatStream 實現雙向串流聊天
func (s *ChatServer) ChatStream(stream pb.ChatService_ChatStreamServer) error {
    // 等待第一條消息以獲取用戶ID
    firstMsg, err := stream.Recv()
    if err != nil {
        return err
    }
    userId := firstMsg.UserId

    // 註冊串流
    s.mu.Lock()
    s.chatStreams[userId] = stream
    s.mu.Unlock()

    // 清理函數
    defer func() {
        s.mu.Lock()
        delete(s.chatStreams, userId)
        s.mu.Unlock()
    }()

    // 廣播系統消息：用戶加入
    s.broadcast(&pb.ChatMessage{
        UserId:    "系統",
        Content:   fmt.Sprintf("用戶 %s 已加入聊天", userId),
        Timestamp: time.Now().Unix(),
    })

    // 持續接收和廣播消息
    for {
        msg, err := stream.Recv()
        if err != nil {
            // 廣播系統消息：用戶離開
            s.broadcast(&pb.ChatMessage{
                UserId:    "系統",
                Content:   fmt.Sprintf("用戶 %s 已離開聊天", userId),
                Timestamp: time.Now().Unix(),
            })
            return err
        }

        // 廣播收到的消息
        s.broadcast(msg)
    }
}

// broadcast 向所有連接的用戶廣播消息
func (s *ChatServer) broadcast(msg *pb.ChatMessage) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    for _, stream := range s.chatStreams {
        if err := stream.Send(msg); err != nil {
            log.Printf("發送消息失敗: %v", err)
        }
    }
}
