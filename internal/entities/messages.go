package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct {
	Payload     []byte             `json:"payload"`
	Timestamp   time.Time          `json:"timestamp"`
	MessageType string             `json:"messageType"`
	ID          primitive.ObjectID `json:"ID"`
}
