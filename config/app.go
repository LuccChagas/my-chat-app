package config

import (
	"database/sql"
	db "github.com/LuccChagas/my-chat-app/db/sqlc"
	"github.com/LuccChagas/my-chat-app/internal/handlers"
	"github.com/LuccChagas/my-chat-app/internal/repository"
	"github.com/LuccChagas/my-chat-app/internal/routers"
	"github.com/LuccChagas/my-chat-app/internal/services"
	"github.com/LuccChagas/my-chat-app/internal/websocket"
	"github.com/rabbitmq/amqp091-go"
)

type App struct {
	Server *routers.Router
}

type RepositoryInstance struct {
	Repository *repository.Repository
}

type ServiceInstance struct {
	UserService *services.UserService
	WsService   *services.WsService
}

type HandlerInstance struct {
	UserHandler *handlers.UserHandler
	WsHandler   *handlers.WsHandler
}

func newRepositoryInstance(sqlDB *sql.DB) *RepositoryInstance {
	return &RepositoryInstance{
		Repository: repository.NewRepository(sqlDB, db.New(sqlDB)),
	}
}

func newHandlerInstance(serviceInstance *ServiceInstance, ws *websocket.Hub) *HandlerInstance {
	return &HandlerInstance{
		UserHandler: handlers.NewUserHandler(serviceInstance.UserService),
		WsHandler:   handlers.NewWsHandler(serviceInstance.WsService, ws),
	}
}

func newServiceInstance(repoInstance *RepositoryInstance, rabbit *amqp091.Connection) *ServiceInstance {
	return &ServiceInstance{
		UserService: services.NewUserService(repoInstance.Repository),
		WsService:   services.NewWsService(repoInstance.Repository, rabbit),
	}
}

func NewApp(db *sql.DB, hub *websocket.Hub, rabbit *amqp091.Connection) *App {

	repoInstance := newRepositoryInstance(db)
	serviceInstance := newServiceInstance(repoInstance, rabbit)
	handlerInstance := newHandlerInstance(serviceInstance, hub)

	server := routers.NewRouter(
		handlerInstance.UserHandler,
		handlerInstance.WsHandler,
	)

	err := serviceInstance.WsService.GetStockResponse(hub)
	if err != nil {
		return nil
	}

	return &App{
		Server: server,
	}
}
