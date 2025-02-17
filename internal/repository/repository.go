package repository

import (
	"context"
	db "github.com/LuccChagas/my-chat-app/db/sqlc"
	"github.com/google/uuid"
)

type Repository struct {
	dbtx    db.DBTX
	queries *db.Queries
}

func NewRepository(dbtx db.DBTX, q *db.Queries) *Repository {
	return &Repository{
		dbtx:    dbtx,
		queries: q,
	}
}

type RepositoryInterface interface {
	CreateUser(ctx context.Context, user db.CreateUsersParams) (db.User, error)
	GetAllUsers(ctx context.Context) ([]db.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (db.User, error)
	GetUserByNickname(ctx context.Context, nickname string) (db.User, error)
}
