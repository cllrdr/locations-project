package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"locations-project/internal/app/ds"
)

func (h *Handler) GetLocations(ctx *gin.Context) {
	var locations []ds.Location
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

	// Получаем информацию о черновике заявки
	draftRequest, chosenLocations, err := h.Repository.GetDraftRequestInfo()
	var draftRequestID uint = 0
	var locationsCount int64 = 0
	if err == nil {
		draftRequestID = draftRequest.ID
		locationsCount = int64(len(chosenLocations))
	}

	ctx.HTML(http.StatusOK, "all-locations.html", gin.H{
		"time":           time.Now().Format("15:04:05"),
		"locations":      locations,
		"query":          searchLocation,
		"draftRequestID": draftRequestID,
		"locationsCount": locationsCount,
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

func (h *Handler) GetPlayersLocations(ctx *gin.Context) {
	idRequest := ctx.Param("id")
	id, err := strconv.Atoi(idRequest)
	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusOK, "fav-locations.html", gin.H{
			"is404": true,
		})
		return
	}

	request, chosenLocations, err := h.Repository.GetPlayersLocationsForRequest(id)
	if err != nil {
		logrus.Error(err)
		ctx.HTML(http.StatusOK, "fav-locations.html", gin.H{
			"is404": true,
		})
		return
	}

	if request.Status != ds.RequestStatusDraft {
		ctx.HTML(http.StatusOK, "fav-locations.html", gin.H{
			"is404": true,
		})
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
		"is404":           false,
	})
}
