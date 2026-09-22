package mock

import (
	"context"
	"user-management-api/internal/core/domain"
	"user-management-api/internal/core/entity"
)

type MockUserRepo struct {
	InsertFn        func(ctx context.Context, u entity.User) error
	GetByIDFn       func(ctx context.Context, id string) (entity.User, error)
	ListFn          func(ctx context.Context, f domain.ListUsersFilter) ([]entity.User, error)
	UpdateFn        func(ctx context.Context, id string, in domain.UpdateUserInput) (entity.User, error)
	DeleteFn        func(ctx context.Context, id string) error
	ExistsByEmailFn func(ctx context.Context, email, excludeID string) (bool, error)
	GetByEmailFn    func(ctx context.Context, email string) (entity.User, error)
	CountFn         func(ctx context.Context) (int64, error)
}

func (m *MockUserRepo) Insert(ctx context.Context, u entity.User) error {
	return m.InsertFn(ctx, u)
}

func (m *MockUserRepo) GetByID(ctx context.Context, id string) (entity.User, error) {
	return m.GetByIDFn(ctx, id)
}

func (m *MockUserRepo) List(ctx context.Context, filter domain.ListUsersFilter) ([]entity.User, error) {
	return m.ListFn(ctx, filter)
}

func (m *MockUserRepo) Update(ctx context.Context, id string, in domain.UpdateUserInput) (entity.User, error) {
	return m.UpdateFn(ctx, id, in)
}

func (m *MockUserRepo) Delete(ctx context.Context, id string) error {
	return m.DeleteFn(ctx, id)
}

func (m *MockUserRepo) ExistsByEmail(ctx context.Context, email, excludeID string) (bool, error) {
	return m.ExistsByEmailFn(ctx, email, excludeID)
}

func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	return m.GetByEmailFn(ctx, email)
}

func (m *MockUserRepo) Count(ctx context.Context) (int64, error) {
	return m.CountFn(ctx)
}

type MockTokenProvider struct {
	GenerateFn func(userID string) (string, error)
	ValidateFn func(token string) (string, error)
}

func (m *MockTokenProvider) GenerateToken(userID string) (string, error) {
	return m.GenerateFn(userID)
}

func (m *MockTokenProvider) ValidateToken(token string) (string, error) {
	return m.ValidateFn(token)
}
