package client

import "user-management-api/internal/core/domain"

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty" validate:"omitempty"`
	Email *string `json:"email,omitempty" validate:"omitempty,email"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (r RegisterRequest) ToRegisterInput() domain.RegisterInput {
	return domain.RegisterInput{
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
	}
}

func (r UpdateUserRequest) ToUpdateUserInput() domain.UpdateUserInput {
	return domain.UpdateUserInput{
		Name:  r.Name,
		Email: r.Email,
	}
}
