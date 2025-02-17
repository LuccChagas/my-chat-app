package main

import (
	"encoding/csv"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

const (
	stockAPIURLTemplate = "https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv"
	requestQueueName    = "mq_stock_code_req"
	responseQueueName   = "mq_stock_code_res"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func ConnRabbit() (*amqp091.Connection, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		os.Getenv("AMQP_USER"),
		os.Getenv("AMQP_PASS"),
		os.Getenv("AMQP_HOST"),
		os.Getenv("AMQP_PORT"))

	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func processStockCmd(stockCode string) (string, error) {
	url := fmt.Sprintf(stockAPIURLTemplate, stockCode)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("error calling API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("error parsing CSV: %w", err)
	}

	if len(records) < 2 {
		return "", fmt.Errorf("incomplete CSV")
	}
	data := records[1]
	if len(data) < 7 {
		return "", fmt.Errorf("unexpected CSV format")
	}
	closePrice := data[6]
	if closePrice == "N/D" {
		return fmt.Sprintf("%s quote is not available", strings.ToUpper(stockCode)), nil
	}

	responseMsg := fmt.Sprintf("%s quote is $%s per share", strings.ToUpper(stockCode), closePrice)
	return responseMsg, nil
}

func main() {
	conn, err := ConnRabbit()
	if err != nil {
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error opening channel: %v", err)
	}
	defer ch.Close()

	reqQueue, err := ch.QueueDeclare(
		requestQueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error declaring request queue: %v", err)
	}

	respQueue, err := ch.QueueDeclare(
		responseQueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error declaring response queue: %v", err)
	}

	msgs, err := ch.Consume(
		reqQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error registering consumer: %v", err)
	}

	log.Printf("Bot started. Waiting for messages on queue %q...", reqQueue.Name)
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			stockCode := string(d.Body)
			log.Printf("Command received: %s", stockCode)

			responseMsg, err := processStockCmd(stockCode)
			if err != nil {
				log.Printf("Error processing command - %s: %v", stockCode, err)
				continue
			}

			err = ch.Publish(
				"",
				respQueue.Name,
				false,
				false,
				amqp091.Publishing{
					ContentType: "text/plain",
					Body:        []byte(responseMsg),
				},
			)
			if err != nil {
				log.Printf("Error publishing response: %v", err)
				continue
			}
			log.Printf("Response published: %s", responseMsg)
		}
	}()

	<-forever
}
