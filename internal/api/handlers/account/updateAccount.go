package account

import (
	"net/http"
	"time"

	"github.com/claytten/golang-simplebank/internal/api"
	"github.com/claytten/golang-simplebank/internal/api/middlewares"
	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/gin-gonic/gin"
)

func PutUpdateAccountRoute(api *api.Server, userRg *gin.RouterGroup) {
	userRg.Use(middlewares.CheckOwnUserUpdate(api.DB))
	userRg.PUT("/update", PutUpdateAccountHandler(api))
}

type UpdateAccountRequest struct {
	ID      int64 `json:"id"`
	Balance int64 `json:"balance"`
}

func PutUpdateAccountHandler(s *api.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req UpdateAccountRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		username := ctx.GetHeader(authorizationUsername)

		args := db.UpdateAccountParams{
			ID:        req.ID,
			Balance:   req.Balance,
			UpdatedAt: time.Now(),
		}

		account, err := s.DB.UpdateAccount(ctx, args)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot update balance account"})
			return
		}

		if account.Owner != username {
			err := s.DB.DeleteAccount(ctx, account.ID)
			if err != nil {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "account doesn't belong to authenticated user"})
				return
			}
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "account doesn't belong to authenticated user"})
			return
		}

		ctx.JSON(http.StatusOK, account)
	}
}
