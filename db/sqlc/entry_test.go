package db

import (
	"context"
	"testing"

	"github.com/ShenPingYuan/go-webdemo/util"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// 创建测试条目
func createRandomEntry(t *testing.T) int64 {
	accountId := createRandomAccountReturnId(t)
	arg := CreateEntryParams{
		AccountID: accountId,
		Amount:    util.RandomMoney(),
	}
	result, err := testQueries.CreateEntry(context.Background(), arg)
	if err != nil {
		t.Error(err)
	}
	dbId, err := result.LastInsertId()
	if err != nil {
		t.Error(err)
	}
	return dbId
}

// 测试创建条目
func TestCreateEntry(t *testing.T) {
	assert := assert.New(t)

	account, err := testQueries.GetAccount(context.Background(), 1)
	var id int64
	if err != nil {
		id = createRandomAccountReturnId(t)
	} else {
		id = account.ID
	}
	arg := CreateEntryParams{
		AccountID: id,
		Amount:    decimal.NewFromInt(util.RandomInt(1, 1000)),
	}

	result, err := testQueries.CreateEntry(context.Background(), arg)
	assert.NoError(err)
	assert.NotEmpty(result)
	dbId, err := result.LastInsertId()
	assert.NoError(err)
	assert.NotZero(dbId)
}

// 测试获取条目
func TestGetEntry(t *testing.T) {
	assert := assert.New(t)

	id := createRandomEntry(t)

	entry, err := testQueries.GetEntry(context.Background(), id)
	assert.NoError(err)
	assert.NotEmpty(entry)
	assert.Equal(id, entry.ID)
	assert.NotZero(entry.AccountID)
	assert.NotZero(entry.Amount)
	assert.NotZero(entry.CreatedAt)
}

// 测试删除条目
func TestDeleteEntry(t *testing.T) {
	assert := assert.New(t)

	id := createRandomEntry(t)

	_,err := testQueries.DeleteEntry(context.Background(), id)
	assert.NoError(err)

	entry, err := testQueries.GetEntry(context.Background(), id)
	assert.Error(err)
	assert.Empty(entry)
}