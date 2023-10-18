package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type RenewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RenewAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	var request RenewAccessTokenRequest
	if err := ctx.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	session, err := server.store.GetSessionByToken(ctx, request.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshToken, err := server.tokenMaker.VerifyRefreshToken(session)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUserById(ctx, refreshToken.UserId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	//生成token
	token, err := server.tokenMaker.CreateAccessToken(user.ID, user.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := RenewAccessTokenResponse{
		AccessToken: token,
	}
	ctx.JSON(http.StatusOK, rsp)
}
