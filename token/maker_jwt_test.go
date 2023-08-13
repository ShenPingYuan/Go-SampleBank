package token

import (
	"log"
	"testing"
	"time"

	"github.com/ShenPingYuan/go-webdemo/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestJWTMaker(t *testing.T) {
	assert := assert.New(t)

	maker, err := NewJWTMaker(util.RandomString(32))
	assert.NoError(err)
	assert.NotEmpty(maker)

	duration := time.Minute

	issueAt := time.Now()
	expiredAt := issueAt.Add(duration)

	username := util.RandomString(16)

	// Create a token
	token, err := maker.CreateToken(username, duration)
	log.Println(token)
	assert.NoError(err)
	assert.NotEmpty(token)

	// Verify the token
	payload, err := maker.VerifyToken(token)
	assert.NoError(err)
	assert.NotEmpty(payload)

	// Verify the token's payload
	assert.Equal(payload.Username, username)
	assert.WithinDuration(payload.IssuedAt, issueAt, time.Second)
	assert.WithinDuration(payload.ExpiredAt, expiredAt, time.Second)

}

func TestExpiredJWTToken(t *testing.T) {
	assert := assert.New(t)

	maker, err := NewJWTMaker(util.RandomString(32))
	assert.NoError(err)
	assert.NotEmpty(maker)

	duration := -time.Minute

	username := util.RandomString(16)

	// Create a token
	token, err := maker.CreateToken(username, duration)
	log.Println(token)
	assert.NoError(err)
	assert.NotEmpty(token)

	// Verify the token
	payload, err := maker.VerifyToken(token)
	assert.Error(err)
	assert.EqualError(err, ErrExpiredToken.Error())
	assert.Empty(payload)
}

func TestInvalidJWTAlgNone(t *testing.T) {
	payload, err := NewPayload(util.RandomString(16), time.Minute)
	assert.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token,err:=jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	assert.NoError(t, err)
	assert.NotEmpty(t, jwtToken)

	maker, err := NewJWTMaker(util.RandomString(32))
	assert.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	log.Println(err)
	assert.Error(t, err)
	assert.Empty(t, payload)

}
