package util

import (
	"math/rand"
	"time"

	"github.com/shopspring/decimal"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	//设置随机种子
	rand.Seed(time.Now().UnixNano())
}

// RandomInt 生成随机数,min和max都包含
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of lenght n
func RandomString(n int) string {
	letterRunes := []rune(alphabet)
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generates a random amount of money
func RandomMoney() decimal.Decimal {
	return decimal.NewFromInt(RandomInt(1, 1000))
}

// RandomCurrency generates a random currency code
func RandomCurrency() string {
	currencies := []string{"CNY", "USD", "EUR"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
