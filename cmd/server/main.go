package main

import (
	"github.com/LuccChagas/my-chat-app/config"
	"github.com/LuccChagas/my-chat-app/internal/websocket"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	hub := websocket.NewHub()
	go hub.Run()

	db, err := config.ConnDB()
	if err != nil {
		return
	}

	rabbit, err := config.ConnRabbit()
	if err != nil {
		return
	}

	app := config.NewApp(db, hub, rabbit)
	app.Server.Serve()
	log.Println("Servidor iniciado...")

	defer rabbit.Close()
	defer db.Close()
}
