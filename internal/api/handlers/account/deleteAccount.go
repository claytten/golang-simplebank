package account

import (
	"net/http"

	"github.com/claytten/golang-simplebank/internal/api"
	"github.com/gin-gonic/gin"
)

func DeleteAccountRoute(api *api.Server, userRg *gin.RouterGroup) {
	userRg.DELETE("/delete", DeleteAccountHandler(api))
}

func DeleteAccountHandler(s *api.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req GetAccountRequest

		if err := ctx.ShouldBindHeader(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := s.DB.DeleteAccount(ctx, req.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot delete account"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Delete Successfully",
		})
	}
}
