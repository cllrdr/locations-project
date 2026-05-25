package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"locations-project/internal/app/ds"
)

// RegisterUserAPI godoc
// @Summary      Регистрация пользователя
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        user  body      ds.RegisterRequest  true  "Данные для регистрации"
// @Success      201   {object}  ds.LoginResponse    "Успешная регистрация"
// @Failure      400   {object}  ds.ErrorResponse    "Неверный формат запроса или пустые поля"
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
// @Success      200   {object}  ds.LoginResponse "Успешная авторизация"
// @Failure      401   {object}  ds.ErrorResponse "Неверный email или пароль"
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
// @Param        Authorization  header  string  true   "Bearer <token>"
// @Success      200  {object}  ds.MessageResponse   "Успешный выход"
// @Failure      401  {object}  ds.ErrorResponse     "Нет токена или он невалидный"
// @Failure      500  {object}  ds.ErrorResponse     "Ошибка сервера"
// @Security     BearerAuth
// @Router       /api/profile/logout [post]
func (h *Handler) LogoutAPI(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	
	// ✅ ВАЖНО: Удаляем "Bearer " так же, как в middleware!
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	
	// ✅ Проверяем только на пустоту (не требуем обязательного "Bearer ")
	if strings.TrimSpace(tokenStr) == "" {
		h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("authorization header required"))
		return
	}

	// Берём соединение из пула (не создаём новое!)
	conn := h.RedisPool.Get()
	defer conn.Close()

	// Парсим токен только для получения времени истечения (exp)
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.Config.JWTSecret), nil
	})

	// Если токен уже невалиден — просто отвечаем 200 (клиент хочет выйти)
	if err != nil || !token.Valid {
		ctx.JSON(http.StatusOK, ds.MessageResponse{Message: "logged out"})
		return
	}

	// ✅ Динамический TTL: сколько секунд осталось до истечения токена
	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl > 0 {
		// SET <token> "blacklisted" EX <seconds>
		_, err := conn.Do("SET", tokenStr, "blacklisted", "EX", int64(ttl.Seconds()))
		if err != nil {
			// Логируем, но не роняем логаут
			// logrus.WithError(err).Warn("failed to blacklist token")
		}
	}

	ctx.JSON(http.StatusOK, ds.MessageResponse{Message: "logged out"})
}