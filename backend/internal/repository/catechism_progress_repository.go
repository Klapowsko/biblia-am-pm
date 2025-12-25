package repository

import (
	"biblia-am-pm/internal/database"
	"biblia-am-pm/internal/models"
	"database/sql"
	"time"
)

type CatechismProgressRepository struct{}

func NewCatechismProgressRepository() *CatechismProgressRepository {
	return &CatechismProgressRepository{}
}

func (r *CatechismProgressRepository) GetByUserAndDate(userID int, questionID int, date time.Time) (*models.CatechismProgress, error) {
	query := `SELECT id, user_id, question_id, date, completed, completed_at 
	          FROM catechism_progress WHERE user_id = $1 AND question_id = $2 AND date = $3`
	
	progress := &models.CatechismProgress{}
	var completedAt sql.NullTime
	
	err := database.DB.QueryRow(query, userID, questionID, date.Format("2006-01-02")).Scan(
		&progress.ID,
		&progress.UserID,
		&progress.QuestionID,
		&progress.Date,
		&progress.Completed,
		&completedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	if err != nil {
		return nil, err
	}
	
	if completedAt.Valid {
		progress.CompletedAt = &completedAt.Time
	}
	
	return progress, nil
}

func (r *CatechismProgressRepository) GetByUserAndQuestionForWeek(userID int, questionID int, weekStart time.Time) ([]*models.CatechismProgress, error) {
	weekEnd := weekStart.AddDate(0, 0, 6) // 6 days after start (7 days total)
	query := `SELECT id, user_id, question_id, date, completed, completed_at 
	          FROM catechism_progress 
	          WHERE user_id = $1 AND question_id = $2 
	          AND date >= $3 AND date <= $4 
	          ORDER BY date`
	
	rows, err := database.DB.Query(query, userID, questionID, weekStart.Format("2006-01-02"), weekEnd.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var progresses []*models.CatechismProgress
	for rows.Next() {
		progress := &models.CatechismProgress{}
		var completedAt sql.NullTime
		
		err := rows.Scan(
			&progress.ID,
			&progress.UserID,
			&progress.QuestionID,
			&progress.Date,
			&progress.Completed,
			&completedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if completedAt.Valid {
			progress.CompletedAt = &completedAt.Time
		}
		
		progresses = append(progresses, progress)
	}
	
	return progresses, rows.Err()
}

func (r *CatechismProgressRepository) CreateOrUpdate(progress *models.CatechismProgress) error {
	query := `INSERT INTO catechism_progress (user_id, question_id, date, completed, completed_at)
	          VALUES ($1, $2, $3, $4, $5)
	          ON CONFLICT (user_id, question_id, date)
	          DO UPDATE SET 
	            completed = EXCLUDED.completed,
	            completed_at = EXCLUDED.completed_at
	          RETURNING id`
	
	var completedAt *time.Time
	if progress.Completed {
		now := time.Now()
		completedAt = &now
	}
	
	err := database.DB.QueryRow(query,
		progress.UserID,
		progress.QuestionID,
		progress.Date.Format("2006-01-02"),
		progress.Completed,
		completedAt,
	).Scan(&progress.ID)
	
	return err
}

func (r *CatechismProgressRepository) GetUserProgress(userID int) ([]*models.CatechismProgress, error) {
	query := `SELECT id, user_id, question_id, date, completed, completed_at 
	          FROM catechism_progress WHERE user_id = $1 ORDER BY date DESC`
	
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var progresses []*models.CatechismProgress
	for rows.Next() {
		progress := &models.CatechismProgress{}
		var completedAt sql.NullTime
		
		err := rows.Scan(
			&progress.ID,
			&progress.UserID,
			&progress.QuestionID,
			&progress.Date,
			&progress.Completed,
			&completedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if completedAt.Valid {
			progress.CompletedAt = &completedAt.Time
		}
		
		progresses = append(progresses, progress)
	}
	
	return progresses, rows.Err()
}

