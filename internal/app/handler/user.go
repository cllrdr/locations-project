package handler

import (
	"fmt"
	"net/http"
	"time"
	"locations-project/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

// RegisterUserAPI godoc
// @Summary      Регистрация пользователя
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        user  body      ds.RegisterRequest  true  "Данные для регистрации"
// @Success      201   {object}  ds.LoginResponse	"Успешная регистрация"
// @Failure      400   {object}  ds.ErrorResponse	"Неверный формат запроса или пустые поля"
// @Router       /api/profile/register [post]
func (h *Handler) RegisterUserAPI(ctx *gin.Context) {
	var req ds.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	newUser := ds.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	createdUser, err := h.Repository.CreateUser(newUser)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	token, err := h.generateToken(createdUser.ID, createdUser.IsModerator)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	response := ds.LoginResponse{
		Token: token,
		User: ds.UserResponse{
			ID:          createdUser.ID,
			Email:       createdUser.Email,
			Name:        createdUser.Name,
			IsModerator: createdUser.IsModerator,
		},
	}

	ctx.JSON(http.StatusCreated, response)
}

// LoginAPI godoc
// @Summary      Вход в систему
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        user  body      ds.LoginRequest  true  "Email и пароль"
// @Success      200   {object}  ds.LoginResponse	"Успешная авторизация"
// @Failure      401   {object}  ds.ErrorResponse	"Неверный email или пароль"
// @Router       /api/profile/login [post]
func (h *Handler) LoginAPI(ctx *gin.Context) {
	var req ds.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	authenticatedUser, err := h.Repository.CheckCredentials(req.Email, req.Password)
	if err != nil {
		h.errorHandler(ctx, http.StatusUnauthorized, err)
		return
	}

	token, err := h.generateToken(authenticatedUser.ID, authenticatedUser.IsModerator)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	response := ds.LoginResponse{
		Token: token,
		User: ds.UserResponse{
			ID:          authenticatedUser.ID,
			Email:       authenticatedUser.Email,
			Name:        authenticatedUser.Name,
			IsModerator: authenticatedUser.IsModerator,
		},
	}

	ctx.JSON(http.StatusOK, response)
}

// LogoutAPI godoc
// @Summary      Выход из аккаунта
// @Tags         Profile
// @Produce      json
// @Param        Authorization  header  string  true  "Bearer <token>"
// @Success      200  {object}  ds.MessageResponse	"Успешный выход"
// @Failure      401  {object}  ds.ErrorResponse	"Нет токена или он невалидный"
// @Failure      500  {object}  ds.ErrorResponse	"Ошибка сервера"
// @Security     BearerAuth
// @Router       /api/profile/logout [post]
func (h *Handler) LogoutAPI(ctx *gin.Context) {
	tokenStr := ctx.GetHeader("Authorization")
	if tokenStr == "" {
		h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("authorization header required"))
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

	ctx.JSON(http.StatusOK, ds.MessageResponse{
		Message: "logged out",
	})
}