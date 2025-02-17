package handlers

import (
	"context"
	"fmt"
	"github.com/LuccChagas/my-chat-app/internal/services"
	socket "github.com/LuccChagas/my-chat-app/internal/websocket"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
)

type WsHandler struct {
	socket  *socket.Hub
	service *services.WsService
}

func NewWsHandler(s *services.WsService, socket *socket.Hub) *WsHandler {
	return &WsHandler{socket: socket, service: s}
}

func (h *WsHandler) WsHandler(c echo.Context) error {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}
	nickname := "Unknown"
	if val, ok := sess.Values["nickname"]; ok {
		nickname = fmt.Sprintf("%v", val)
	}

	client := &socket.Client{
		Hub:      h.socket,
		Conn:     ws,
		Send:     make(chan []byte, 256),
		Nickname: nickname,
	}

	client.Hub.Register <- client

	for _, msg := range client.Hub.Messages {
		formattedMsg := fmt.Sprintf("[%s] %s: %s \n", msg.Timestamp.Format("15:04:05"), msg.Author, msg.Content)
		client.Send <- []byte(formattedMsg)
	}

	ctx := context.Background()
	go h.service.ReadingPool(ctx, client)
	go h.service.WritingPool(ctx, client)

	return nil
}
