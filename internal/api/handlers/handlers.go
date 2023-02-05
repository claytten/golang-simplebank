package handlers

import (
	"github.com/claytten/golang-simplebank/internal/api"
	"github.com/claytten/golang-simplebank/internal/api/handlers/auth"
	"github.com/claytten/golang-simplebank/internal/api/middlewares"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	api *api.Server
	rg  *gin.RouterGroup
}

func NewHandler(api *api.Server, rg *gin.RouterGroup) *Handler {
	return &Handler{api, rg}
}

func (h *Handler) ApplyAllAuthRoutes() {
	user := h.rg.Group("auth")
	{
		auth.PostLoginUserRoute(h.api, user)
		user.Use(middlewares.AuthMiddleware(h.api.Token))
		// just middleware basic authentication
		auth.GetUserRoute(h.api, user)
		auth.PostCreateUserRoute(h.api, user)

		// adding middleware for checking username and old password
		user.Use(middlewares.CheckOwnUserUpdate(h.api.DB))
		auth.UpdateUserProfileRoute(h.api, user)
		auth.UpdateUserPasswordRoute(h.api, user)
	}
}
