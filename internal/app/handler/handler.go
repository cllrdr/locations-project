package handler

import (
	"locations-project/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
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

		requests := api.Group("/requests")
		{
			requests.GET("/cart", h.DraftRequestInfoAPI)
			requests.GET("", h.GetRequestsAPI)
			requests.GET("/:id", h.GetRequestAPI)
			requests.PUT("/:id", h.UpdateRequestAPI)
			requests.DELETE("/:id", h.DeleteRequestAPI)
			requests.PUT("/:id/form", h.FormRequestAPI)
			requests.PUT("/:id/complete", h.CompleteRequestAPI)
			requests.POST("/:id/locations/:locationId", h.AddLocationToRequestAPI)
			requests.PUT("/:id/locations/:locationId", h.UpdateLocationPriorityAPI)
			requests.DELETE("/:id/locations/:locationId", h.RemoveLocationFromRequestAPI)
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
