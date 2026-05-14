package handler

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"locations-project/internal/app/ds"

	"github.com/gin-gonic/gin"
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
	// Чтение полей из form-data
	name := ctx.PostForm("name")
	description := ctx.PostForm("description")
	players := ctx.PostForm("players")

	if name == "" || description == "" || players == "" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("обязательные поля: name, description, players"))
		return
	}

	// Загружаем изображение
	imageFile, err := ctx.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("ошибка при загрузке изображения: %v", err))
		return
	}

	// Загружаем видео
	videoFile, err := ctx.FormFile("video")
	if err != nil && err != http.ErrMissingFile {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("ошибка при загрузке видео: %v", err))
		return
	}

	var imagePath, videoPath *string

	// Загружаем изображение в MinIO
	if imageFile != nil {
		imageName := generateFileName("location_image")
		imageSrc, err := imageFile.Open()
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
		defer imageSrc.Close()

		err = h.Repository.UploadFileToMinIO(
			ctx.Request.Context(),
			imageName,
			imageSrc,
			imageFile.Size,
			imageFile.Header.Get("Content-Type"),
		)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}

		imageURL := "http://127.0.0.1:9000/locations/" + imageName
		imagePath = &imageURL
	}

	// Загружаем видео в MinIO
	if videoFile != nil {
		videoName := generateFileName("location_video")
		videoSrc, err := videoFile.Open()
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
		defer videoSrc.Close()

		err = h.Repository.UploadFileToMinIO(
			ctx.Request.Context(),
			videoName,
			videoSrc,
			videoFile.Size,
			videoFile.Header.Get("Content-Type"),
		)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}

		videoURL := "http://127.0.0.1:9000/locations/" + videoName
		videoPath = &videoURL
	}

	// Создаём локацию
	location := ds.Location{
		Name:             name,
		Description:      description,
		ShortDescription: description, // используем description как short_description
		Players:          players,
	}

	if imagePath != nil {
		location.ImagePath = *imagePath
	}
	if videoPath != nil {
		location.VideoPath = *videoPath
	}

	createdLocation, err := h.Repository.CreateLocation(location)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, createdLocation)
}

// generateFileName генерирует имя файла на латинице
func generateFileName(prefix string) string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 16)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return prefix + "_" + string(b)
}
