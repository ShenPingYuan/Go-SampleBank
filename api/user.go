package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	db "github.com/ShenPingYuan/go-webdemo/db/sqlc"
	"github.com/ShenPingYuan/go-webdemo/token"
	"github.com/ShenPingYuan/go-webdemo/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type ResponseUser struct {
	Username         string    `json:"username"`
	FullName         string    `json:"fullname"`
	Email            string    `json:"email"`
	PasswordChangeAt time.Time `json:"password_change_at"`
	CreatedAt        time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}
	result, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	user, err := server.store.GetUserById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := ResponseUser{
		Username:         user.Username,
		FullName:         user.FullName,
		Email:            user.Email,
		PasswordChangeAt: user.PasswordChangedAt,
		CreatedAt:        user.CreatedAt,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type UserDetailDto struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	// HashedPassword    string    `json:"hashedPassword"`
	Email             string    `json:"email"`
	FullName          string    `json:"fullName"`
	PasswordChangedAt time.Time `json:"passwordChangedAt"`
	CreatedAt         time.Time `json:"createdAt"`
}

func (server *Server) getUserInfo(ctx *gin.Context) {
	currentUser := ctx.MustGet(currentUserKey).(*token.Payload)

	user, err := server.store.GetUserById(ctx, currentUser.UserId)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := UserDetailDto{
		ID:                user.ID,
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}

type LoginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserInfo     LoginUserDto
}
type LoginUserDto struct {
	Username string `json:"username"`
	FullName string `json:"fullname"`
	Email    string `json:"email"`
}

func (server *Server) login(ctx *gin.Context) {
	var request LoginUserRequest
	if err := ctx.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUserByUsername(ctx, request.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if !util.CheckPassword(request.Password, user.HashedPassword) {
		//密码错误
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("invalid user information")))
		return
	}
	//生成token
	token, err := server.tokenMaker.CreateAccessToken(user.ID, user.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, err := server.tokenMaker.CreateRefreshToken(user.ID, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// 把refreshToken存到数据库
	_, err = server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshToken.Id.String(),
		UserID:       user.ID,
		RefreshToken: refreshToken.RefreshToken,
		UserAgent:    ctx.Request.UserAgent(), //TODO
		ClientIp:     ctx.ClientIP(), //TODO
		IsBlocked:    false,
		ExpireTime:   refreshToken.ExpiredAt,
		CreatedAt:    time.Now(),
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rspUser := LoginUserDto{
		Username: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
	}

	rsp := LoginResponse{
		RefreshToken: refreshToken.RefreshToken,
		AccessToken:  token,
		UserInfo:     rspUser,
	}
	ctx.JSON(http.StatusOK, rsp)
}
