package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"locations-project/internal/app/ds"

	"github.com/gin-gonic/gin"
)

// UpdateLocationPriorityAPI godoc
// @Summary      Обновить приоритет локации в заявке
// @Description  Изменяет приоритет локации в черновике заявки. Доступ владельцу или модератору.
// @Tags         GamLoc
// @Accept       json
// @Produce      json
// @Param        id          path      int     true  "ID заявки"
// @Param        locationId  path      int     true  "ID локации"
// @Param        body        body      object  true  "{\"priority\": 1}"
// @Success      200         {object}  ds.MessageResponse
// @Failure      400         {object}  ds.ErrorResponse
// @Failure      403         {object}  ds.ErrorResponse
// @Failure      404         {object}  ds.ErrorResponse
// @Failure      500         {object}  ds.ErrorResponse
// @Security     BearerAuth
// @Router       /api/gamloc/{id}/locations/{locationId} [put]
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

	// 🔐 Проверка владения
	if !h.IsOwnerOrModerator(ctx, request.CreatorID) {
		h.errorHandler(ctx, http.StatusForbidden, fmt.Errorf("access denied"))
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

	ctx.JSON(http.StatusOK, ds.MessageResponse{
		Message: "Приоритет обновлен",
	})
}

// RemoveLocationFromGameAPI godoc
// @Summary      Удалить локацию из заявки
// @Description  Удаляет локацию из черновика заявки. Доступ владельцу или модератору.
// @Tags         GamLoc
// @Produce      json
// @Param        id          path      int  true  "ID заявки"
// @Param        locationId  path      int  true  "ID локации"
// @Success      200         {object}  ds.MessageResponse
// @Failure      400         {object}  ds.ErrorResponse
// @Failure      403         {object}  ds.ErrorResponse
// @Failure      404         {object}  ds.ErrorResponse
// @Failure      500         {object}  ds.ErrorResponse
// @Security     BearerAuth
// @Router       /api/gamloc/{id}/locations/{locationId} [delete]
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

	// 🔐 Проверка владения
	if !h.IsOwnerOrModerator(ctx, request.CreatorID) {
		h.errorHandler(ctx, http.StatusForbidden, fmt.Errorf("access denied"))
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

	ctx.JSON(http.StatusOK, ds.MessageResponse{
		Message: "Локация удалена из заявки",
	})
}