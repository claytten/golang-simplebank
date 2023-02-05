package auth

import (
	"database/sql"
	"net/http"

	"github.com/claytten/golang-simplebank/internal/api"
	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/gin-gonic/gin"
)

type loginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

func PostLoginUserRoute(api *api.Server, userRg *gin.RouterGroup) {
	userRg.POST("/login", PostLoginUserHandler(api))
}

func PostLoginUserHandler(s *api.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req loginUserRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Some Params are not found"})
			return
		}

		user, err := s.DB.GetUserUsingEmail(ctx, req.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot Login!. username/email or password is not match."})
			return
		}

		err = util.ComparePassword(user.HashedPassword, req.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot Login!. username/email or password is not match."})
			return
		}

		accessToken, _, err := s.Token.CreateToken(user.Email, s.Config.AccessTokenDuration)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		refreshToken, _, err := s.Token.CreateToken(user.Email, s.Config.AccessTokenDuration)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		res := loginUserResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			User:         *NewUserResponse(user),
		}

		ctx.JSON(http.StatusOK, res)
	}
}
