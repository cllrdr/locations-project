package handler

import (
	"fmt"
	"net/http"
	"time"

	"locations-project/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
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

	token, err := h.generateToken(createdUser.ID, createdUser.IsModerator)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"user":  createdUser,
		"token": token,
	})
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

	token, err := h.generateToken(authenticatedUser.ID, authenticatedUser.IsModerator)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":  authenticatedUser,
		"token": token,
	})
}

func (h *Handler) LogoutAPI(ctx *gin.Context) {
	tokenStr := ctx.GetHeader("Authorization")
	if tokenStr == "" {
		h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("missing authorization header"))
		return
	}

	conn, err := redis.Dial("tcp", h.Config.RedisAddr)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	defer conn.Close()

	_, err = conn.Do("SET", tokenStr, "blacklisted", "EX", int64(time.Hour*24))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "logged out",
	})
}
