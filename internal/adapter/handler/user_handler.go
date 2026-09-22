package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"user-management-api/internal/adapter/client"
	"user-management-api/internal/constant"
	"user-management-api/internal/core/domain"
	"user-management-api/internal/core/entity"
	"user-management-api/internal/core/port"
	"user-management-api/pkg/validator"
)

type UserHandler struct {
	userService port.UserService
}

func NewUserHandler(userService port.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	user, err := h.userService.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			client.WriteFail(w, http.StatusNotFound, constant.NOT_FOUND_CODE, err.Error())
			return
		}
		client.WriteFail(w, http.StatusInternalServerError, constant.INTERNAL_ERROR_CODE, err.Error())
		return
	}

	client.WriteSuccess(w, http.StatusOK, constant.SUCCESS_CODE, client.ToUserResponse(user))
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req client.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		client.WriteFail(w, http.StatusBadRequest, constant.BAD_REQUEST_CODE, constant.INVALID_REQUEST_BODY_MSG)
		return
	}

	if err := validator.ValidateStruct(req); err != nil {
		client.WriteFail(w, http.StatusBadRequest, constant.BAD_REQUEST_CODE, err.Error())
		return
	}

	id, err := h.userService.Register(r.Context(), req.ToRegisterInput())
	if err != nil {
		if errors.Is(err, entity.ErrEmailAlreadyExists) {
			client.WriteFail(w, http.StatusBadRequest, constant.BAD_REQUEST_CODE, err.Error())
			return
		}
		client.WriteFail(w, http.StatusInternalServerError, constant.INTERNAL_ERROR_CODE, err.Error())
		return
	}

	client.WriteSuccess(w, http.StatusCreated, constant.CREATED_CODE, client.RegisterResponse{ID: id})
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	filter := resovleQuery(r.URL.Query())
	users, err := h.userService.GetList(r.Context(), &filter)
	if err != nil {
		client.WriteFail(w, http.StatusInternalServerError, constant.INTERNAL_ERROR_CODE, err.Error())
		return
	}
	client.WriteSuccessList(w, http.StatusOK, constant.SUCCESS_CODE, client.ToUserResponseList(users), client.Meta{
		Offset: filter.Offset,
		Limit:  filter.Limit,
	})

}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.userService.Delete(r.Context(), id); err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			client.WriteFail(w, http.StatusNotFound, constant.NOT_FOUND_CODE, err.Error())
			return
		}
		client.WriteFail(w, http.StatusInternalServerError, constant.INTERNAL_ERROR_CODE, err.Error())
		return
	}
	client.WriteMessage(w, http.StatusOK, constant.SUCCESS_CODE, "user deleted successfully")
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req client.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		client.WriteFail(w, http.StatusBadRequest, constant.BAD_REQUEST_CODE, constant.INVALID_REQUEST_BODY_MSG)
		return
	}

	if err := validator.ValidateStruct(req); err != nil {
		client.WriteFail(w, http.StatusBadRequest, constant.BAD_REQUEST_CODE, err.Error())
		return
	}

	user, err := h.userService.Update(r.Context(), id, req.ToUpdateUserInput())
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrUserNotFound):
			client.WriteFail(w, http.StatusNotFound, constant.NOT_FOUND_CODE, err.Error())
		case errors.Is(err, entity.ErrEmailAlreadyExists), errors.Is(err, entity.ErrNoFieldsToUpdate):
			client.WriteFail(w, http.StatusBadRequest, constant.BAD_REQUEST_CODE, err.Error())
		default:
			client.WriteFail(w, http.StatusInternalServerError, constant.INTERNAL_ERROR_CODE, err.Error())
		}
		return
	}

	client.WriteSuccess(w, http.StatusOK, constant.SUCCESS_CODE, client.ToUserResponse(user))
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req client.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		client.WriteFail(w, http.StatusBadRequest, constant.BAD_REQUEST_CODE, constant.INVALID_REQUEST_BODY_MSG)
		return
	}

	if err := validator.ValidateStruct(req); err != nil {
		client.WriteFail(w, http.StatusBadRequest, constant.BAD_REQUEST_CODE, err.Error())
		return
	}

	token, err := h.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, entity.ErrInvalidCredentials) {
			client.WriteFail(w, http.StatusUnauthorized, constant.UNAUTHORIZED_CODE, err.Error())
			return
		}
		client.WriteFail(w, http.StatusInternalServerError, constant.INTERNAL_ERROR_CODE, err.Error())
		return
	}

	client.WriteSuccess(w, http.StatusOK, constant.SUCCESS_CODE, client.LoginResponse{Token: token})
}

func resovleQuery(query url.Values) domain.ListUsersFilter {
	filter := domain.ListUsersFilter{}

	if limitStr := query.Get("limit"); limitStr != "" {
		if limit, err := strconv.ParseInt(limitStr, 10, 64); err == nil {
			filter.Limit = limit
		}
	}
	if offsetStr := query.Get("offset"); offsetStr != "" {
		if offset, err := strconv.ParseInt(offsetStr, 10, 64); err == nil {
			filter.Offset = offset
		}
	}
	if fromStr := query.Get("from"); fromStr != "" {
		if from, err := time.Parse(time.RFC3339, fromStr); err == nil {
			filter.From = from
		}
	}
	if toStr := query.Get("to"); toStr != "" {
		if to, err := time.Parse(time.RFC3339, toStr); err == nil {
			filter.To = to
		}
	}

	return filter

}
