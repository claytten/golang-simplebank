package auth

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/claytten/golang-simplebank/internal/api"
	"github.com/claytten/golang-simplebank/internal/api/middlewares"
	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
	authorizationUsername   = "username"
)

func UpdateUserProfileRoute(api *api.Server, userRg *gin.RouterGroup) {
	userRg.Use(middlewares.CheckOwnUserUpdate(api.DB))
	userRg.PUT("/profile", UpdateUserProfileHandler(api))
}

type UpdateUserProfileRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func UpdateUserProfileHandler(s *api.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req UpdateUserProfileRequest
		username := ctx.GetHeader(authorizationUsername)

		// binding update user profile
		if err := ctx.ShouldBindJSON(&req); err != nil {
			var errors []string
			errs, ok := err.(validator.ValidationErrors)
			if !ok {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			for _, e := range errs {
				switch e.Tag() {
				case "required":
					errors = append(errors, e.Field()+" is required")
				case "email":
					errors = append(errors, "Password must be at least 6 characters long")
				}
			}
			ctx.JSON(http.StatusBadRequest, gin.H{"error": errors})
			return
		}

		user, err := s.DB.GetUser(ctx, username)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User Not Found"})
			return
		}
		if req.FullName == "" {
			req.FullName = user.FullName
		}

		if req.Email == "" {
			req.Email = user.Email
		}

		updatedUser, err := s.DB.UpdateUser(ctx, db.UpdateUserParams{
			Username: user.Username,
			FullName: sql.NullString{
				String: req.FullName,
				Valid:  true,
			},
			Email: sql.NullString{
				String: req.Email,
				Valid:  true,
			},
			UpdatedAt: time.Now(),
		})

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User Cannot Updated"})
			return
		}

		res := NewUserResponse(updatedUser)
		ctx.JSON(http.StatusOK, res)
	}
}
