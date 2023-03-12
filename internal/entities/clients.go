package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Client struct {
	ID           primitive.ObjectID `json:"ID" bson:"_id"`
	LocationName string             `json:"locationName" bson:"locationName"`
	FullName     string             `json:"fullName" bson:"fullName"`
	Latitude     float64            `json:"latitude" bson:"latitude"`
	Longitude    float64            `json:"longitude" bson:"longitude"`
}
