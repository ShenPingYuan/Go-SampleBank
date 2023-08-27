package token

import "time"

type Maker interface {
	//创建token
	CreateToken(userId int64, username string, duration time.Duration) (string, error)
	//验证token
	VerifyToken(token string) (*Payload, error)
}
