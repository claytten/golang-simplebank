package auth

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/claytten/golang-simplebank/internal/api"
	"github.com/gin-gonic/gin"
)

type RenewAccessRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RenewAccessResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func PostRenewTokenUserRoute(api *api.Server, userRg *gin.RouterGroup) {
	userRg.POST("/renew-token", PostRenewTokenUserHandler(api))
}

func PostRenewTokenUserHandler(s *api.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req RenewAccessRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		refreshPayload, err := s.Token.VerifyToken(req.RefreshToken)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		session, err := s.DB.GetSession(ctx, refreshPayload.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// session checking
		// check block sessin
		if session.IsBlocked {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "blocked session"})
			return
		}

		// check email is matching
		if session.Email != refreshPayload.Email {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect session user"})
			return
		}

		// check refreshToken is matching
		if session.RefreshToken != req.RefreshToken {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "mismatched session token"})
			return
		}

		// check expired token
		if time.Now().After(session.ExpiresAt) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "expired session"})
			return
		}

		// create new token
		accessToken, accessTokenPayload, err := s.Token.CreateToken(
			refreshPayload.Email,
			s.Config.AccessTokenDuration,
		)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := RenewAccessResponse{
			AccessToken:          accessToken,
			AccessTokenExpiresAt: accessTokenPayload.ExpiredAt,
		}

		ctx.JSON(http.StatusOK, response)
	}
}
