package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"locations-project/internal/app/ds"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetRequestsAPI(ctx *gin.Context) {
	var status *ds.RequestStatus
	var startDate, endDate *time.Time

	if statusStr := ctx.Query("status"); statusStr != "" {
		requestStatus := ds.RequestStatus(statusStr)
		status = &requestStatus
	}

	if startDateStr := ctx.Query("start_date"); startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &parsed
		}
	}

	if endDateStr := ctx.Query("end_date"); endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &parsed
		}
	}

	requests, err := h.Repository.GetRequests(status, startDate, endDate)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	var simplifiedRequests []gin.H
	for _, req := range requests {
		creatorName := ""
		if req.Creator.Name != "" {
			creatorName = req.Creator.Name
		}

		moderatorName := ""
		if req.Moderator.Name != "" {
			moderatorName = req.Moderator.Name
		}

		simplifiedRequests = append(simplifiedRequests, gin.H{
			"id":             req.ID,
			"nickname":       req.Nickname,
			"status":         req.Status,
			"created_at":     req.CreatedAt,
			"formed_at":      req.FormedAt,
			"completed_at":   req.CompletedAt,
			"creator_name":   creatorName,
			"moderator_name": moderatorName,
		})
	}

	ctx.JSON(http.StatusOK, simplifiedRequests)
}

func (h *Handler) GetRequestAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	request, locations, err := h.Repository.GetRequestWithLocations(uint(id))
	if err != nil {
		if err.Error() == "record not found" {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	var simplifiedLocations []gin.H
	for _, loc := range locations {
		simplifiedLocations = append(simplifiedLocations, gin.H{
			"id":             loc.ID,
			"location_id":    loc.LocationID,
			"priority":       loc.Priority,
			"location_name":  loc.Location.Name,
			"location_image": loc.Location.ImagePath,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":           request.ID,
		"nickname":     request.Nickname,
		"status":       request.Status,
		"created_at":   request.CreatedAt,
		"formed_at":    request.FormedAt,
		"completed_at": request.CompletedAt,
		"locations":    simplifiedLocations,
	})
}

func (h *Handler) UpdateRequestAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var request ds.PlayersLocationRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Получаем текущую заявку
	currentRequest, err := h.Repository.GetRequest(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	// Обновляем только разрешённые поля (системные поля не меняются)
	currentRequest.Nickname = request.Nickname

	err = h.Repository.UpdateRequest(uint(id), currentRequest)
	if err != nil {
		if err.Error() == "record not found" {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	// TODO: Удалить этот блок для продакшена
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Заявка обновлена",
	})
}

func (h *Handler) DeleteRequestAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.DeleteRequest(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// TODO: Удалить этот блок для продакшена
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Заявка удалена",
	})
}

func (h *Handler) DraftRequestInfoAPI(ctx *gin.Context) {
	draft, locations, err := h.Repository.GetDraftRequestInfo()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"draft_id":      0,
			"locations_cnt": 0,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"draft_id":      draft.ID,
		"locations_cnt": len(locations),
	})
}

func (h *Handler) FormRequestAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if err := h.Repository.FormRequest(uint(id)); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	// TODO: Удалить этот блок для продакшена
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Заявка сформирована",
	})
}

func (h *Handler) CompleteRequestAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var req struct {
		Approve bool `json:"approve"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.CompleteRequest(uint(id), req.Approve); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	message := "Заявка отклонена"
	if req.Approve {
		message = "Заявка одобрена и обработана"
	}

	// TODO: Удалить этот блок для продакшена
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": message,
	})
}

func (h *Handler) AddLocationToRequestAPI(ctx *gin.Context) {
	requestIDStr := ctx.Param("id")
	requestID, err := strconv.ParseUint(requestIDStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	locationIDStr := ctx.Param("locationId")
	locationID, err := strconv.ParseUint(locationIDStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Получаем заявку
	request, err := h.Repository.GetRequest(uint(requestID))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	// Проверяем, что заявка в статусе черновика
	if request.Status != ds.RequestStatusDraft {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("заявка должна быть в статусе черновика"))
		return
	}

	// Добавляем локацию в заявку
	if err := h.Repository.AddLocationToRequest(uint(requestID), uint(locationID)); err != nil {
		if err.Error() == "локация уже добавлена в заявку" {
			h.errorHandler(ctx, http.StatusBadRequest, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Локация добавлена в заявку",
	})
}
