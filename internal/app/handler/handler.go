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

	searchLocation := ctx.Query("location-search") 
	if searchLocation == "" {
		locations, err = h.Repository.GetLocations()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		locations, err = h.Repository.GetLocationsByName(searchLocation)
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.HTML(http.StatusOK, "all-locations.html", gin.H{
		"time": time.Now().Format("15:04:05"),
		"locations": locations,
		"query": searchLocation,
	})
}

func (h *Handler) GetLocation(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
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

func (h *Handler) GetPlayersLocations(ctx *gin.Context) {
	idRequest := ctx.Param("id")
	id, err := strconv.Atoi(idRequest)
	if err != nil {
		logrus.Error(err)
		return
	}

	request, chosenLocations, err := h.Repository.GetPlayersLocationsForRequest(id)
	if err != nil {
		logrus.Error(err)
		return
	}

	locations, err := h.Repository.GetLocations()
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "fav-locations.html", gin.H{
		"playerRequest":   request,
		"chosenLocations": chosenLocations,
		"locations":       locations,
	})
}


