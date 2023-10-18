package token

import (
	"errors"
	"fmt"
	"time"

	db "github.com/ShenPingYuan/go-webdemo/db/sqlc"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

// CreateAccessToken implements Maker.
func (maker *JWTMaker) CreateAccessToken(userId int64, username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(userId, username, duration)
	if err != nil {
		return "", err
	}
	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", err
	}
	return jwtToken, nil
}

// VerifyAccessToken implements Maker.
func (maker *JWTMaker) VerifyAccessToken(token string) (*Payload, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token: %s", token.Header["alg"])
		}
		return []byte(maker.secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return payload, nil
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateRefreshToken(userId int64, duration time.Duration) (*RefreshToken, error) {
	return NewRefreshToken(userId, duration)
}

func (maker *JWTMaker) VerifyRefreshToken(token db.Session) (*RefreshToken, error) {
	// 1. Check if the token is blocked
	if token.IsBlocked {
		return nil, errors.New("token is blocked")
	}

	// 2. Check if the token is expired
	if time.Now().After(token.ExpireTime) {
		return nil, errors.New("token has expired")
	}

	// 3. Check if the refreshToken exists
	if token.RefreshToken == "" {
		return nil, errors.New("refresh token is missing")
	}

	// If all checks pass, return the RefreshToken (you might want to extract or construct it based on the Session)
	id, err := uuid.Parse(token.ID)
	if err != nil {
		return nil, err
	}
	refreshToken := &RefreshToken{
		// ... populate the RefreshToken fields based on the Session or other logic ...
		Id:           id,
		UserId:       token.UserID,
		RefreshToken: token.RefreshToken,
		ExpiredAt:    token.ExpireTime,
	}

	return refreshToken, nil
}
