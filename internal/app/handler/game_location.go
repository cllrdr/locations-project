package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"locations-project/internal/app/ds"

	"github.com/gin-gonic/gin"
)

func (h *Handler) UpdateLocationPriorityAPI(ctx *gin.Context) {
	requestIDStr := ctx.Param("id")
	locationIDStr := ctx.Param("locationId")
	requestID, err := strconv.ParseUint(requestIDStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	locationID, err := strconv.ParseUint(locationIDStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	request, err := h.Repository.GetRequest(uint(requestID))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if request.Status != ds.GameStatusDraft {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("заявка должна быть в статусе черновика"))
		return
	}

	var body struct {
		Priority int `json:"priority" binding:"required,gt=0"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("приоритет должен быть больше 0"))
		return
	}

	if err := h.Repository.UpdateLocationPriority(uint(requestID), uint(locationID), body.Priority); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// TODO: Удалить этот блок для продакшена
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Приоритет обновлен",
	})
}

func (h *Handler) RemoveLocationFromGameAPI(ctx *gin.Context) {
	requestIDStr := ctx.Param("id")
	locationIDStr := ctx.Param("locationId")
	requestID, err := strconv.ParseUint(requestIDStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	locationID, err := strconv.ParseUint(locationIDStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	request, err := h.Repository.GetRequest(uint(requestID))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if request.Status != ds.GameStatusDraft {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("заявка должна быть в статусе черновика"))
		return
	}

	err = h.Repository.RemoveLocationFromRequest(uint(requestID), uint(locationID))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// TODO: Удалить этот блок для продакшена
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Локация удалена из заявки",
	})
}
