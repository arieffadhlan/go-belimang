package handlers

import (
	"belimang/internal/dto"
	"belimang/internal/entities"
	"belimang/internal/services"
	"belimang/internal/utils"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type MerchantHandler struct {
	service    services.MerchantService
	validation *validator.Validate
}

func NewMerchantHandler(service services.MerchantService, validation *validator.Validate) MerchantHandler {
	return MerchantHandler{
		service:    service,
		validation: validation,
	}
}

func (h MerchantHandler) CreateMerchant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := dto.CreateMerchantRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validation.Struct(req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	merchant := entities.Merchant{
		ID:       uuid.New().String(),
		Name:     req.Name,
		ImageURL: req.ImageURL,
		Category: req.MerchantCategory,
		Location: entities.Location{Lat: req.Location.Lat, Long: req.Location.Long},
	}

	merchantId, err := h.service.CreateMerchant(ctx, merchant)
	if err != nil {
		if appErr, ok := err.(utils.AppError); ok {
			utils.SendErrorResponse(w, appErr.StatusCode, appErr.Message)
		} else {
			utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.SendResponse(w, http.StatusCreated, merchantId)
}

func (h MerchantHandler) GetAllMerchant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()

	limit := 5
	if limStr := q.Get("limit"); limStr != "" {
		if limVal, err := strconv.Atoi(limStr); err == nil && limVal > 0 {
			 limit = limVal
		}
	}

	offset := 0
	if offStr := q.Get("offset"); offStr != "" {
		if offVal, err := strconv.Atoi(offStr); err == nil && offVal > 0 {
			 offset = offVal
		}
	}

	createdAt := strings.ToLower(q.Get("createdAt"))
	if createdAt != "" && createdAt != "asc" && createdAt != "desc" {
		 createdAt = ""
	}

	filter := entities.MerchantFilter{
		Limit:            limit,
		CreatedAt:        createdAt,
		Name:             q.Get("name"),
		MerchantID:       q.Get("merchantId"),
		MerchantCategory: q.Get("merchantCategory"),
		Offset:           offset,
	}

	merchant, err := h.service.GetAllMerchant(ctx, filter)
	if err != nil {
		if appErr, ok := err.(utils.AppError); ok {
			utils.SendErrorResponse(w, appErr.StatusCode, appErr.Message)
		} else {
			utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.SendResponse(w, http.StatusOK, merchant)
}

func (h MerchantHandler) CreateMercItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := dto.CreateMercItemRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validation.Struct(req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	merchantId := chi.URLParam(r, "merchantId")

	item := entities.MerchantItem{
		ID:       	uuid.New().String(),
		MerchantID: merchantId,
		Name:       req.Name,
		Category:   req.ProductCategory,
		ImageURL:   req.ImageURL,
		Price:      req.Price,
	}

	merchantItemId, err := h.service.CreateMercItem(ctx, item)
	if err != nil {
		if appErr, ok := err.(utils.AppError); ok {
			utils.SendErrorResponse(w, appErr.StatusCode, appErr.Message)
		} else {
			utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.SendResponse(w, http.StatusCreated, merchantItemId)
}

func (h MerchantHandler) GetAllMercItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()

	limit := 5
	if limStr := q.Get("limit"); limStr != "" {
		if limVal, err := strconv.Atoi(limStr); err == nil && limVal > 0 {
			 limit = limVal
		}
	}

	offset := 0
	if offStr := q.Get("offset"); offStr != "" {
		if offVal, err := strconv.Atoi(offStr); err == nil && offVal > 0 {
			 offset = offVal
		}
	}

	createdAt := strings.ToLower(q.Get("createdAt"))
	if createdAt != "" && createdAt != "asc" && createdAt != "desc" {
		 createdAt = ""
	}

	filter := entities.MerchantItemFilter{
		Limit:           limit,
		CreatedAt:       createdAt,
		Name:            q.Get("name"),
		ItemID:          q.Get("itemId"),
		ProductCategory: q.Get("productCategory"),
		Offset:          offset,
	}

	merchantId := chi.URLParam(r, "merchantId")

	merchantItems, err := h.service.GetAllMercItem(ctx, merchantId, filter)
	if err != nil {
		if appErr, ok := err.(utils.AppError); ok {
			utils.SendErrorResponse(w, appErr.StatusCode, appErr.Message)
		} else {
			utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.SendResponse(w, http.StatusOK, merchantItems)
}
