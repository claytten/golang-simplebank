package auth

import (
	"net/http"

	"github.com/claytten/golang-simplebank/internal/api"
	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func PostCreateUserRoute(api *api.Server, userRg *gin.RouterGroup) {
	userRg.POST("/create", PostCreateUserHandler(api))
}

type UserCreateRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func PostCreateUserHandler(s *api.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req UserCreateRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hashedPassword, err := util.HashingPassword(req.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Password Can't Use"})
			return
		}

		user, err := s.DB.CreateUser(ctx, db.CreateUserParams{
			Username:       req.Username,
			Email:          req.Email,
			FullName:       req.FullName,
			HashedPassword: hashedPassword,
		})

		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				switch pqErr.Code.Name() {
				case "unique_violation":
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User Cannot use. Please Try again"})
					return
				}
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User Already Use"})
			return
		}

		res := NewUserResponse(user)
		ctx.JSON(http.StatusOK, res)
	}
}
