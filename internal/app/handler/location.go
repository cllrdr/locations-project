package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"locations-project/internal/app/ds"
)

func (h *Handler) GetLocationsAPI(ctx *gin.Context) {
	locationName := ctx.Query("location")

	locations, err := h.Repository.GetLocations(locationName)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, locations)
}

func (h *Handler) GetLocationAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	location, err := h.Repository.GetLocation(uint(id))
	if err != nil {
		if err.Error() == "location not found" {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, location)
}

func (h *Handler) CreateLocationAPI(ctx *gin.Context) {
	type locationInput struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
		Players     string `json:"players" binding:"required"`
	}

	var input locationInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	location := ds.Location{
		Name:        input.Name,
		Description: input.Description,
		Players:     input.Players,
	}

	createdLocation, err := h.Repository.CreateLocation(location)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, createdLocation)
}

func (h *Handler) UpdateLocationAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	type locationUpdateRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Players     string `json:"players"`
	}

	var req locationUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	location := ds.Location{
		Name:        req.Name,
		Description: req.Description,
		Players:     req.Players,
	}

	err = h.Repository.UpdateLocation(uint(id), location)
	if err != nil {
		if err.Error() == "location not found" {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	// TODO: Удалить этот блок для продакшена, оставить только ctx.Status(http.StatusOK)
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Локация обновлена",
	})
}

func (h *Handler) DeleteLocationAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	location, err := h.Repository.GetLocation(uint(id))
	if err != nil {
		if err.Error() == "location not found" {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	// Удаляем изображение из MinIO, если оно существует
	if location.ImagePath.Valid {
		fileName := location.ImagePath.String
		parts := strings.Split(fileName, "/")
		if len(parts) > 0 {
			fileName = parts[len(parts)-1]
		}
		if fileName != "" {
			err = h.Repository.DeleteFileFromMinIO(ctx.Request.Context(), fileName)
			if err != nil {
				log.Printf("Warning: failed to delete image from MinIO: %v", err)
			}
		}
	}

	err = h.Repository.DeleteLocation(uint(id))
	if err != nil {
		if err.Error() == "location not found" {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	// TODO: Удалить этот блок для продакшена, оставить только ctx.Status(http.StatusOK)
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Локация удалена",
	})
}

func (h *Handler) UploadLocationImageAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	_, err = h.Repository.GetLocation(uint(id))
	if err != nil {
		if err.Error() == "location not found" {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	src, err := file.Open()
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	defer src.Close()

	fileName := "location_" + strconv.FormatUint(uint64(id), 10) + ".png"

	// Загружаем файл в MinIO
	err = h.Repository.UploadFileToMinIO(
		context.Background(),
		fileName,
		src,
		file.Size,
		file.Header.Get("Content-Type"),
	)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Сохраняем путь к файлу в базе данных
	imagePath := "http://127.0.0.1:9000/locations/" + fileName
	err = h.Repository.UpdateLocationImage(uint(id), imagePath)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"message":    "Изображение загружено",
		"image_path": imagePath,
	})
}

func (h *Handler) AddLocationToRequestAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	draft, _, err := h.Repository.GetDraftRequestInfo()
	if err != nil {
		created, createErr := h.Repository.CreateRequestWithLocation(uint(id))
		if createErr != nil {
			if strings.Contains(createErr.Error(), "duplicate key value violates unique constraint") {
				h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("локация уже добавлена в заявку"))
				return
			}
			h.errorHandler(ctx, http.StatusInternalServerError, createErr)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status":   "success",
			"message":  "Локация добавлена в новую заявку",
			"draft_id": created.ID,
		})
		return
	}

	if draft.Status != ds.RequestStatusDraft {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("заявка должна быть в статусе черновика"))
		return
	}

	if err := h.Repository.AddLocationToRequest(draft.ID, uint(id)); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("локация уже добавлена в заявку"))
			return
		}
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"message":  "Локация добавлена в заявку",
		"draft_id": draft.ID,
	})
}
