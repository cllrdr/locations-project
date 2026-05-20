package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"locations-project/internal/app/ds"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetGamesAPI(ctx *gin.Context) {
	var status *ds.GameStatus
	var startDate, endDate *time.Time

	if statusStr := ctx.Query("status"); statusStr != "" {
		requestStatus := ds.GameStatus(statusStr)
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

		selectedLocationsCount := 0
		if req.Status == ds.GameStatusCompleted {
			count, err := h.Repository.GetRandomedLocationsCount(req.ID)
			if err == nil {
				selectedLocationsCount = count
			}
		}

		simplifiedRequests = append(simplifiedRequests, gin.H{
			"id":              req.ID,
			"nickname":        req.Nickname,
			"status":          req.Status,
			"created_at":      req.CreatedAt,
			"formed_at":       req.FormedAt,
			"completed_at":    req.CompletedAt,
			"creator_name":    creatorName,
			"moderator_name":  moderatorName,
			"random_pool":     selectedLocationsCount,
		})
	}

	ctx.JSON(http.StatusOK, simplifiedRequests)
}

func (h *Handler) GetGameAPI(ctx *gin.Context) {
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

	// 🔐 Проверка владения: модератор или создатель?
	if !h.IsOwnerOrModerator(ctx, request.CreatorID) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
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
			"is_randomed":    loc.IsRandomed,
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

func (h *Handler) UpdateGameAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var request ds.PlayersLocationGame
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	currentRequest, err := h.Repository.GetRequest(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	// 🔐 Проверка владения
	if !h.IsOwnerOrModerator(ctx, currentRequest.CreatorID) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

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

	ctx.JSON(http.StatusOK, gin.H{"message": "Заявка обновлена"})
}

func (h *Handler) DeleteGameAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// 🔐 Загружаем заявку ПЕРЕД удалением для проверки владения
	request, err := h.Repository.GetRequest(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	// 🔐 Проверка владения
	if !h.IsOwnerOrModerator(ctx, request.CreatorID) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	err = h.Repository.DeleteRequest(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Заявка удалена"})
}

func (h *Handler) DraftGameInfoAPI(ctx *gin.Context) {
	// ✅ Берем ID текущего пользователя из контекста
	userID, exists := ctx.Get("user_id")
	if !exists {
		h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("user id not found"))
		return
	}

	// ✅ Передаем ID в репозиторий
	draft, locations, err := h.Repository.GetDraftRequestInfo(userID.(uint))
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

func (h *Handler) FormGameAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	request, err := h.Repository.GetRequest(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	// 🔐 Проверка владения
	if !h.IsOwnerOrModerator(ctx, request.CreatorID) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	if err := h.Repository.FormRequest(uint(id)); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Заявка сформирована"})
}

func (h *Handler) CompleteGameAPI(ctx *gin.Context) {
	// Этот эндпоинт защищён RequireModerator() в handler.go
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

	ctx.JSON(http.StatusOK, gin.H{"message": message})
}

func (h *Handler) AddLocationToGameAPI(ctx *gin.Context) {
	locationIDStr := ctx.Param("locationId")
	locationID, err := strconv.ParseUint(locationIDStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// ✅ Берем ID текущего пользователя из контекста
	userID, exists := ctx.Get("user_id")
	if !exists {
		h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("user id not found"))
		return
	}
	
	currentUserID := userID.(uint)

	// ✅ Передаем ID в репозиторий
	draft, _, err := h.Repository.GetDraftRequestInfo(currentUserID)
	var request ds.PlayersLocationGame

	if err != nil {
		// Черновика нет — создаём новую заявку (владелец = текущий пользователь)
		// ✅ Передаем ID создателя в репозиторий
		request, err = h.Repository.CreateRequestWithLocation(uint(locationID), currentUserID)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"status":     "success",
			"message":    "Заявка создана и локация добавлена",
			"request_id": request.ID,
		})
		return
	}

	// Черновик есть — проверяем владение перед добавлением
	// 🔐 Проверка: модератор или создатель черновика?
	if !h.IsOwnerOrModerator(ctx, draft.CreatorID) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	requestID := draft.ID

	if draft.Status != ds.GameStatusDraft {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("заявка должна быть в статусе черновика"))
		return
	}

	if err := h.Repository.AddLocationToRequest(requestID, uint(locationID)); err != nil {
		if err.Error() == "локация уже добавлена в заявку" {
			h.errorHandler(ctx, http.StatusBadRequest, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":    "Локация добавлена в заявку",
		"request_id": requestID,
	})
}