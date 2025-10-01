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

	"github.com/go-playground/validator/v10"
)

type MerchantHandler struct {
	service    services.MerchantService
	validation *validator.Validate
}

func NewMerchantHandler(service services.MerchantService) MerchantHandler {
	return MerchantHandler{
		service: service,
	}
}

func (h MerchantHandler) CreateMerchant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := dto.MerchantCreateRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validation.Struct(req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	merchant := entities.Merchant{
		Name:             req.Name,
		MerchantCategory: entities.MerchantCategory(req.MerchantCategory),
		ImageURL:         req.ImageURL,
		Location: entities.Location{
			Lat:  req.Location.Lat,
			Long: req.Location.Long,
		},
	}

	mrc, err := h.service.CreateMerchant(ctx, merchant)
	if err != nil {
		if appErr, ok := err.(utils.AppError); ok {
			utils.SendErrorResponse(w, appErr.StatusCode, appErr.Message)
		} else {
			utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		}

		return
	}

	utils.SendResponse(w, http.StatusCreated, mrc)
}

func (h MerchantHandler) GetAllMerchant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()

	filter := entities.MerchantFilter{
		MerchantID:       q.Get("merchantId"),
		Name:             q.Get("name"),
		MerchantCategory: q.Get("merchantCategory"),
		SortCreatedAt:    strings.ToLower(q.Get("createdAt")),
	}

	// Parse limit
	if l := q.Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			filter.Limit = val
		} else {
			utils.SendErrorResponse(w, http.StatusBadRequest, "invalid limit parameter")
			return
		}
	}

	// Parse offset
	if o := q.Get("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil {
			filter.Offset = val
		} else {
			utils.SendErrorResponse(w, http.StatusBadRequest, "invalid offset parameter")
			return
		}
	}

	// Call service
	merchants, err := h.service.GetMerchants(ctx, filter)
	if err != nil {
		if appErr, ok := err.(utils.AppError); ok {
			utils.SendErrorResponse(w, appErr.StatusCode, appErr.Message)
		} else {
			utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Always return 200 with list (can be empty)
	utils.SendResponse(w, http.StatusOK, merchants)
}
