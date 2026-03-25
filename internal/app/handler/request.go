package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AddLocationToCart добавляет локацию в корзину (заявку) и делает выбор рандома
func (h *Handler) AddLocationToCart(ctx *gin.Context) {
	locationIDStr := ctx.Param("id")
	locationID, err := strconv.Atoi(locationIDStr)
	if err != nil {
		logrus.Error("Error converting location ID:", err)
	}

	// Пытаемся получить существующий черновик
	draftRequest, _, err := h.Repository.GetDraftRequestInfo()
	if err != nil {
		// Черновика нет, создаём новую заявку с локацией
		_, err := h.Repository.CreateRequestWithLocation(uint(locationID))
		if err != nil {
			logrus.Error("Error creating new draft request:", err)
		}
	} else {
		// Черновик есть, добавляем локацию в него
		err = h.Repository.AddLocationToRequest(draftRequest.ID, uint(locationID))
		if err != nil {
			logrus.Error("Error adding location to existing request:", err)
		}
	}

	// Делаем выбор рандома среди всех локаций в корзине
	_, err = h.Repository.ChooseRandomLocationForUser(1)
	if err != nil {
		logrus.Error("Error choosing random location:", err)
	}

	ctx.Redirect(http.StatusFound, "/all-locations")
}

// DeleteRequest удаляет заявку (меняет статус на "удалён")
func (h *Handler) DeleteRequest(ctx *gin.Context) {
	idRequest := ctx.Param("id")
	id, err := strconv.Atoi(idRequest)
	if err != nil {
		logrus.Errorf("Error converting request ID: %v", err)
	}

	err = h.Repository.DeleteRequest(uint(id))
	if err != nil {
		logrus.Errorf("Error deleting request: %v", err)
	}

	ctx.Redirect(http.StatusFound, "/all-locations")
}
