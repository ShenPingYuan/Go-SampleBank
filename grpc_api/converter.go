package grpc_api

import (
	db "github.com/ShenPingYuan/go-webdemo/db/sqlc"
	pb "github.com/ShenPingYuan/go-webdemo/protobuffer"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(dbUser db.User) *pb.User {
	return &pb.User{
		Username:          dbUser.Username,
		FullName:          dbUser.FullName,
		Email:             dbUser.Email,
		PasswordChangedAt: timestamppb.New(dbUser.PasswordChangedAt),
		CreatedAt:         timestamppb.New(dbUser.CreatedAt),
	}
}
