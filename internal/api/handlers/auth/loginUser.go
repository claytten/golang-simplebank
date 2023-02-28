package auth

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/claytten/golang-simplebank/internal/api"
	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func NewUserResponse(user db.Users) *UserResponse {
	return &UserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
	}
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  UserResponse `json:"user"`
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

		accessToken, accessPayload, err := s.Token.CreateToken(user.Email, s.Config.AccessTokenDuration)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		refreshToken, refreshPayload, err := s.Token.CreateToken(user.Email, s.Config.RefreshTokenDuration)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		session, err := s.DB.CreateSession(ctx, db.CreateSessionParams{
			ID:           refreshPayload.ID,
			Email:        user.Email,
			RefreshToken: refreshToken,
			UserAgent:    ctx.Request.UserAgent(),
			ClientIp:     ctx.ClientIP(),
			IsBlocked:    false,
			ExpiresAt:    refreshPayload.ExpiredAt,
			CreatedAt:    time.Now(),
		})

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		res := loginUserResponse{
			SessionID:             session.ID,
			AccessToken:           accessToken,
			AccessTokenExpiresAt:  accessPayload.ExpiredAt,
			RefreshToken:          refreshToken,
			RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
			User:                  *NewUserResponse(user),
		}

		ctx.JSON(http.StatusOK, res)
	}
}
