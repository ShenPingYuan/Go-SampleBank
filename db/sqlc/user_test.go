package db

import (
	"context"
	"testing"

	"github.com/ShenPingYuan/go-webdemo/util"
	"github.com/stretchr/testify/assert"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.RandomUsername(),
		HashedPassword: util.RandomPassword(),
		Email:          util.RandomEmail(),
		FullName:       util.RandomUsername(),
	}
	createResult, err := testQueries.CreateUser(context.Background(), arg)
	if err != nil {
		t.Error(err)
	}
	id, err := createResult.LastInsertId()
	if err != nil {
		t.Error(err)
	}
	user := User{
		ID:             id,
		Username:       arg.Username,
		HashedPassword: arg.HashedPassword,
		Email:          arg.Email,
		FullName:       arg.FullName,
	}
	return user
}

// 测试创建用户
func TestCreateUser(t *testing.T) {
	assert := assert.New(t)
	arg := CreateUserParams{
		Username:       util.RandomUsername(),
		HashedPassword: util.RandomPassword(),
		Email:          util.RandomEmail(),
		FullName:       util.RandomUsername(),
	}

	result, err := testQueries.CreateUser(context.Background(), arg)

	if err != nil {
		t.Error(err)
	}
	if result == nil {
		t.Error("result is nil")
	}
	var userId, _ = result.LastInsertId()
	assert.NoError(err)
	assert.NotEmpty(userId)
	assert.NotZero(result)
}

// 测试获取用户
func TestGetUserById(t *testing.T) {
	assert := assert.New(t)
	user := createRandomUser(t)
	userInDb, err := testQueries.GetUserById(context.Background(), user.ID)
	assert.NoError(err)
	assert.NotEmpty(userInDb)
	assert.Equal(user.ID, userInDb.ID)
	assert.Equal(user.Username, userInDb.Username)
	assert.Equal(user.HashedPassword, userInDb.HashedPassword)
	assert.Equal(user.Email, userInDb.Email)
	assert.Equal(user.FullName, userInDb.FullName)
}

// 测试获取用户
func TestGetUserByUsername(t *testing.T) {
	assert := assert.New(t)
	user := createRandomUser(t)

	userBefore, err := testQueries.GetUserById(context.Background(), user.ID)

	assert.NotNil(userBefore)
	assert.NoError(err)
	userLater, err := testQueries.GetUserByUsername(context.Background(), userBefore.Username)
	assert.NoError(err)
	assert.NotEmpty(userLater)
	assert.Equal(user.ID, userLater.ID)
	assert.Equal(userBefore.ID, userLater.ID)
}
