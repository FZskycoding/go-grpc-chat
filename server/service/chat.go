package service

import (
	"context"
	"sync"
	"time"

	pb "GoChatRPC/proto/chat/v1"
	"github.com/google/uuid"
)

// ChatServer 實現 gRPC 服務器
type ChatServer struct {
	pb.UnimplementedChatServiceServer                             // 必須嵌入未實現的服務器
	activeConnections                 map[string]chan *pb.Message // 儲存在線用戶的消息通道
	mu                                sync.RWMutex                // 保護 map 的併發訪問
}

// NewChatServer 創建一個新的聊天服務器實例
func NewChatServer() *ChatServer {
	return &ChatServer{
		// 初始化在線用戶的消息通道 map
		activeConnections: make(map[string]chan *pb.Message),
	}
}

// ReceiveMessages 實現消息接收的串流服務
func (s *ChatServer) ReceiveMessages(req *pb.ReceiveRequest, stream pb.ChatService_ReceiveMessagesServer) error {
	// 為用戶創建消息通道
	msgChan := make(chan *pb.Message, 100) // 緩衝區大小為 100，避免阻塞
	userId := req.UserId

	// 註冊用戶連接
	s.mu.Lock()
	s.activeConnections[userId] = msgChan
	s.mu.Unlock()

	// 當函數返回時清理資源
	defer func() {
		s.mu.Lock()
		delete(s.activeConnections, userId) // 移除用戶連接
		close(msgChan)                      // 關閉通道
		s.mu.Unlock()
	}()

	// 持續監聽消息通道並發送給客戶端
	for {
		select {
		case msg := <-msgChan:
			// 收到新消息，發送給客戶端
			if err := stream.Send(msg); err != nil {
				return err
			}
		case <-stream.Context().Done():
			// 客戶端斷開連接或取消請求
			return nil
		}
	}
}

// SendMessage 實現發送消息的方法
func (s *ChatServer) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	// 創建新的消息
	msg := &pb.Message{
		Id:        uuid.New().String(), // 生成唯一的消息ID
		Sender:    "user123",           // 這裡應該是真實的用戶ID
		Content:   req.Content,         // 使用請求中的消息內容
		Timestamp: time.Now().Unix(),   // 設置當前時間戳
	}

	// 檢查接收者是否在線並發送消息
	s.mu.RLock()
	if msgChan, ok := s.activeConnections[req.Receiver]; ok {
		// 嘗試發送消息到接收者的通道
		select {
		case msgChan <- msg:
			// 消息成功發送到通道
		default:
			// 通道已滿，消息可能無法即時發送
			// 在實際應用中，可能需要處理這種情況（例如：存儲到數據庫）
		}
	}
	s.mu.RUnlock()

	// 返回響應
	return &pb.SendMessageResponse{
		Message: msg,
		Success: true,
	}, nil
}
