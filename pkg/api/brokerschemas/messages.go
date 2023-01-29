package brokerschemas

import (
	"github.com/Imm0bilize/gunshot-api-service/internal/entities"
	"github.com/google/uuid"
)

type AudioMessage struct {
	RequestID uuid.UUID             `json:"requestID"`
	Payload   entities.AudioMessage `json:"payload"`
}
