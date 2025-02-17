package services

import (
	"context"
	"github.com/LuccChagas/my-chat-app/internal/models"
	"github.com/LuccChagas/my-chat-app/internal/websocket"
	socket "github.com/LuccChagas/my-chat-app/internal/websocket"
)

type UserServiceInterface interface {
	CreateUser(ctx context.Context, user models.UserRequest) (models.UserResponse, error)
	GetUser(ctx context.Context, ID string) (models.UserResponse, error)
	GetAllUsers(ctx context.Context) ([]models.UserResponse, error)
	GetUserByUsername(ctx context.Context, username string) (models.UserResponse, error)
}

type WsServiceInterface interface {
	ReadingPool(ctx context.Context, client *websocket.Client)
	WritingPool(ctx context.Context, client *websocket.Client)
	PublishStockRequest(stock string)
	GetStockResponse(hub *socket.Hub) error
}
