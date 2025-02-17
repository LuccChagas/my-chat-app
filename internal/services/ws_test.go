package services_test

import (
	"context"
	"fmt"
	"github.com/LuccChagas/my-chat-app/internal/models"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"

	"github.com/LuccChagas/my-chat-app/internal/services"
	ws "github.com/LuccChagas/my-chat-app/internal/websocket"
)

// FakeWSConn implementa a interface WSConn definida em ws
type FakeWSConn struct {
	readMessages [][]byte
	readIndex    int
	closed       bool
}

func (f *FakeWSConn) ReadMessage() (int, []byte, error) {
	if f.readIndex >= len(f.readMessages) {
		return 0, nil, fmt.Errorf("no message")
	}
	msg := f.readMessages[f.readIndex]
	f.readIndex++
	return websocket.TextMessage, msg, nil
}

func (f *FakeWSConn) SetReadLimit(limit int64)            {}
func (f *FakeWSConn) SetReadDeadline(t time.Time) error   { return nil }
func (f *FakeWSConn) SetWriteDeadline(t time.Time) error  { return nil }
func (f *FakeWSConn) SetPongHandler(h func(string) error) {}
func (f *FakeWSConn) Close() error {
	f.closed = true
	return nil
}

// Para os testes de ReadingPool não precisamos usar NextWriter e WriteMessage,
// mas podemos implementar métodos mínimos se necessário:
func (f *FakeWSConn) NextWriter(messageType int) (io.WriteCloser, error) {
	return &FakeWriteCloser{}, nil
}

func (f *FakeWSConn) WriteMessage(messageType int, data []byte) error {
	return nil
}

type FakeWriteCloser struct{}

func (fwc *FakeWriteCloser) Write(p []byte) (int, error) {
	return len(p), nil
}
func (fwc *FakeWriteCloser) Close() error {
	return nil
}

// FakeRabbitConn implementa a interface RabbitMQConnection
type FakeRabbitConn struct{}

func (f *FakeRabbitConn) Channel() (*amqp091.Channel, error) {
	return nil, fmt.Errorf("fake rabbit error")
}

func TestReadingPool_NonStockMessage(t *testing.T) {
	fakeConn := &FakeWSConn{
		readMessages: [][]byte{[]byte("Hello")},
	}

	fakeHub := &ws.Hub{
		Broadcast:  make(chan []byte, 10),
		Messages:   []models.Message{},
		Register:   make(chan *ws.Client, 10),
		Unregister: make(chan *ws.Client, 10),
	}

	client := &ws.Client{
		Hub:      fakeHub,
		Conn:     fakeConn,
		Send:     make(chan []byte, 10),
		Nickname: "TestUser",
	}

	svc := services.NewWsService(nil, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	go svc.ReadingPool(ctx, client)
	time.Sleep(200 * time.Millisecond)

	select {
	case msg := <-fakeHub.Broadcast:
		assert.True(t, strings.Contains(string(msg), "TestUser: Hello"), "The message should contain the nickname and content")
	default:
		t.Error("No message was broadcasted")
	}
}

func TestReadingPool_StockCommand(t *testing.T) {
	fakeConn := &FakeWSConn{
		readMessages: [][]byte{[]byte("/stock=GOOGL.US")},
	}

	fakeHub := &ws.Hub{
		Broadcast:  make(chan []byte, 10),
		Messages:   []models.Message{},
		Register:   make(chan *ws.Client, 10),
		Unregister: make(chan *ws.Client, 10),
	}

	client := &ws.Client{
		Hub:      fakeHub,
		Conn:     fakeConn,
		Send:     make(chan []byte, 10),
		Nickname: "TestUser",
	}

	fakeRabbit := &FakeRabbitConn{}
	svc := services.NewWsService(nil, fakeRabbit)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	go svc.ReadingPool(ctx, client)
	time.Sleep(200 * time.Millisecond)

	select {
	case msg := <-fakeHub.Broadcast:
		assert.True(t, strings.Contains(string(msg), "Processing command for stock code: GOOGL.US"),
			"The confirmation message should be broadcasted")
	default:
		t.Error("No confirmation message was broadcasted")
	}
}
