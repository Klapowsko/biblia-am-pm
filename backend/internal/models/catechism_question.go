package models

type CatechismQuestion struct {
	ID            int    `json:"id"`
	QuestionNumber int   `json:"question_number"`
	QuestionText   string `json:"question_text"`
	AnswerText     string `json:"answer_text"`
}

