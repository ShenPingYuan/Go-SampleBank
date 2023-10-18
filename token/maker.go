package token

import (
	"time"

	db "github.com/ShenPingYuan/go-webdemo/db/sqlc"
)

type Maker interface {
	//创建token
	CreateAccessToken(userId int64, username string, duration time.Duration) (string, error)
	//验证token
	VerifyAccessToken(token string) (*Payload, error)
	//创建RefreshToken
	CreateRefreshToken(userId int64, duration time.Duration) (*RefreshToken, error)
	VerifyRefreshToken(token db.Session) (*RefreshToken, error)
}
