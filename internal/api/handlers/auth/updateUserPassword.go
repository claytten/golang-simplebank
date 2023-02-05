package auth

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/claytten/golang-simplebank/internal/api"
	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/gin-gonic/gin"
)

func UpdateUserPasswordRoute(api *api.Server, userRg *gin.RouterGroup) {
	userRg.PUT("/profile/password", UpdateUserPasswordHandler(api))
}

type UpdateUserPasswordRequest struct {
	Username string `header:"username" binding:"required,alphanum"`
	Password string `header:"password" binding:"required,min=6"`
}

func UpdateUserPasswordHandler(s *api.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req UpdateUserPasswordRequest
		username := ctx.GetHeader(authorizationUsername)

		// checking username on header
		if err := ctx.ShouldBindHeader(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := s.DB.GetUser(ctx, username)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User Not Found"})
			return
		}

		newPassword, err := util.HashingPassword(req.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Password can't use. Please Try Another Password"})
			return
		}

		updatedUser, err := s.DB.UpdateUser(ctx, db.UpdateUserParams{
			Username: user.Username,
			HashedPassword: sql.NullString{
				String: newPassword,
				Valid:  true,
			},
			PasswordChangedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
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
