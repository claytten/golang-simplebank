package account

import (
	"net/http"
	"time"

	"github.com/claytten/golang-simplebank/internal/api"
	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/gin-gonic/gin"
)

func PutUpdateAccountRoute(api *api.Server, userRg *gin.RouterGroup) {
	userRg.PUT("/update", PutUpdateAccountHandler(api))
}

type UpdateAccountRequest struct {
	Balance int64 `json:"balance" binding:"required"`
}

func PutUpdateAccountHandler(s *api.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req UpdateAccountRequest
		var id GetAccountRequest

		if err := ctx.ShouldBindHeader(&id); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		args := db.UpdateAccountParams{
			ID:        id.ID,
			Balance:   req.Balance,
			UpdatedAt: time.Now(),
		}

		account, err := s.DB.UpdateAccount(ctx, args)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot update balance account"})
			return
		}

		ctx.JSON(http.StatusOK, account)
	}
}
