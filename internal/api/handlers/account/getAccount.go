package account

import (
	"database/sql"
	"net/http"

	"github.com/claytten/golang-simplebank/internal/api"
	"github.com/gin-gonic/gin"
)

func GetAccountRoute(api *api.Server, userRg *gin.RouterGroup) {
	userRg.GET("/getAccount", GetAccountHandler(api))
}

type GetAccountRequest struct {
	ID int64 `header:"id" binding:"required,min=1"`
}

func GetAccountHandler(s *api.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req GetAccountRequest
		if err := ctx.ShouldBindHeader(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		account, err := s.DB.GetAccount(ctx, req.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, account)
	}
}
