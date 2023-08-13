package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPassword(t *testing.T) {
	assert := assert.New(t)

	randomPassword:=RandomString(8)
	hashedPassword, err := HashPassword(randomPassword)

	assert.NoError(err)
	assert.NotEmpty(hashedPassword)
	assert.NotEqual(randomPassword, hashedPassword)

	ok := CheckPassword(randomPassword, hashedPassword)
	assert.True(ok)

	wrongPassword := RandomString(8)
	ok = CheckPassword(wrongPassword, hashedPassword)
	assert.False(ok)
}
