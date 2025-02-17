package repository

import (
	"context"
	db "github.com/LuccChagas/my-chat-app/db/sqlc"
	"github.com/google/uuid"
)

func (r *Repository) CreateUser(ctx context.Context, user db.CreateUsersParams) (db.User, error) {
	u, err := r.queries.CreateUsers(ctx, user)
	if err != nil {
		return db.User{}, err
	}

	return u, nil
}
func (r *Repository) GetAllUsers(ctx context.Context) ([]db.User, error) {
	users, err := r.queries.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}
func (r *Repository) GetUser(ctx context.Context, id uuid.UUID) (db.User, error) {
	u, err := r.queries.GetUser(ctx, id)
	if err != nil {
		return db.User{}, err
	}

	return u, nil
}

func (r *Repository) GetUserByNickname(ctx context.Context, nickname string) (db.User, error) {
	u, err := r.queries.GetUserByNickname(ctx, nickname)
	if err != nil {
		return db.User{}, err
	}

	return u, nil
}
