package service_test

import (
	"context"
	"errors"
	"testing"
	"user-management-api/internal/core/domain"
	"user-management-api/internal/core/entity"
	"user-management-api/internal/core/mock"
	"user-management-api/internal/core/service"
	"user-management-api/pkg/hash"
)

func TestUserService_Register(t *testing.T) {
	initName := "UserTest"
	initEmail := "user-test@mail.com"
	initPwd := "iniPwd"
	ctx := context.Background()

	t.Run("register success", func(t *testing.T) {
		var userInsert entity.User
		repo := &mock.MockUserRepo{
			ExistsByEmailFn: func(ctx context.Context, email, excludeID string) (bool, error) { return false, nil },
			InsertFn:        func(ctx context.Context, u entity.User) error { userInsert = u; return nil },
		}
		us := service.NewUserService(repo, &mock.MockTokenProvider{})

		id, err := us.Register(ctx, domain.RegisterInput{Name: initName, Email: initEmail, Password: initPwd})

		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if id == "" {
			t.Fatal("expected non-empty id")
		}
		if userInsert.Email != initEmail || userInsert.PasswordHash == initPwd {
			t.Fatalf("insert got unexpected user: %+v", userInsert)
		}
	})

	t.Run("email already exists", func(t *testing.T) {
		isInserted := false
		repo := &mock.MockUserRepo{
			ExistsByEmailFn: func(ctx context.Context, email, excludeID string) (bool, error) { return true, nil },
			InsertFn:        func(ctx context.Context, u entity.User) error { isInserted = true; return nil },
		}

		us := service.NewUserService(repo, &mock.MockTokenProvider{})

		_, err := us.Register(ctx, domain.RegisterInput{Email: initEmail})

		if !errors.Is(err, entity.ErrEmailAlreadyExists) {
			t.Fatalf("expected ErrEmailAlreadyExists, err: %v", err)
		}
		if isInserted {
			t.Fatal("insert should not be called when email exists")
		}
	})
}

func TestUserService_Login(t *testing.T) {
	initEmail := "user1@mail.com"
	initPwd := "pwd123456"

	ctx := context.Background()
	hashPwd, _ := hash.HashPassword(initPwd)

	token := &mock.MockTokenProvider{
		GenerateFn: func(userID string) (string, error) {
			return "token-generate", nil
		},
	}
	t.Run("success", func(t *testing.T) {
		repo := &mock.MockUserRepo{
			GetByEmailFn: func(ctx context.Context, email string) (entity.User, error) {
				return entity.User{ID: "uuid", Email: email, PasswordHash: hashPwd}, nil
			},
		}

		us := service.NewUserService(repo, token)

		token, err := us.Login(ctx, initEmail, initPwd)

		if err != nil {
			t.Fatalf("error: %v", err)
		}

		if token != "token-generate" {
			t.Fatalf("expected token %q, input %q", "token", token)
		}

	})

	t.Run("user not found", func(t *testing.T) {
		repo := &mock.MockUserRepo{
			GetByEmailFn: func(ctx context.Context, email string) (entity.User, error) {
				return entity.User{}, entity.ErrUserNotFound
			},
		}

		us := service.NewUserService(repo, token)

		_, err := us.Login(ctx, initEmail, initPwd)

		if !errors.Is(err, entity.ErrInvalidCredentials) {
			t.Fatalf("expected ErrInvalidCredentials ,err: %v", err)
		}

	})

	t.Run("invalid credentials", func(t *testing.T) {
		initWrongPwd := "wrong-password"
		repo := &mock.MockUserRepo{
			GetByEmailFn: func(ctx context.Context, email string) (entity.User, error) {
				return entity.User{ID: "uuid", Email: email, PasswordHash: hashPwd}, nil
			},
		}
		us := service.NewUserService(repo, token)

		_, err := us.Login(ctx, initEmail, initWrongPwd)

		if !errors.Is(err, entity.ErrInvalidCredentials) {
			t.Fatalf("expected ErrInvalidCredentials ,err: %v", err)
		}
	})
}

func TestUserService_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("get user by id success", func(t *testing.T) {
		userObj := entity.User{ID: "uuid-1", Email: "email@mail.com", PasswordHash: "pwd-hashed-1"}

		repo := &mock.MockUserRepo{
			GetByIDFn: func(ctx context.Context, id string) (entity.User, error) {
				return userObj, nil
			},
		}

		us := service.NewUserService(repo, &mock.MockTokenProvider{})
		user, err := us.Get(ctx, "uuid-1")

		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if user != userObj {
			t.Fatalf("expected %+v, got %+v", userObj, user)
		}
	})

	t.Run("get user by id not found", func(t *testing.T) {
		repo := &mock.MockUserRepo{
			GetByIDFn: func(ctx context.Context, id string) (entity.User, error) {
				return entity.User{}, entity.ErrUserNotFound
			},
		}

		us := service.NewUserService(repo, &mock.MockTokenProvider{})
		_, err := us.Get(ctx, "user-1")

		if !errors.Is(err, entity.ErrUserNotFound) {
			t.Fatalf("expected ErrUserNotFound, err: %v", err)
		}
	})

}

func TestUserService_Update(t *testing.T) {
	ctx := context.Background()

	initNameUpdate := "user updated"
	initEmailUpdate := "update@mail.com"

	t.Run("update name success", func(t *testing.T) {
		userInput := domain.UpdateUserInput{ID: "uuid-1", Name: &initNameUpdate}
		repo := &mock.MockUserRepo{
			UpdateFn: func(ctx context.Context, id string, in domain.UpdateUserInput) (entity.User, error) {
				return entity.User{ID: "uuid-1", Name: initNameUpdate}, nil
			},
		}

		us := service.NewUserService(repo, &mock.MockTokenProvider{})
		user, err := us.Update(ctx, "uuid-1", userInput)

		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if user.Name != *userInput.Name {
			t.Fatalf("expected name is %s", initNameUpdate)
		}
	})

	t.Run("update email success", func(t *testing.T) {
		userInput := domain.UpdateUserInput{ID: "uuid-1", Email: &initEmailUpdate}
		repo := &mock.MockUserRepo{
			ExistsByEmailFn: func(ctx context.Context, email, excludeID string) (bool, error) { return false, nil },
			UpdateFn: func(ctx context.Context, id string, in domain.UpdateUserInput) (entity.User, error) {
				return entity.User{ID: "uuid-1", Email: initEmailUpdate}, nil
			},
		}

		us := service.NewUserService(repo, &mock.MockTokenProvider{})
		user, err := us.Update(ctx, "uuid-1", userInput)

		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if user.Email != *userInput.Email {
			t.Fatalf("expected email is %s", initEmailUpdate)
		}
	})

	t.Run("update with no fields to update", func(t *testing.T) {
		userInput := domain.UpdateUserInput{ID: "uuid-1"}
		isUpdated := false
		repo := &mock.MockUserRepo{
			UpdateFn: func(ctx context.Context, id string, in domain.UpdateUserInput) (entity.User, error) {
				isUpdated = true
				return entity.User{}, nil
			},
		}

		us := service.NewUserService(repo, &mock.MockTokenProvider{})
		_, err := us.Update(ctx, "uuid-1", userInput)

		if !errors.Is(err, entity.ErrNoFieldsToUpdate) {
			t.Fatalf("expected ErrNoFieldsToUpdate, err: %v", err)
		}
		if isUpdated {
			t.Fatal("update should not be called when there are no fields to update")
		}
	})

	t.Run("update with email already exists", func(t *testing.T) {
		userInput := domain.UpdateUserInput{ID: "uuid-1", Email: &initEmailUpdate}
		isUpdated := false
		repo := &mock.MockUserRepo{
			ExistsByEmailFn: func(ctx context.Context, email, excludeID string) (bool, error) { return true, nil },
			UpdateFn: func(ctx context.Context, id string, in domain.UpdateUserInput) (entity.User, error) {
				isUpdated = true
				return entity.User{}, nil
			},
		}

		us := service.NewUserService(repo, &mock.MockTokenProvider{})
		_, err := us.Update(ctx, "uuid-1", userInput)

		if !errors.Is(err, entity.ErrEmailAlreadyExists) {
			t.Fatalf("expected ErrEmailAlreadyExists, err: %v", err)
		}
		if isUpdated {
			t.Fatal("update should not be called when email already exists")
		}
	})

	t.Run("existsByEmail returns an error", func(t *testing.T) {
		wantErr := errors.New("db error")
		userInput := domain.UpdateUserInput{ID: "uuid-1", Email: &initEmailUpdate}
		isUpdated := false
		repo := &mock.MockUserRepo{
			ExistsByEmailFn: func(ctx context.Context, email, excludeID string) (bool, error) { return false, wantErr },
			UpdateFn: func(ctx context.Context, id string, in domain.UpdateUserInput) (entity.User, error) {
				isUpdated = true
				return entity.User{}, nil
			},
		}

		us := service.NewUserService(repo, &mock.MockTokenProvider{})
		_, err := us.Update(ctx, "uuid-1", userInput)

		if !errors.Is(err, wantErr) {
			t.Fatalf("expected %v, err: %v", wantErr, err)
		}
		if isUpdated {
			t.Fatal("update should not be called when existsByEmail fails")
		}
	})

	t.Run("update with user not found", func(t *testing.T) {
		userInput := domain.UpdateUserInput{ID: "uuid-1", Name: &initNameUpdate}
		repo := &mock.MockUserRepo{
			UpdateFn: func(ctx context.Context, id string, in domain.UpdateUserInput) (entity.User, error) {
				return entity.User{}, entity.ErrUserNotFound
			},
		}

		us := service.NewUserService(repo, &mock.MockTokenProvider{})
		_, err := us.Update(ctx, "uuid-1", userInput)

		if !errors.Is(err, entity.ErrUserNotFound) {
			t.Fatalf("expected ErrUserNotFound, err: %v", err)
		}
	})
}

func TestUserService_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("delete success", func(t *testing.T) {
		repo := &mock.MockUserRepo{
			DeleteFn: func(ctx context.Context, id string) error { return nil },
		}

		us := service.NewUserService(repo, &mock.MockTokenProvider{})

		err := us.Delete(ctx, "uuid-1")

		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
	})

	t.Run("delete with user not found", func(t *testing.T) {
		repo := &mock.MockUserRepo{
			DeleteFn: func(ctx context.Context, id string) error { return entity.ErrUserNotFound },
		}

		us := service.NewUserService(repo, &mock.MockTokenProvider{})

		err := us.Delete(ctx, "uuid-0")

		if !errors.Is(err, entity.ErrUserNotFound) {
			t.Fatalf("expected ErrUserNotFound, err: %v", err)
		}
	})
}

func TestUserService_GetList(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		usersResult := []entity.User{
			{ID: "uuid-1", Email: "user1@mail.com"},
			{ID: "uuid-2", Email: "user2@mail.com"},
		}
		repo := &mock.MockUserRepo{
			ListFn: func(ctx context.Context, f domain.ListUsersFilter) ([]entity.User, error) {
				return usersResult, nil
			},
		}

		us := service.NewUserService(repo, &mock.MockTokenProvider{})
		filter := domain.ListUsersFilter{}
		users, err := us.GetList(ctx, &filter)

		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if len(users) != len(usersResult) {
			t.Fatalf("expected %d users", len(usersResult))
		}
	})

	t.Run("default limit and date range when empty", func(t *testing.T) {
		var gotFilter domain.ListUsersFilter
		repo := &mock.MockUserRepo{
			ListFn: func(ctx context.Context, f domain.ListUsersFilter) ([]entity.User, error) {
				gotFilter = f
				return nil, nil
			},
		}

		us := service.NewUserService(repo, &mock.MockTokenProvider{})
		filter := domain.ListUsersFilter{}
		_, err := us.GetList(ctx, &filter)

		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if gotFilter.Limit != 20 {
			t.Fatalf("expected default limit 20, got %d", gotFilter.Limit)
		}
		if gotFilter.To.IsZero() || gotFilter.From.IsZero() {
			t.Fatal("expected From and To to be default")
		}
		if !gotFilter.From.Before(gotFilter.To) {
			t.Fatalf("expected From (%v) before To (%v)", gotFilter.From, gotFilter.To)
		}
	})

	t.Run("get user results empty", func(t *testing.T) {
		repo := &mock.MockUserRepo{
			ListFn: func(ctx context.Context, f domain.ListUsersFilter) ([]entity.User, error) {
				return []entity.User{}, nil
			},
		}

		us := service.NewUserService(repo, &mock.MockTokenProvider{})
		filter := domain.ListUsersFilter{}
		users, err := us.GetList(ctx, &filter)

		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if len(users) != 0 {
			t.Fatal("expected empty result")
		}
	})
}

func TestUserService_CountUsers(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &mock.MockUserRepo{
			CountFn: func(ctx context.Context) (int64, error) { return 42, nil },
		}

		us := service.NewUserService(repo, &mock.MockTokenProvider{})
		got, err := us.CountUsers(ctx)

		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if got != 42 {
			t.Fatal("expected 42")
		}
	})
}
