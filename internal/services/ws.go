package services

import "C"
import (
	"context"
	"fmt"
	"github.com/LuccChagas/my-chat-app/internal/models"
	"github.com/LuccChagas/my-chat-app/internal/repository"
	ws "github.com/LuccChagas/my-chat-app/internal/websocket"
	"github.com/LuccChagas/my-chat-app/pkg"
	"github.com/gorilla/websocket"
	"github.com/labstack/gommon/log"
	"github.com/rabbitmq/amqp091-go"
	"strings"
	"time"
)

const (
	pongWait   = 60 * time.Second
	writeWait  = 10 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

var newline = []byte{'\n'}

type WsService struct {
	repository repository.RepositoryInterface
	rabbit     pkg.RabbitMQConnection
}

func NewWsService(repository repository.RepositoryInterface, rabbit pkg.RabbitMQConnection) *WsService {
	return &WsService{
		repository: repository,
		rabbit:     rabbit,
	}
}

func (s *WsService) ReadingPool(ctx context.Context, client *ws.Client) {
	defer func() {
		client.Hub.Unregister <- client
		client.Conn.Close()
	}()

	client.Conn.SetReadLimit(280)
	client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected error: %v", err)
			}
			break
		}

		msgStr := string(message)

		if strings.HasPrefix(msgStr, "/stock=") {
			stockCode := strings.TrimPrefix(msgStr, "/stock=")

			err := s.PublishStockRequest(stockCode)
			if err != nil {
				log.Printf("Error publishing stock: %v", err)
			}

			confirmationMsg := fmt.Sprintf("Processing command for stock code: %s", stockCode)
			timestamp := time.Now().Format("15:04:05")
			client.Hub.Broadcast <- []byte(fmt.Sprintf("[%s] %s", timestamp, confirmationMsg))
			continue
		}

		timestamp := time.Now().Format("15:04:05")
		formattedMsg := fmt.Sprintf("[%s] %s: %s", timestamp, client.Nickname, msgStr)

		newMsg := models.Message{
			Timestamp: time.Now(),
			Content:   msgStr,
			Author:    client.Nickname,
		}
		client.Hub.Messages = append(client.Hub.Messages, newMsg)
		if len(client.Hub.Messages) > 50 {
			client.Hub.Messages = client.Hub.Messages[len(client.Hub.Messages)-50:]
		}

		client.Hub.Broadcast <- []byte(formattedMsg)
	}
}

func (s *WsService) WritingPool(ctx context.Context, client *ws.Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			return

		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {

				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(client.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.Send)
			}
			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (s *WsService) PublishStockRequest(stockCode string) error {
	conn := s.rabbit
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	queue, err := ch.QueueDeclare(
		"mq_stock_code_req",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	body := stockCode
	err = ch.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		return err
	}

	log.Printf("Stock command - %s - published", stockCode)
	return nil
}

func (s *WsService) GetStockResponse(hub *ws.Hub) error {
	conn := s.rabbit

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	queue, err := ch.QueueDeclare(
		"mq_stock_code_res",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	msgs, err := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	go func() {
		for d := range msgs {
			log.Printf("Received response from queue: %s", string(d.Body))
			hub.Broadcast <- d.Body
		}
	}()
	return nil
}
