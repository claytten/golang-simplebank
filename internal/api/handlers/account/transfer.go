package account

import (
	"database/sql"
	"net/http"

	"github.com/claytten/golang-simplebank/internal/api"
	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/gin-gonic/gin"
)

func PostCreateTransferAccountRoute(api *api.Server, userRg *gin.RouterGroup) {
	userRg.POST("/transfer", PostCreateTransferAccountHandler(api))
}

type TransferHeaderRequest struct {
	FromAccountID int64 `header:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64 `header:"to_account_id" binding:"required,min=1"`
}

type TransferBodyRequest struct {
	Amount   int64  `json:"amount" binding:"required,gt=0"`
	Currency string `json:"currency" binding:"required,currency"`
}

func PostCreateTransferAccountHandler(s *api.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body TransferBodyRequest
		var head TransferHeaderRequest

		if err := ctx.ShouldBindHeader(&head); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ok := head.ValidAccount(ctx, s.DB, body.Currency)
		if !ok {
			return
		}

		arg := db.TransferTxParams{
			FromAccountID: head.FromAccountID,
			ToAccountID:   head.ToAccountID,
			Amount:        body.Amount,
		}

		result, err := s.DB.TransferTx(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, result)
	}
}

func (h *TransferHeaderRequest) ValidAccount(ctx *gin.Context, db db.Store, currency string) bool {
	fromAccount, err := db.GetAccount(ctx, h.FromAccountID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return false
	}

	toAccount, err := db.GetAccount(ctx, h.ToAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return false
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return false
	}

	if fromAccount.Currency != currency || toAccount.Currency != currency {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "From/To Account Mismatch Currency"})
		return false
	}
	return true
}
