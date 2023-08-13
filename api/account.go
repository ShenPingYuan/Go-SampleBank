package api

import (
	"database/sql"
	"net/http"

	db "github.com/ShenPingYuan/go-webdemo/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/shopspring/decimal"
)

type CreateAccountRequest struct {
	OwnerID  int64  `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.CreateAccountParams{
		OwnerID:  req.OwnerID,
		Currency: req.Currency,
		Balance:  decimal.NewFromFloat(0),
	}
	result, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if sqlError, ok := err.(*mysql.MySQLError); ok {
			switch sqlError.Number {
			case 1062:
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	account, err := server.store.GetAccount(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//rsp := NewAccountResponse(account)
	ctx.JSON(http.StatusOK, account)
}

type GetAccountRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	account, err := server.store.GetAccount(ctx, req.Id)
	if err == sql.ErrNoRows {
		//使用http.NotFoundHandler()处理404错误
		//ctx.JSON(http.StatusNotFound, errorResponse(err))
		http.Error(ctx.Writer, "account not found", http.StatusNotFound)
		//http.NotFound(ctx.Writer, ctx.Request)
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

// 分页查询请求参数，从查询参数中获取
type GetPagedAccountsRequest struct {
	PageIndex int32 `form:"pageIndex" binding:"required,min=1"`
	PageSize  int32 `form:"pageSize" binding:"required,min=1,max=100"`
}

func (server *Server) getPagedAccounts(ctx *gin.Context) {
	var req GetPagedAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	accounts, err := server.store.ListAccounts(ctx, db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageIndex - 1) * req.PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}

type DeleteAccountRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req DeleteAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err := server.store.DeleteAccount(ctx, req.Id)
	if err == sql.ErrNoRows {
		http.Error(ctx.Writer, "account not found", http.StatusNotFound)
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.Status(http.StatusOK)
}
