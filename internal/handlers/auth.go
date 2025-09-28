package handlers

import (
	"belimang/internal/dto"
	"belimang/internal/entities"
	"belimang/internal/services"
	"belimang/internal/utils"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	service services.AuthService
	validation *validator.Validate
}

func NewAuthHandler(service services.AuthService, validation *validator.Validate) AuthHandler {
	return AuthHandler{
		service: service,
		validation: validation,
	}
}

func (h AuthHandler) SignUp(w http.ResponseWriter, r *http.Request, role string) {
	ctx := r.Context()
	req := dto.SignUpRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validation.Struct(req); err != nil {
		SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	usr := entities.User{
		Email: 	req.Email,
		IsAdmin:  role == "admin",
		Username: req.Username,
		Password: req.Password,
	}

	t, err := h.service.SignUp(ctx, usr)
	if err != nil {
		if appErr, ok := err.(utils.AppError); ok {
			SendErrorResponse(w, appErr.StatusCode, appErr.Message)
		} else {
			SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		}

		return
	}

	SendResponse(w, http.StatusCreated, t)
}

func (h AuthHandler) SignIn(w http.ResponseWriter, r *http.Request, role string) {
	ctx := r.Context()
	req := dto.SignInRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validation.Struct(req); err != nil {
		SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	usr := entities.User{
		IsAdmin:  role == "admin",
		Username: req.Username,
		Password: req.Password,
	}

	t, err := h.service.SignIn(ctx, usr)
	if err != nil {
		if appErr, ok := err.(utils.AppError); ok {
			SendErrorResponse(w, appErr.StatusCode, appErr.Message)
		} else {
			SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		}

		return
	}

	SendResponse(w, http.StatusOK, t)
}
