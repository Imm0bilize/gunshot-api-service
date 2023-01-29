package dto

import "time"

type ClientInfo struct {
	LocationName        string   `json:"locationName" binding:"required"`
	FullName            string   `json:"fullName" binding:"required"`
	Latitude            float64  `json:"latitude" binding:"required"`
	Longitude           float64  `json:"longitude" binding:"required"`
	NotificationMethods []string `json:"notificationMethods" binding:"required"`
}

type UploadAudioRequest struct {
	Timestamp time.Time `url:"ts" binding:"required"`
	ID        string    `url:"id" binding:"required"`
}
