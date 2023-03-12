package brokerschemas

import (
	"github.com/Imm0bilize/gunshot-api-service/internal/entities"
	"github.com/google/uuid"
)

type AudioMessage struct {
	Payload   entities.Message `json:"payload"`
	RequestID uuid.UUID        `json:"requestID"`
}
