package internal

import (
	"time"
)

type User struct {
	ID         string    `bson:"_id,omitempty" json:"id,omitempty"`
	Email      string    `bson:"email" json:"email"`
	Name       string    `bson:"name" json:"name"`
	Role       string    `bson:"role" json:"role"`
	Password   string    `bson:"passwordHash" json:"password"`
	RegisterAt time.Time `bson:"registerAt" json:"registerAt"`
}

type Category struct {
	ID   string `bson:"_id,omitempty" json:"id,omitempty"`
	Name string `bson:"name" json:"name"`
}

type Location struct {
	Latitude  float64 `bson:"latitude" json:"latitude"`
	Longitude float64 `bson:"longitude" json:"longitude"`
	Address   string  `bson:"address" json:"address"`
}

type Restaurant struct {
	ID         string   `bson:"_id,omitempty" json:"id,omitempty"`
	Name       string   `bson:"name" json:"name"`
	CategoryID string   `bson:"categoryId" json:"categoryId"`
	Location   Location `bson:"location" json:"location"`
}

type Review struct {
	ID           string    `bson:"_id,omitempty" json:"id,omitempty"`
	UserID       string    `bson:"userId" json:"userId"`
	RestaurantID string    `bson:"restaurantId" json:"restaurantId"`
	Text         string    `bson:"text" json:"text"`
	Rating       float64   `bson:"rating" json:"rating"`
	CreatedAt    time.Time `bson:"createdAt" json:"createdAt"`
}

type NLPResult struct {
	ID        string   `bson:"_id,omitempty" json:"id,omitempty"`
	ReviewID  string   `bson:"reviewId" json:"reviewId"`
	Sentiment string   `bson:"sentiment" json:"sentiment"`
	Keywords  []string `bson:"keywords" json:"keywords"`
}

type Rating struct {
	ID            string  `bson:"_id,omitempty" json:"id,omitempty"`
	RestaurantID  string  `bson:"restaurantId" json:"restaurantId"`
	AverageRating float64 `bson:"averageRating" json:"averageRating"`
	ReviewCount   int     `bson:"reviewCount" json:"reviewCount"`
}

type AdminLog struct {
	ID         string    `bson:"_id,omitempty" json:"id,omitempty"`
	AdminID    string    `bson:"adminId" json:"adminId"`
	ActionType string    `bson:"actionType" json:"actionType"`
	Details    string    `bson:"details" json:"details"`
	CreatedAt  time.Time `bson:"createdAt" json:"createdAt"`
}

type Favorite struct {
	ID           string    `bson:"_id,omitempty" json:"id,omitempty"`
	UserID       string    `bson:"userId" json:"userId"`
	RestaurantID string    `bson:"restaurantId" json:"restaurantId"`
	AddedAt      time.Time `bson:"addedAt" json:"addedAt"`
}

type RatingResponse struct {
	TextReview string  `json:"review"`
	Status     bool    `json:"status"`
	Rating     float64 `json:"rating"`
}
