package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AudioMessage struct {
	ID        primitive.ObjectID `json:"ID"`
	Timestamp time.Time          `json:"timestamp"`
	Payload   []byte             `json:"payload"`
}
