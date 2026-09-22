package client

import (
	"time"
	"user-management-api/internal/core/entity"
)

type RegisterResponse struct {
	ID string `json:"id"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func ToUserResponse(u entity.User) UserResponse {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		loc = time.UTC
	}
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.In(loc),
	}
}

func ToUserResponseList(users []entity.User) []UserResponse {
	res := make([]UserResponse, len(users))
	for i, u := range users {
		res[i] = ToUserResponse(u)
	}
	return res
}
