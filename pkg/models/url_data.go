package models

import "time"

type UrlData struct {
	ShortURL  string    `json:"short_url"      dynamodbav:"short_url"`
	FullURL   string    `json:"full_url"      dynamodbav:"full_url"`
	CreatedAt time.Time `json:"created_at"  dynamodbav:"created_at"`
}
