package dto

import (
	"github.com/google/uuid"
	"time"
)

type NewClientRequest struct {
	RequestID         uuid.UUID  `json:"request_id"`
	ClientInformation ClientInfo `json:"client_info"`
}

type DeleteClientRequest struct {
	RequestID uuid.UUID `json:"request_id"`
	ClientID  uuid.UUID `json:"client_id"`
}

type UpdateClientRequest struct {
	RequestID      uuid.UUID  `json:"request_id"`
	ClientID       uuid.UUID  `json:"client_id"`
	NewInformation ClientInfo `json:"client_info"`
}

type AudioJsonRequest struct {
	RequestID uuid.UUID `json:"request_id"`
	ClientID  uuid.UUID `json:"client_id"`
	Timestamp time.Time `json:"timestamp"`
	Payload   string    `json:"payload"`
}

type AudioProtoRequest struct {
	RequestID uuid.UUID `json:"request_id"`
	ClientID  uuid.UUID `json:"client_id"`
	Timestamp time.Time `json:"timestamp"`
	Payload   []byte    `json:"payload"`
}

type ClientInfo struct {
	FullName            string   `json:"full_name"`
	LocationName        string   `json:"location_name"`
	Latitude            string   `json:"latitude"`
	Longitude           string   `json:"longitude"`
	NotificationMethods []string `json:"notification_methods"`
}
