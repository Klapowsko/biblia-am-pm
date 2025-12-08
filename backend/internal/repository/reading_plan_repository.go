package repository

import (
	"biblia-am-pm/internal/database"
	"biblia-am-pm/internal/models"
	"database/sql"
)

type ReadingPlanRepository struct{}

func NewReadingPlanRepository() *ReadingPlanRepository {
	return &ReadingPlanRepository{}
}

func (r *ReadingPlanRepository) GetByDayOfYear(dayOfYear int) (*models.ReadingPlan, error) {
	query := `SELECT id, day_of_year, old_testament_ref, new_testament_ref, psalms_ref, proverbs_ref 
	          FROM reading_plans WHERE day_of_year = $1`
	
	plan := &models.ReadingPlan{}
	err := database.DB.QueryRow(query, dayOfYear).Scan(
		&plan.ID,
		&plan.DayOfYear,
		&plan.OldTestamentRef,
		&plan.NewTestamentRef,
		&plan.PsalmsRef,
		&plan.ProverbsRef,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	if err != nil {
		return nil, err
	}
	
	return plan, nil
}

func (r *ReadingPlanRepository) Create(plan *models.ReadingPlan) error {
	query := `INSERT INTO reading_plans (day_of_year, old_testament_ref, new_testament_ref, psalms_ref, proverbs_ref) 
	          VALUES ($1, $2, $3, $4, $5)
	          ON CONFLICT (day_of_year) 
	          DO UPDATE SET 
	            old_testament_ref = EXCLUDED.old_testament_ref,
	            new_testament_ref = EXCLUDED.new_testament_ref,
	            psalms_ref = EXCLUDED.psalms_ref,
	            proverbs_ref = EXCLUDED.proverbs_ref
	          RETURNING id`
	
	err := database.DB.QueryRow(query,
		plan.DayOfYear,
		plan.OldTestamentRef,
		plan.NewTestamentRef,
		plan.PsalmsRef,
		plan.ProverbsRef,
	).Scan(&plan.ID)
	
	return err
}

func (r *ReadingPlanRepository) GetAll() ([]*models.ReadingPlan, error) {
	query := `SELECT id, day_of_year, old_testament_ref, new_testament_ref, psalms_ref, proverbs_ref 
	          FROM reading_plans ORDER BY day_of_year`
	
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var plans []*models.ReadingPlan
	for rows.Next() {
		plan := &models.ReadingPlan{}
		err := rows.Scan(
			&plan.ID,
			&plan.DayOfYear,
			&plan.OldTestamentRef,
			&plan.NewTestamentRef,
			&plan.PsalmsRef,
			&plan.ProverbsRef,
		)
		if err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}
	
	return plans, rows.Err()
}

