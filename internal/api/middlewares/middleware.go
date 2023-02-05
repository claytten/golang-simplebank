package middlewares

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/claytten/golang-simplebank/internal/api/token"
	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
	authorizationPassword   = "password"
	authorizationUsername   = "username"
)

func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		//split string into slice
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

type UpdateUserPasswordRequest struct {
	Username    string `header:"username" binding:"required,alphanum"`
	OldPassword string `header:"oldPassword" binding:"required,min=6"`
}

func CheckOwnUserUpdate(db db.Store) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
		var req UpdateUserPasswordRequest

		// checking username on header
		if err := ctx.ShouldBindHeader(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Some params header not provided"})
			return
		}

		//finding user by email
		userHeader, err := db.GetUserUsingEmail(ctx, authPayload.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User Not Found and Not Authorized"})
				return
			}
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "account doesn't belong to authenticated user"})
			return
		}

		// checking if username is provided at header
		if req.Username != userHeader.Username {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User Status Unauthorized"})
			return
		}

		// checking if user typing old password and new password is same
		err = util.ComparePassword(userHeader.HashedPassword, req.OldPassword)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "new password / old password doesn't match"})
			return
		}

		ctx.Set(authorizationUsername, userHeader.Username)
		ctx.Next()
	}
}
