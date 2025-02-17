package models

import "time"

type Message struct {
	Timestamp time.Time
	Content   string
	Author    string
}
