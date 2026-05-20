package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"locations-project/internal/app/ds"

	"github.com/gin-gonic/gin"
)

// GetGamesAPI godoc
// @Summary      Список заявок
// @Description  Модераторы видят все заявки, обычные пользователи — только свои
// @Tags         Games
// @Accept       json
// @Produce      json
// @Param        status      query     string  false  "Фильтр по статусу"
// @Param        start_date  query     string  false  "Дата начала (YYYY-MM-DD)"
// @Param        end_date    query     string  false  "Дата конца (YYYY-MM-DD)"
// @Success      200         {array}   ds.GamesListResponse
// @Failure      401         {object}  ds.ErrorResponse
// @Failure      500         {object}  ds.ErrorResponse
// @Security     BearerAuth
// @Router       /api/games [get]
func (h *Handler) GetGamesAPI(ctx *gin.Context) {
	// 🔐 Получаем данные пользователя из контекста
	userID, exists := ctx.Get("user_id")
	if !exists {
		h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("user id not found"))
		return
	}
	isModerator, _ := ctx.Get("is_moderator")

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

	var requests []ds.PlayersLocationGame
	var err error

	// 🔐 Если НЕ модератор — фильтруем только по своим заявкам
	if isModerator != true {
		requests, err = h.Repository.GetRequestsByCreator(userID.(uint), status, startDate, endDate)
	} else {
		// Модератор видит всё
		requests, err = h.Repository.GetRequests(status, startDate, endDate)
	}

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	var simplifiedRequests []ds.GamesListResponse
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

		simplifiedRequests = append(simplifiedRequests, ds.GamesListResponse{
			ID:            req.ID,
			Nickname:      req.Nickname,
			Status:        string(req.Status),
			CreatedAt:     req.CreatedAt,
			FormedAt:      req.FormedAt,
			CompletedAt:   req.CompletedAt,
			CreatorName:   creatorName,
			ModeratorName: moderatorName,
			RandomPool:    selectedLocationsCount,
		})
	}

	ctx.JSON(http.StatusOK, simplifiedRequests)
}

// GetGameAPI godoc
// @Summary      Получить заявку по ID
// @Description  Возвращает детальную информацию о заявке с локациями. Доступ владельцу или модератору.
// @Tags         Games
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "ID заявки"
// @Success      200 {object}  ds.GameResponse
// @Failure      400 {object}  ds.ErrorResponse
// @Failure      403 {object}  ds.ErrorResponse
// @Failure      404 {object}  ds.ErrorResponse
// @Failure      500 {object}  ds.ErrorResponse
// @Security     BearerAuth
// @Router       /api/games/{id} [get]
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

	// 🔐 Проверка владения
	if !h.IsOwnerOrModerator(ctx, request.CreatorID) {
		h.errorHandler(ctx, http.StatusForbidden, fmt.Errorf("access denied"))
		return
	}

	var simplifiedLocations []ds.LocationResponse
	for _, loc := range locations {
		simplifiedLocations = append(simplifiedLocations, ds.LocationResponse{
			ID:            loc.ID,
			LocationID:    loc.LocationID,
			Priority:      loc.Priority,
			LocationName:  loc.Location.Name,
			LocationImage: loc.Location.ImagePath,
			IsRandomed:    loc.IsRandomed,
		})
	}

	ctx.JSON(http.StatusOK, ds.GameResponse{
		ID:          request.ID,
		Nickname:    request.Nickname,
		Status:      string(request.Status),
		CreatedAt:   request.CreatedAt,
		FormedAt:    request.FormedAt,
		CompletedAt: request.CompletedAt,
		Locations:   simplifiedLocations,
	})
}

// UpdateGameAPI godoc
// @Summary      Обновить заявку
// @Description  Обновляет nickname заявки. Доступ владельцу или модератору.
// @Tags         Games
// @Accept       json
// @Produce      json
// @Param        id    path      int              true  "ID заявки"
// @Param        body  body      ds.PlayersLocationGame  true  "Данные для обновления"
// @Success      200   {object}  ds.MessageResponse
// @Failure      400   {object}  ds.ErrorResponse
// @Failure      403   {object}  ds.ErrorResponse
// @Failure      404   {object}  ds.ErrorResponse
// @Failure      500   {object}  ds.ErrorResponse
// @Security     BearerAuth
// @Router       /api/games/{id} [put]
func (h *Handler) UpdateGameAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var req ds.PlayersLocationGame
	if err := ctx.ShouldBindJSON(&req); err != nil {
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
		h.errorHandler(ctx, http.StatusForbidden, fmt.Errorf("access denied"))
		return
	}

	currentRequest.Nickname = req.Nickname

	err = h.Repository.UpdateRequest(uint(id), currentRequest)
	if err != nil {
		if err.Error() == "record not found" {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, ds.MessageResponse{
		Message: "Заявка обновлена",
	})
}

// DeleteGameAPI godoc
// @Summary      Удалить заявку
// @Description  Мягкое удаление (смена статуса на 'удалён'). Доступ владельцу или модератору.
// @Tags         Games
// @Produce      json
// @Param        id  path      int  true  "ID заявки"
// @Success      200 {object}  ds.MessageResponse
// @Failure      400 {object}  ds.ErrorResponse
// @Failure      403 {object}  ds.ErrorResponse
// @Failure      404 {object}  ds.ErrorResponse
// @Failure      500 {object}  ds.ErrorResponse
// @Security     BearerAuth
// @Router       /api/games/{id} [delete]
func (h *Handler) DeleteGameAPI(ctx *gin.Context) {
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
		h.errorHandler(ctx, http.StatusForbidden, fmt.Errorf("access denied"))
		return
	}

	err = h.Repository.DeleteRequest(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, ds.MessageResponse{
		Message: "Заявка удалена",
	})
}

// DraftGameInfoAPI godoc
// @Summary      Информация о черновике
// @Description  Возвращает ID черновика и количество локаций в нём.
// @Tags         Games
// @Produce      json
// @Success      200 {object}  ds.DraftInfoResponse
// @Failure      401 {object}  ds.ErrorResponse
// @Failure      500 {object}  ds.ErrorResponse
// @Security     BearerAuth
// @Router       /api/games/cart [get]
func (h *Handler) DraftGameInfoAPI(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("user id not found"))
		return
	}

	draft, locations, err := h.Repository.GetDraftRequestInfo(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusOK, ds.DraftInfoResponse{
			DraftID:      0,
			LocationsCnt: 0,
		})
		return
	}

	ctx.JSON(http.StatusOK, ds.DraftInfoResponse{
		DraftID:      draft.ID,
		LocationsCnt: len(locations),
	})
}

// FormGameAPI godoc
// @Summary      Сформировать заявку
// @Description  Переводит черновик в статус 'сформирован'. Доступ владельцу или модератору.
// @Tags         Games
// @Produce      json
// @Param        id  path      int  true  "ID заявки"
// @Success      200 {object}  ds.MessageResponse
// @Failure      400 {object}  ds.ErrorResponse
// @Failure      403 {object}  ds.ErrorResponse
// @Failure      404 {object}  ds.ErrorResponse
// @Failure      500 {object}  ds.ErrorResponse
// @Security     BearerAuth
// @Router       /api/games/{id}/form [put]
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
		h.errorHandler(ctx, http.StatusForbidden, fmt.Errorf("access denied"))
		return
	}

	if err := h.Repository.FormRequest(uint(id)); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, ds.MessageResponse{
		Message: "Заявка сформирована",
	})
}

// CompleteGameAPI godoc
// @Summary      Завершить или отклонить заявку
// @Description  Доступно только модераторам. approve=true завершает, false отклоняет.
// @Tags         Games
// @Accept       json
// @Produce      json
// @Param        id    path      int  true  "ID заявки"
// @Param        body  body      object  true  "{\"approve\": true/false}"
// @Success      200   {object}  ds.MessageResponse
// @Failure      400   {object}  ds.ErrorResponse
// @Failure      403   {object}  ds.ErrorResponse
// @Failure      500   {object}  ds.ErrorResponse
// @Security     BearerAuth
// @Router       /api/games/{id}/complete [put]
func (h *Handler) CompleteGameAPI(ctx *gin.Context) {
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

	ctx.JSON(http.StatusOK, ds.MessageResponse{
		Message: message,
	})
}

// AddLocationToGameAPI godoc
// @Summary      Добавить локацию в заявку
// @Description  Если черновика нет — создаёт новую заявку. Доступ владельцу или модератору.
// @Tags         GamLoc
// @Produce      json
// @Param        locationId  path      int  true  "ID локации"
// @Success      201         {object}  ds.MessageResponse
// @Failure      400         {object}  ds.ErrorResponse
// @Failure      401         {object}  ds.ErrorResponse
// @Failure      403         {object}  ds.ErrorResponse
// @Failure      500         {object}  ds.ErrorResponse
// @Security     BearerAuth
// @Router       /api/gamloc/locations/{locationId} [post]
func (h *Handler) AddLocationToGameAPI(ctx *gin.Context) {
	locationIDStr := ctx.Param("locationId")
	locationID, err := strconv.ParseUint(locationIDStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("user id not found"))
		return
	}

	currentUserID := userID.(uint)

	draft, _, err := h.Repository.GetDraftRequestInfo(currentUserID)
	var request ds.PlayersLocationGame

	if err != nil {
		// Черновика нет — создаём новую заявку
		request, err = h.Repository.CreateRequestWithLocation(uint(locationID), currentUserID)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusCreated, ds.MessageResponse{
			Message: fmt.Sprintf("Заявка создана и локация добавлена. ID: %d", request.ID),
		})
		return
	}

	// Черновик есть — проверяем владение
	if !h.IsOwnerOrModerator(ctx, draft.CreatorID) {
		h.errorHandler(ctx, http.StatusForbidden, fmt.Errorf("access denied"))
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

	ctx.JSON(http.StatusCreated, ds.MessageResponse{
		Message: fmt.Sprintf("Локация добавлена в заявку. ID заявки: %d", requestID),
	})
}