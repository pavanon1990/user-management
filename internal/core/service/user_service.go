package service

import (
	"context"
	"errors"
	"log"
	"time"
	"user-management-api/internal/core/domain"
	"user-management-api/internal/core/entity"
	"user-management-api/internal/core/port"
	"user-management-api/pkg/hash"
	"user-management-api/pkg/uuid"
)

type userService struct {
	userRepo      port.UserRepository
	tokenProvider port.TokenProvider
}

func NewUserService(userRepo port.UserRepository, tokenProvider port.TokenProvider) port.UserService {
	return &userService{userRepo: userRepo, tokenProvider: tokenProvider}
}

func (u *userService) Register(ctx context.Context, user domain.RegisterInput) (string, error) {
	exists, err := u.userRepo.ExistsByEmail(ctx, user.Email, "")
	if err != nil {
		return "", err
	}
	if exists {
		return "", entity.ErrEmailAlreadyExists
	}

	hashedPwd, err := hash.HashPassword(user.Password)
	if err != nil {
		return "", err
	}
	newUser := entity.User{
		ID:           uuid.NewUUID(),
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: hashedPwd,
		CreatedAt:    time.Now().UTC(),
	}

	if err := u.userRepo.Insert(ctx, newUser); err != nil {
		return "", err
	}

	return newUser.ID, nil
}

func (u *userService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			return "", entity.ErrUserNotFound
		}
		return "", err
	}

	if err := hash.ComparePassword(user.PasswordHash, password); err != nil {
		return "", entity.ErrInvalidCredentials
	}

	return u.tokenProvider.GenerateToken(user.ID)
}

func (u *userService) Get(ctx context.Context, userID string) (entity.User, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		log.Printf("Get UserByID err: %v", err)
		return entity.User{}, err
	}
	return user, nil
}

func (u *userService) GetList(ctx context.Context, filter *domain.ListUsersFilter) ([]entity.User, error) {
	if filter.To.IsZero() {
		filter.To = time.Now()
	}
	if filter.From.IsZero() {
		filter.From = filter.To.AddDate(0, -1, 0)
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	return u.userRepo.List(ctx, *filter)
}

func (u *userService) Update(ctx context.Context, id string, input domain.UpdateUserInput) (entity.User, error) {

	if input.Name == nil && input.Email == nil {
		return entity.User{}, entity.ErrNoFieldsToUpdate
	}

	if input.Email != nil {
		exist, err := u.userRepo.ExistsByEmail(ctx, *input.Email, id)
		if err != nil {
			return entity.User{}, err
		}
		if exist {
			return entity.User{}, entity.ErrEmailAlreadyExists
		}
	}

	return u.userRepo.Update(ctx, id, input)

}

func (u *userService) Delete(ctx context.Context, userID string) error {
	if err := u.userRepo.Delete(ctx, userID); err != nil {
		return err
	}
	return nil
}

func (u *userService) CountUsers(ctx context.Context) (int64, error) {
	return u.userRepo.Count(ctx)
}
