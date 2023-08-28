package models

import "time"

type User struct {
	Name string `json:"name"`
}

type Segment struct {
	Slug             string    `json:"slug"`
	ExpirationDate   time.Time `json:"expiration_date"`
	RandomPercentage float64   `json:"random_percentage"`
}

type DeleteUserRequest struct {
	UserId int `json:"user_id"`
}

type UserSegmentsRequest struct {
	UserId int `json:"user_id"`
}

type UpdateSegmentsRequest struct {
	UserId int       `json:"user_id"`
	Add    []Segment `json:"add"`
	Remove []string  `json:"remove"`
}

type ReportRequest struct {
	UserId    int    `json:"user_id"`
	YearMonth string `json:"yearMonth"`
}
