# Go gRPC 聊天系統

這是一個使用 Go 語言和 gRPC 實現的聊天系統，支持單向和雙向即時通訊。

## 功能特點

- 支持單向消息發送和接收
- 支持雙向即時聊天
- 系統消息通知（用戶加入/離開）
- 使用 gRPC 串流實現即時通訊
- 支持多個用戶同時在線

## 技術棧

- Go 1.24.0
- gRPC
- Protocol Buffers
- Google UUID

## 功能說明

### 單向聊天

單向聊天模式提供了簡單的消息發送和接收功能：

- **發送者 (Sender)**：可以發送單條消息給指定接收者
- **接收者 (Receiver)**：持續監聽並接收發送給自己的消息

### 雙向聊天

雙向聊天模式提供了完整的即時通訊功能：

- 支持多人同時在線聊天
- 即時消息推送
- 系統消息通知（用戶加入/離開聊天室）
- 優雅的程式退出處理

## 如何使用

### 1. 啟動伺服器

```bash
go run server/main.go
```

### 2. 使用雙向聊天

```bash
# 啟動聊天客戶端（需要提供用戶ID）
go run client/TwoWay_chat/main.go <使用者ID>
```

### 3. 使用單向聊天

```bash
# 啟動接收者
go run client/OneWay_chat/receiver/receiver.go

# 啟動發送者
go run client/OneWay_chat/sender/sender.go
```

## 通訊協議

系統使用 Protocol Buffers 定義了以下主要訊息類型：

- `Message`：基本消息結構
- `ChatMessage`：聊天消息
- `SendMessageRequest/Response`：發送消息請求/回應
- `ReceiveRequest`：接收消息請求
- `ConnectRequest`：連接請求

## 系統特性

- **併發安全**：使用互斥鎖確保多用戶操作的安全性
- **優雅退出**：支持正確的連接關閉和資源釋放
- **可擴展性**：模組化設計，易於添加新功能
- **即時通訊**：使用 gRPC 串流實現低延遲通訊
