package auth

import (
	"net/http"

	"github.com/claytten/golang-simplebank/internal/api"
	"github.com/claytten/golang-simplebank/internal/api/token"
	"github.com/gin-gonic/gin"
)

func GetUserRoute(api *api.Server, userRg *gin.RouterGroup) {
	userRg.GET("/getUser", GetUserHandler(api))
}

type GetUserRequest struct {
	Username string `header:"username" binding:"required,alphanum"`
}

func GetUserHandler(s *api.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req GetUserRequest

		if err := ctx.ShouldBindHeader(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
		userHeader, err := s.DB.GetUserUsingEmail(ctx, authPayload.Email)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "account doesn't belong to authenticated user"})
			return
		}

		if userHeader.Username != req.Username {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User Status Unauthorized"})
			return
		}

		user, err := s.DB.GetUser(ctx, req.Username)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User Not Found"})
			return
		}

		res := NewUserResponse(user)
		ctx.JSON(http.StatusOK, res)
	}
}
