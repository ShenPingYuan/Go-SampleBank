package db

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	assert := assert.New(t)

	accountId1 := createRandomAccountReturnId(t)
	account1, err := store.GetAccount(testContext, accountId1)
	assert.NoError(err)
	assert.NotEmpty(account1)
	accountId2 := createRandomAccountReturnId(t)
	account2, err := store.GetAccount(testContext, accountId2)
	assert.NoError(err)
	assert.NotEmpty(account2)

	//assert.Error(err)
	fmt.Println("转账前：", account1.Balance, account2.Balance)

	n := 5
	amount := decimal.NewFromFloat(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(testContext, TransferTxParams{
				FromAccountID: accountId1,
				ToAccountID:   accountId2,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool) //用于判断是否已经存在

	// Check results
	for i := 0; i < n; i++ {

		err := <-errs
		assert.NoError(err)
		result := <-results
		assert.NotEmpty(result)

		fmt.Println("第", i, "次转账：", result.FromAccount.Balance, result.ToAccount.Balance)

		// Check transfer
		transfer := result.Transfer
		assert.NotEmpty(transfer)
		assert.Equal(accountId1, transfer.FromAccountID)
		assert.Equal(accountId2, transfer.ToAccountID)
		assert.True(amount.Equal(transfer.Amount))
		assert.NotZero(transfer.ID)
		assert.NotZero(transfer.CreatedAt)

		_, err = store.GetTransfer(testContext, transfer.ID)
		assert.NoError(err)

		// Check entries
		fromEntry := result.FromEntry
		assert.NotEmpty(fromEntry)
		assert.Equal(accountId1, fromEntry.AccountID)
		assert.True(amount.Neg().Equal(fromEntry.Amount))
		assert.NotZero(fromEntry.ID)
		assert.NotZero(fromEntry.CreatedAt)

		_, err = store.GetEntry(testContext, fromEntry.ID)
		assert.NoError(err)

		toEntry := result.ToEntry
		assert.NotEmpty(toEntry)
		assert.Equal(accountId2, toEntry.AccountID)
		assert.True(amount.Equal(toEntry.Amount))
		assert.NotZero(toEntry.ID)
		assert.NotZero(toEntry.CreatedAt)

		_, err = store.GetEntry(testContext, toEntry.ID)
		assert.NoError(err)

		// Check account
		fromAccount := result.FromAccount
		assert.NotEmpty(fromAccount)
		assert.Equal(accountId1, fromAccount.ID)

		toAccount := result.ToAccount
		assert.NotEmpty(toAccount)
		assert.Equal(accountId2, toAccount.ID)

		// Check account balance
		diff1 := account1.Balance.Sub(fromAccount.Balance) // 从账户1扣除的金额
		diff2 := toAccount.Balance.Sub(account2.Balance)   // 转入账户2的金额
		assert.Equal(diff1, diff2)
		assert.True(diff1.GreaterThan(decimal.NewFromInt(0)))        // 从账户1扣除的金额大于0
		assert.True(diff1.Mod(amount).Equals(decimal.NewFromInt(0))) // 从账户1扣除的金额是amount的整数倍

		k := diff1.Div(amount) // 从账户1扣除的金额是amount的整数倍
		assert.True(k.GreaterThanOrEqual(decimal.NewFromInt(1)) && k.LessThanOrEqual(decimal.NewFromInt(int64(n))))

		assert.NotContains(existed, int(k.IntPart())) //判断是否已经存在
		existed[int(k.IntPart())] = true              //记录已经存在的倍数
	}

	// Check the final updated balance
	updatedAccount1, err := store.GetAccount(testContext, accountId1)
	assert.NoError(err)
	assert.NotEmpty(updatedAccount1)

	updatedAccount2, err := store.GetAccount(testContext, accountId2)
	assert.NoError(err)
	assert.NotEmpty(updatedAccount2)

	fmt.Println("转账后：", updatedAccount1.Balance, updatedAccount2.Balance)

	assert.Equal(account1.Balance.Sub(decimal.NewFromInt(int64(n)).Mul(amount)), updatedAccount1.Balance) //账户1扣除的金额
	assert.Equal(account2.Balance.Add(decimal.NewFromInt(int64(n)).Mul(amount)), updatedAccount2.Balance) //账户2增加的金额

}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)
	assert := assert.New(t)

	accountId1 := createRandomAccountReturnId(t)
	account1, err := store.GetAccount(testContext, accountId1)
	assert.NoError(err)
	assert.NotEmpty(account1)
	accountId2 := createRandomAccountReturnId(t)
	account2, err := store.GetAccount(testContext, accountId2)
	assert.NoError(err)
	assert.NotEmpty(account2)

	//assert.Error(err)
	fmt.Println("转账前：", account1.Balance, account2.Balance)

	n := 2
	amount := decimal.NewFromFloat(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountId := accountId1
		toAccountId := accountId2
		if i%2 == 1 {
			fromAccountId = accountId2
			toAccountId = accountId1
		}
		go func() {
			_, err := store.TransferTx(testContext, TransferTxParams{
				FromAccountID: fromAccountId,
				ToAccountID:   toAccountId,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {

		err := <-errs
		assert.NoError(err)
	}

	// Check the final updated balance
	updatedAccount1, err := store.GetAccount(testContext, accountId1)
	assert.NoError(err)
	assert.NotEmpty(updatedAccount1)

	updatedAccount2, err := store.GetAccount(testContext, accountId2)
	assert.NoError(err)
	assert.NotEmpty(updatedAccount2)

	fmt.Println("转账后：", updatedAccount1.Balance, updatedAccount2.Balance)

	assert.Equal(account1.Balance, updatedAccount1.Balance) //账户1扣除的金额
	assert.Equal(account2.Balance, updatedAccount2.Balance) //账户2增加的金额

}
