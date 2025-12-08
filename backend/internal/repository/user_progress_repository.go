package repository

import (
	"biblia-am-pm/internal/database"
	"biblia-am-pm/internal/models"
	"database/sql"
	"time"
)

type UserProgressRepository struct{}

func NewUserProgressRepository() *UserProgressRepository {
	return &UserProgressRepository{}
}

func (r *UserProgressRepository) GetByUserAndDate(userID int, date time.Time) (*models.UserProgress, error) {
	query := `SELECT id, user_id, reading_plan_id, date, morning_completed, evening_completed, completed_at 
	          FROM user_progress WHERE user_id = $1 AND date = $2`
	
	progress := &models.UserProgress{}
	var completedAt sql.NullTime
	
	err := database.DB.QueryRow(query, userID, date.Format("2006-01-02")).Scan(
		&progress.ID,
		&progress.UserID,
		&progress.ReadingPlanID,
		&progress.Date,
		&progress.MorningCompleted,
		&progress.EveningCompleted,
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

func (r *UserProgressRepository) CreateOrUpdate(progress *models.UserProgress) error {
	query := `INSERT INTO user_progress (user_id, reading_plan_id, date, morning_completed, evening_completed, completed_at)
	          VALUES ($1, $2, $3, $4, $5, $6)
	          ON CONFLICT (user_id, reading_plan_id, date)
	          DO UPDATE SET 
	            morning_completed = EXCLUDED.morning_completed,
	            evening_completed = EXCLUDED.evening_completed,
	            completed_at = EXCLUDED.completed_at
	          RETURNING id`
	
	now := time.Now()
	var completedAt *time.Time
	if progress.MorningCompleted && progress.EveningCompleted {
		completedAt = &now
	}
	
	err := database.DB.QueryRow(query,
		progress.UserID,
		progress.ReadingPlanID,
		progress.Date.Format("2006-01-02"),
		progress.MorningCompleted,
		progress.EveningCompleted,
		completedAt,
	).Scan(&progress.ID)
	
	return err
}

func (r *UserProgressRepository) GetUserProgress(userID int) ([]*models.UserProgress, error) {
	query := `SELECT id, user_id, reading_plan_id, date, morning_completed, evening_completed, completed_at 
	          FROM user_progress WHERE user_id = $1 ORDER BY date DESC`
	
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var progresses []*models.UserProgress
	for rows.Next() {
		progress := &models.UserProgress{}
		var completedAt sql.NullTime
		
		err := rows.Scan(
			&progress.ID,
			&progress.UserID,
			&progress.ReadingPlanID,
			&progress.Date,
			&progress.MorningCompleted,
			&progress.EveningCompleted,
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

