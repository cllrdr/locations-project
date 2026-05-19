package handler

import (
	"locations-project/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
	Config     *config.Config
}

func NewHandler(r *repository.Repository, cfg *config.Config) *Handler {
	return &Handler{
		Repository: r,
		Config:     cfg,
	}
}

func (h *Handler) RegisterAPI(router *gin.Engine) {
	api := router.Group("/api")
	{
		profile := api.Group("/profile")
		{
			profile.POST("/register", h.RegisterUserAPI)
			profile.POST("/login", h.LoginAPI)
			profile.POST("/logout", h.LogoutAPI)
		}

		locations := api.Group("/locations")
		{
			locations.GET("", h.GetLocationsAPI)
			locations.GET("/:id", h.GetLocationAPI)
			locations.POST("", h.CreateLocationAPI)
		}

		games := api.Group("/games")
		{
			games.GET("/cart", h.DraftGameInfoAPI)
			games.GET("", h.GetGamesAPI)
			games.GET("/:id", h.GetGameAPI)
			games.PUT("/:id", h.UpdateGameAPI)
			games.DELETE("/:id", h.DeleteGameAPI)
			games.PUT("/:id/form", h.FormGameAPI)
			games.PUT("/:id/complete", h.CompleteGameAPI)
		}

		gamloc := api.Group("/gamloc")
		{
			gamloc.POST("/locations/:locationId", h.AddLocationToGameAPI)
			gamloc.PUT("/:id/locations/:locationId", h.UpdateLocationPriorityAPI)
			gamloc.DELETE("/:id/locations/:locationId", h.RemoveLocationFromGameAPI)
		}
	}
}

// RegisterStatic То же самое, что и с маршрутами, регистрируем статику
func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources")
}

// errorHandler для более удобного вывода ошибок
func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
