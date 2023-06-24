package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ShenPingYuan/go-webdemo/util"
	"github.com/stretchr/testify/assert"
)

func createRandomAccountReturnId(t *testing.T) int64 {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	createResult, err := testQueries.CreateAccount(context.Background(), arg)
	if err != nil {
		t.Error(err)
	}
	id, err := createResult.LastInsertId()
	if err != nil {
		t.Error(err)
	}
	return id
}

// 测试获取最后插入的Id
func TestGetLastInsertId(t *testing.T) {
	assert := assert.New(t)
	id := createRandomAccountReturnId(t)
	//根据Id查询账户
	account, err := testQueries.GetAccount(context.Background(), id)

	id2 := createRandomAccountReturnId(t)

	assert.NoError(err)
	assert.NotEmpty(account)
	assert.Equal(id, account.ID)
	assert.True(id2 > id)
}

// 测试创建账户
func TestCreateAccount(t *testing.T) {
	assert := assert.New(t)
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	var result, _ = account.LastInsertId()
	assert.NoError(err)
	assert.NotEmpty(account)
	assert.NotZero(result)
}

// 测试获取账户
func TestGetAccount(t *testing.T) {
	assert := assert.New(t)
	id := createRandomAccountReturnId(t)
	//根据Id查询账户
	account, err := testQueries.GetAccount(context.Background(), id)
	assert.NoError(err)
	assert.NotEmpty(account)
	assert.NotZero(account.Balance)
	assert.NotZero(account.ID)
	assert.NotEmpty(account.Owner)
	assert.NotEmpty(account.Currency)
	assert.NotZero(account.CreatedAt)
	assert.Equal(id, account.ID)
}

// 测试更新账户余额
func TestUpdateAccountBalance(t *testing.T) {
	assert := assert.New(t)
	id := createRandomAccountReturnId(t)
	arg := UpdateAccountBalanceParams{
		ID:      id,
		Balance: util.RandomMoney(),
	}
	//更新账户余额
	updateResult, err := testQueries.UpdateAccountBalance(context.Background(), arg)
	assert.NoError(err)
	//查询更新后的账户
	account, err := testQueries.GetAccount(context.Background(), id)
	assert.NoError(err)
	assert.NotEmpty(account)
	assert.Equal(id, account.ID)
	argBalance, err := arg.Balance.Value()
	assert.NoError(err)
	accountBalance, err := account.Balance.Value()
	assert.NoError(err)
	assert.Equal(argBalance, accountBalance)

	affectedRows, err := updateResult.RowsAffected()
	assert.NoError(err)
	assert.NotEmpty(updateResult)
	assert.Equal(int64(1), affectedRows)
}

// 测试删除账户
func TestDeleteAccount(t *testing.T) {
	assert := assert.New(t)
	id := createRandomAccountReturnId(t)
	//删除账户
	result, err := testQueries.DeleteAccount(context.Background(), id)
	assert.NoError(err)
	affectedRows, err := result.RowsAffected()
	assert.NoError(err)
	assert.NotEmpty(result)
	assert.Equal(int64(1), affectedRows)
	//查询账户
	account, err := testQueries.GetAccount(context.Background(), id)
	assert.Error(err)
	assert.EqualError(err, sql.ErrNoRows.Error())
	assert.Empty(account)
}

// 测试列出所有账户
func TestListAccounts(t *testing.T) {
	assert := assert.New(t)
	//创建账户
	for i := 0; i < 10; i++ {
		createRandomAccountReturnId(t)
	}
	//列出所有账户
	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	assert.NoError(err)
	assert.NotEmpty(accounts)
	assert.Len(accounts, 5)
	for _, account := range accounts {
		assert.NotEmpty(account)
		assert.NotZero(account.Balance)
		assert.NotZero(account.ID)
		assert.NotEmpty(account.Owner)
		assert.NotEmpty(account.Currency)
		assert.NotZero(account.CreatedAt)
	}
}
