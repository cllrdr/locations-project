package handler

import (
	"locations-project/internal/app/config"
	"locations-project/internal/app/repository"
    "locations-project/internal/app/ds"

	"github.com/gin-gonic/gin"
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
		// === PROFILE ===
		profile := api.Group("/profile")
		{
			profile.POST("/register", h.RegisterUserAPI)                    // public
			profile.POST("/login", h.LoginAPI)                              // public
			profile.POST("/logout", h.AuthMiddleware(), h.LogoutAPI)        // auth only
		}

		// === LOCATIONS ===
		locations := api.Group("/locations")
		{
			locations.GET("", h.GetLocationsAPI)                            // public
			locations.GET("/:id", h.GetLocationAPI)                         // public
			locations.POST("", h.AuthMiddleware(), RequireModerator(), h.CreateLocationAPI) // moderator only
		}

		// === GAMES ===
		games := api.Group("/games")
		{
			games.GET("/cart", h.AuthMiddleware(), h.DraftGameInfoAPI)      // auth only
			games.GET("", h.AuthMiddleware(), h.GetGamesAPI) // auth + ownership/moderator
			games.GET("/:id", h.AuthMiddleware(), h.GetGameAPI)             // auth + ownership/moderator
			games.PUT("/:id", h.AuthMiddleware(), h.UpdateGameAPI)          // auth + ownership/moderator
			games.DELETE("/:id", h.AuthMiddleware(), h.DeleteGameAPI)       // auth + ownership/moderator
			games.PUT("/:id/form", h.AuthMiddleware(), h.FormGameAPI)       // auth + ownership/moderator
			games.PUT("/:id/complete", h.AuthMiddleware(), RequireModerator(), h.CompleteGameAPI) // moderator only
		}

		// === GAMLOC ===
		gamloc := api.Group("/gamloc")
		{
			gamloc.POST("/locations/:locationId", h.AuthMiddleware(), h.AddLocationToGameAPI) // auth + ownership/moderator
			gamloc.PUT("/:id/locations/:locationId", h.AuthMiddleware(), h.UpdateLocationPriorityAPI) // auth + ownership/moderator
			gamloc.DELETE("/:id/locations/:locationId", h.AuthMiddleware(), h.RemoveLocationFromGameAPI) // auth + ownership/moderator
		}
	}
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	ctx.JSON(errorStatusCode, ds.ErrorResponse{
		Error: err.Error(),
	})
}