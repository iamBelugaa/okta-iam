package users_handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/iamNilotpal/iam/internal/models"
	"github.com/iamNilotpal/iam/internal/okta"
	"github.com/iamNilotpal/iam/pkg/response"
	"go.uber.org/zap"
)

type handler struct {
	okta *okta.Service
	log  *zap.SugaredLogger
}

func New(okta *okta.Service, log *zap.SugaredLogger) *handler {
	return &handler{okta: okta, log: log}
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, response.NewAPIError("Invalid request payload", nil))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	user, err := h.okta.CreateUser(ctx, &req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, response.NewAPIError("Failed to create user", nil))
		return
	}

	response.Success(w, http.StatusCreated, "Success", user)
}

func (h *handler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		response.Error(w, http.StatusBadRequest, response.NewAPIError("User ID is required", nil))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	user, err := h.okta.GetUser(ctx, userID)
	if err != nil {
		response.Error(w, http.StatusNotFound, response.NewAPIError("Failed to retrieve user", nil))
		return
	}

	response.Success(w, http.StatusOK, "Success", user)
}

func (h *handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	users, err := h.okta.ListUsers(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, response.NewAPIError("Failed to list users", nil))
		return
	}

	response.Success(w, http.StatusOK, "Success", users)
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		response.Error(w, http.StatusBadRequest, response.NewAPIError("User ID is required", nil))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	err := h.okta.DeactivateUser(ctx, userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, response.NewAPIError("Failed to deactivate user", nil))
		return
	}

	err = h.okta.DeleteUser(ctx, userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, response.NewAPIError("Failed to delete user after deactivation", nil))
		return
	}

	response.Success(w, http.StatusOK, "Success", nil)
}
