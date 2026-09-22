package port

import (
	"context"
	"user-management-api/internal/core/domain"
	"user-management-api/internal/core/entity"
)

type UserService interface {
	Register(ctx context.Context, user domain.RegisterInput) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
	Get(ctx context.Context, userID string) (entity.User, error)
	GetList(ctx context.Context, filter *domain.ListUsersFilter) ([]entity.User, error)
	Update(ctx context.Context, userID string, user domain.UpdateUserInput) (entity.User, error)
	Delete(ctx context.Context, userID string) error
	CountUsers(ctx context.Context) (int64, error)
}
