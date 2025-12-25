package models

import "time"

type CatechismProgress struct {
	ID          int       `json:"id"`
	UserID    int       `json:"user_id"`
	QuestionID  int       `json:"question_id"`
	Date        time.Time `json:"date"`
	Completed   bool      `json:"completed"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

