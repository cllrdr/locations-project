package handler

import (
	"net/http"

	"locations-project/internal/app/ds"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterUserAPI(ctx *gin.Context) {
	var user ds.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	createdUser, err := h.Repository.CreateUser(user)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusCreated, createdUser)
}

func (h *Handler) LoginAPI(ctx *gin.Context) {
	var user ds.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	authenticatedUser, err := h.Repository.CheckCredentials(user.Email, user.Password)
	if err != nil {
		h.errorHandler(ctx, http.StatusUnauthorized, err)
		return
	}
	ctx.JSON(http.StatusOK, authenticatedUser)
}

func (h *Handler) LogoutAPI(ctx *gin.Context) {
	// TODO: Удалить этот блок для продакшена
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "logged out",
	})
}
