package port

import (
	"context"
	"user-management-api/internal/core/domain"
	"user-management-api/internal/core/entity"
)

type UserRepository interface {
	Insert(ctx context.Context, user entity.User) error
	GetByID(ctx context.Context, id string) (entity.User, error)
	List(ctx context.Context, filter domain.ListUsersFilter) ([]entity.User, error)
	Update(ctx context.Context, id string, user domain.UpdateUserInput) (entity.User, error)
	Delete(ctx context.Context, id string) error
	ExistsByEmail(ctx context.Context, email, excludeID string) (bool, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
	Count(ctx context.Context) (int64, error)
}
