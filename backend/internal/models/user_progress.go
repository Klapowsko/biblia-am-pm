package models

import "time"

type UserProgress struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	ReadingPlanID   int       `json:"reading_plan_id"`
	Date            time.Time `json:"date"`
	MorningCompleted bool     `json:"morning_completed"`
	EveningCompleted bool     `json:"evening_completed"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
}

