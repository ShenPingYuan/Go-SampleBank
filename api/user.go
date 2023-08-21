package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/ShenPingYuan/go-webdemo/db/sqlc"
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

type LoginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	UserInfo    LoginUserDto
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
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}
	//生成token
	token, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
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
		AccessToken: token,
		UserInfo:    rspUser,
	}
	ctx.JSON(http.StatusOK, rsp)
}
