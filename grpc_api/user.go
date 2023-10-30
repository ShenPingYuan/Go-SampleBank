package grpc_api

import (
	"context"
	"database/sql"
	"time"

	db "github.com/ShenPingYuan/go-webdemo/db/sqlc"
	pb "github.com/ShenPingYuan/go-webdemo/protobuffer"
	"github.com/ShenPingYuan/go-webdemo/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot hash password,err:%s", err)
	}
	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}
	result, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create user,err:%s", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get last insert id,err:%s", err)
	}
	user, err := server.store.GetUserById(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get user by id,err:%s", err)
	}
	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}

func (server *Server) LoginUser(ctx context.Context, request *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUserByUsername(ctx, request.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not exists,err:%s", err)
		}
		return nil, status.Errorf(codes.Internal, "cannot get user by username,err:%s", err)
	}
	if !util.CheckPassword(request.Password, user.HashedPassword) {
		//密码错误
		return nil, status.Errorf(codes.Unauthenticated, "incorrect password,err:%s", err)
	}
	//生成token
	token, err := server.tokenMaker.CreateAccessToken(user.ID, user.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create access token,err:%s", err)
	}

	refreshToken, err := server.tokenMaker.CreateRefreshToken(user.ID, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create refresh token,err:%s", err)
	}
	// 把refreshToken存到数据库
	_, err = server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshToken.Id.String(),
		UserID:       user.ID,
		RefreshToken: refreshToken.RefreshToken,
		UserAgent:    "", //ctx.Request.UserAgent(), //TODO
		ClientIp:     "", //ctx.ClientIP(),          //TODO
		IsBlocked:    false,
		ExpireTime:   refreshToken.ExpiredAt,
		CreatedAt:    time.Now(),
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create session,err:%s", err)
	}

	rsp := &pb.LoginUserResponse{
		User:         convertUser(user),
		AccessToken:  token,
		RefreshToken: refreshToken.RefreshToken,
	}
	return rsp, nil
}
