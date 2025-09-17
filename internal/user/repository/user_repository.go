package repository

import (
	"context"
	"database/sql"

	"go-clean-arch/internal/user/domain/entity"
	"go-clean-arch/internal/user/domain/repository"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (id, email, password, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Email, user.Password, user.Name, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	query := `SELECT id, email, password, name, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.GetContext(ctx, &user, query, id)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &user, err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	query := `SELECT id, email, password, name, created_at, updated_at FROM users WHERE email = $1`
	err := r.db.GetContext(ctx, &user, query, email)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &user, err
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users 
		SET email = $2, name = $3, updated_at = $4
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Email, user.Name, user.UpdatedAt)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	var users []*entity.User
	query := `SELECT id, email, password, name, created_at, updated_at FROM users LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &users, query, limit, offset)
	return users, err
}
