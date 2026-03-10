package handler

import (
  "github.com/gin-gonic/gin"
  "github.com/sirupsen/logrus"
  "locations-project/internal/app/repository"
  "strconv"
  "net/http"
  "time"
)

type Handler struct {
  Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
  return &Handler{
    Repository: r,
  }
}

func (h *Handler) GetLocations(ctx *gin.Context) {
	var locations []repository.Location
	var err error

	searchQuery := ctx.Query("query") // получаем значение из поля поиска
	if searchQuery == "" {            // если поле поиска пусто, то просто получаем из репозитория все записи
		locations, err = h.Repository.GetLocations()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		locations, err = h.Repository.GetLocationsByName(searchQuery) // в ином случае ищем локацию по имени
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.HTML(http.StatusOK, "all-locations.html", gin.H{
		"time":      time.Now().Format("15:04:05"),
		"locations": locations,
		"query":     searchQuery, // передаем введенный запрос обратно на страницу
		// в ином случае оно будет очищаться при нажатии на кнопку
	})
}

func (h *Handler) GetLocation(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id локации из урла (то есть из /location/:id)
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}

	location, err := h.Repository.GetLocation(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "location.html", gin.H{
		"location": location,
	})
}

func (h *Handler) GetFavorites(ctx *gin.Context) {
	// получаем локацию с ID=3
	location, err := h.Repository.GetLocation(3)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "fav-locations.html", gin.H{
		"location": location,
	})
}
