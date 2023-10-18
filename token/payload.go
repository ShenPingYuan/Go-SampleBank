package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrExpiredToken = errors.New("token has expired")

type Payload struct {
	Id        uuid.UUID `json:"id"`
	UserId    int64     `json:"user_id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"iat"`
	ExpiredAt time.Time `json:"exp"`
}

type RefreshToken struct {
	Id           uuid.UUID `json:"id"`
	UserId       int64     `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	ExpiredAt    time.Time `json:"exp"`
}

// Valid implements jwt.Claims.
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

func NewPayload(userId int64, username string, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Payload{
		Id:        tokenId,
		UserId:    userId,
		Username:  username,
		IssuedAt:  now,
		ExpiredAt: now.Add(duration),
	}, nil
}

func NewRefreshToken(userId int64, duration time.Duration) (*RefreshToken, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &RefreshToken{
		Id:           tokenId,
		UserId:       userId,
		RefreshToken: uuid.NewString(),
		ExpiredAt:    now.Add(duration),
	}, nil
}
