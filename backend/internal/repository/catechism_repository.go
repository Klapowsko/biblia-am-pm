package repository

import (
	"biblia-am-pm/internal/database"
	"biblia-am-pm/internal/models"
	"database/sql"
)

type CatechismRepository struct{}

func NewCatechismRepository() *CatechismRepository {
	return &CatechismRepository{}
}

func (r *CatechismRepository) GetByQuestionNumber(questionNumber int) (*models.CatechismQuestion, error) {
	query := `SELECT id, question_number, question_text, answer_text 
	          FROM westminster_catechism WHERE question_number = $1`
	
	question := &models.CatechismQuestion{}
	err := database.DB.QueryRow(query, questionNumber).Scan(
		&question.ID,
		&question.QuestionNumber,
		&question.QuestionText,
		&question.AnswerText,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	if err != nil {
		return nil, err
	}
	
	return question, nil
}

func (r *CatechismRepository) GetAll() ([]*models.CatechismQuestion, error) {
	query := `SELECT id, question_number, question_text, answer_text 
	          FROM westminster_catechism ORDER BY question_number`
	
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var questions []*models.CatechismQuestion
	for rows.Next() {
		question := &models.CatechismQuestion{}
		err := rows.Scan(
			&question.ID,
			&question.QuestionNumber,
			&question.QuestionText,
			&question.AnswerText,
		)
		if err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}
	
	return questions, rows.Err()
}

func (r *CatechismRepository) Create(question *models.CatechismQuestion) error {
	query := `INSERT INTO westminster_catechism (question_number, question_text, answer_text) 
	          VALUES ($1, $2, $3)
	          ON CONFLICT (question_number) 
	          DO UPDATE SET 
	            question_text = EXCLUDED.question_text,
	            answer_text = EXCLUDED.answer_text
	          RETURNING id`
	
	err := database.DB.QueryRow(query,
		question.QuestionNumber,
		question.QuestionText,
		question.AnswerText,
	).Scan(&question.ID)
	
	return err
}

func (r *CatechismRepository) CreateBatch(questions []*models.CatechismQuestion) error {
	for _, question := range questions {
		if err := r.Create(question); err != nil {
			return err
		}
	}
	return nil
}

