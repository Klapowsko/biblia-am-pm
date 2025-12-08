package repository

import (
	"biblia-am-pm/internal/database"
	"biblia-am-pm/internal/models"
	"database/sql"
	"time"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) CreateUser(email, hashedPassword string) (*models.User, error) {
	query := `INSERT INTO users (email, password, created_at) VALUES ($1, $2, $3) RETURNING id, email, created_at`
	
	user := &models.User{}
	err := database.DB.QueryRow(query, email, hashedPassword, time.Now()).Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, password, created_at FROM users WHERE email = $1`
	
	user := &models.User{}
	err := database.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	query := `SELECT id, email, created_at FROM users WHERE id = $1`
	
	user := &models.User{}
	err := database.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

