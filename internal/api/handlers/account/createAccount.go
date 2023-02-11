package account

import (
	"net/http"

	"github.com/claytten/golang-simplebank/internal/api"
	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

const (
	authorizationHeaderKey   = "authorization"
	authorizationTypeBearer  = "bearer"
	authorizationPayloadKey  = "authorization_payload"
	authorizationPassword    = "password"
	authorizationOldPassword = "oldPassword"
	authorizationUsername    = "username"
)

func PostCreateAccountRoute(api *api.Server, userRg *gin.RouterGroup) {
	userRg.POST("/create", PostCreateAccountHandler(api))
}

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func PostCreateAccountHandler(s *api.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req createAccountRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Some params are not found"})
			return
		}
		username := ctx.GetHeader(authorizationUsername)

		arg := db.CreateAccountParams{
			Owner:    username,
			Currency: req.Currency,
			Balance:  0,
		}

		account, err := s.DB.CreateAccount(ctx, arg)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				switch pqErr.Code.Name() {
				case "foreign_key_violation", "unique_violation":
					ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
					return
				}
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot Create Account"})
			return
		}

		ctx.JSON(http.StatusOK, account)
	}
}
