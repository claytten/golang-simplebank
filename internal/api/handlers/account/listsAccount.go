package account

import (
	"database/sql"
	"math"
	"net/http"

	"github.com/claytten/golang-simplebank/internal/api"
	"github.com/claytten/golang-simplebank/internal/api/token"
	db "github.com/claytten/golang-simplebank/internal/db/sqlc"
	"github.com/gin-gonic/gin"
)

func ListsAccountsRoute(api *api.Server, userRg *gin.RouterGroup) {
	userRg.GET("/lists", ListsAccountsHandler(api))
}

type ListsAccountsRequest struct {
	PageID   int64 `form:"page_id" binding:"required,min=1"`
	PageSize int64 `form:"page_size" binding:"required,min=5,max=10"`
}

func ListsAccountsHandler(s *api.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req ListsAccountsRequest
		if err := ctx.ShouldBindQuery(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
		user, err := s.DB.GetUserUsingEmail(ctx, authPayload.Email)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		args := db.ListsAccountsParams{
			Owner:  user.Username,
			Limit:  int32(req.PageSize),
			Offset: (int32(req.PageID) - 1) * int32(req.PageSize),
		}

		accounts, err := s.DB.ListsAccounts(ctx, args)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Accounts not found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error while get account"})
			return
		}

		allAccounts, _ := s.DB.GetTotalPageListsAccounts(ctx, user.Username)
		var total_page float64
		if allAccounts == 0 {
			total_page = 0
		} else if allAccounts <= req.PageSize {
			total_page = 1
		} else {
			total_page = math.Ceil(float64(allAccounts) / float64(args.Limit))
		}

		ctx.JSON(http.StatusOK, gin.H{
			"total_page":   total_page,
			"current_page": req.PageID,
			"limit":        req.PageSize,
			"data":         accounts,
		})
	}
}
