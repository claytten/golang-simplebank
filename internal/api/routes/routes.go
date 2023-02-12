package routes

import (
	"github.com/claytten/golang-simplebank/internal/api"
	"github.com/claytten/golang-simplebank/internal/api/forms"
	"github.com/claytten/golang-simplebank/internal/api/handlers"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func ApplyAllPublicRoutes(s *api.Server) {
	s.Engine = gin.Default()

	rg := s.Engine.Group("/api")

	rg1 := rg.Group("/v1")

	//Custom form validator
	binding.Validator = new(forms.DefaultValidator)

	handlers := handlers.NewHandler(s, rg1)

	handlers.ApplyAllAuthRoutes()
	handlers.ApplyAllAccountRoutes()
}
