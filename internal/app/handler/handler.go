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
			profile.GET("/me", h.GetMeAPI)
			profile.PUT("/me", h.UpdateMeAPI)
		}

		locations := api.Group("/locations")
		{
			locations.GET("", h.GetLocationsAPI)
			locations.GET("/:id", h.GetLocationAPI)
			locations.POST("", h.CreateLocationAPI)
			locations.PUT("/:id", h.UpdateLocationAPI)
			locations.DELETE("/:id", h.DeleteLocationAPI)
			locations.POST("/:id/image", h.UploadLocationImageAPI)
			locations.POST("/:id/addLocation", h.AddLocationToRequestAPI)
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
		}

		requestLocations := api.Group("/requestlocations")
		{
			requestLocations.PUT("/:id/location/:locationId", h.UpdateLocationPriorityAPI)
			requestLocations.DELETE("/:id/location/:locationId", h.RemoveLocationFromRequestAPI)
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
