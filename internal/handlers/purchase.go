package handlers

import (
	"belimang/internal/dto"
	"belimang/internal/entities"
	"belimang/internal/services"
	"belimang/internal/middleware"
	"belimang/internal/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type PurchaseHandler struct {
	service services.PurchaseService
	validation *validator.Validate
}

func NewPurchaseHandler(service services.PurchaseService, validation *validator.Validate) PurchaseHandler {
	return PurchaseHandler{
		service: service,
		validation: validation,
	}
}

func (h PurchaseHandler) GetNearbyMerchants(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()

	authCtx, ok := middleware.GetAuthContext(ctx)
	if !ok {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	latStr := chi.URLParam(r, "lat")
	lonStr := chi.URLParam(r, "lon")

	lat, errLat := strconv.ParseFloat(latStr, 64)
	lon, errLon := strconv.ParseFloat(lonStr, 64)
	if errLat != nil || errLon != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "lat or lon is not valid")
		return
	}

	limit := 5
	if limStr := q.Get("limit"); limStr != "" {
		if limVal, err := strconv.Atoi(limStr); err == nil && limVal > 0 {
			limit = limVal
		}
	}

	offset := 0
	if offStr := q.Get("offset"); offStr != "" {
		if offVal, err := strconv.Atoi(offStr); err == nil && offVal >= 0 {
			offset = offVal
		}
	}

	filter := entities.MerchantNearbyFilter{
		Limit:            limit,
		Name:             q.Get("name"),
		MerchantID:       q.Get("merchantId"),
		MerchantCategory: q.Get("merchantCategory"),
		Offset:           offset,
		UserID:           authCtx.ID,
		Lat:              lat,
		Lon:              lon,
	}

	resp, err := h.service.GetNearbyMerchants(ctx, filter)
	if err != nil {
		if appErr, ok := err.(utils.AppError); ok {
			utils.SendErrorResponse(w, appErr.StatusCode, appErr.Message)
		} else {
			utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.SendResponse(w, http.StatusOK, resp)
}

func (h PurchaseHandler) CreateEstimate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := dto.EstimateReq{}

	authCtx, ok := middleware.GetAuthContext(ctx)
	if !ok {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	req.UserID = authCtx.ID

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validation.Struct(req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(req.UserPurchase) == 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, "orders must not be empty")
		return
	}

	response, err := h.service.CreateEstimate(ctx, req)
	if err != nil {
		if appErr, ok := err.(utils.AppError); ok {
			utils.SendErrorResponse(w, appErr.StatusCode, appErr.Message)
		} else {
			utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.SendResponse(w, http.StatusOK, response)
}

func (h PurchaseHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := dto.CreateOrderRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validation.Struct(req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.service.CreateOrder(ctx, req)
	if err != nil {
		if appErr, ok := err.(utils.AppError); ok {
			utils.SendErrorResponse(w, appErr.StatusCode, appErr.Message)
		} else {
			utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.SendResponse(w, http.StatusCreated, response)
}

func (h PurchaseHandler) GetAllOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()

	authCtx, ok := middleware.GetAuthContext(ctx)
	if !ok {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	limit := 5
	if limStr := q.Get("limit"); limStr != "" {
		if limVal, err := strconv.Atoi(limStr); err == nil && limVal > 0 {
			 limit = limVal
		}
	}

	offset := 0
	if offStr := q.Get("offset"); offStr != "" {
		if offVal, err := strconv.Atoi(offStr); err == nil && offVal >= 0 {
			 offset = offVal
		}
	}

	filter := entities.OrderFilter{
		Limit:            limit,
		Name:             q.Get("name"),
		MerchantID:       q.Get("merchantId"),
		MerchantCategory: q.Get("merchantCategory"),
		UserID:           authCtx.ID,
		Offset:           offset,
	}

	response, err := h.service.GetAllOrder(ctx, filter)
	if err != nil {
		if appErr, ok := err.(utils.AppError); ok {
			utils.SendErrorResponse(w, appErr.StatusCode, appErr.Message)
		} else {
			utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.SendResponse(w, http.StatusOK, response)
}
