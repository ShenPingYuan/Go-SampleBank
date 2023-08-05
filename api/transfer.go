package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/ShenPingYuan/go-webdemo/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type transferRequest struct {
	FromAccountID int64   `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64   `json:"to_account_id" binding:"required,min=1"`
	Amount        float64 `json:"amount" binding:"required,numeric,gt=0"`
	Currency      string  `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.validAccount(ctx, req.FromAccountID, req.Currency) || !server.validAccount(ctx, req.ToAccountID, req.Currency) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        decimal.NewFromFloat(req.Amount),
	}
	result, err := server.store.TransferTx(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountId int64, Currency string) bool {
	account, err := server.store.GetAccount(ctx, accountId)
	if err == sql.ErrNoRows {
		//使用http.NotFoundHandler()处理404错误
		//ctx.JSON(http.StatusNotFound, errorResponse(err))
		http.Error(ctx.Writer, "account not found", http.StatusNotFound)
		//http.NotFound(ctx.Writer, ctx.Request)
		return false
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}
	if account.Currency != Currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountId, account.Currency, Currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true
}
