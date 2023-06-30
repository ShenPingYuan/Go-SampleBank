package api

import (
	"net/http"

	db "github.com/ShenPingYuan/go-webdemo/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR CNY"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  decimal.NewFromFloat(0),
	}
	result, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	id,err:=result.LastInsertId()
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	account,err:=server.store.GetAccount(ctx,id)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//rsp := NewAccountResponse(account)
	ctx.JSON(http.StatusOK, account)
}
