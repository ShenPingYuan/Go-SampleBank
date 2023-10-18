// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	ID        int64           `json:"id"`
	Balance   decimal.Decimal `json:"balance"`
	Currency  string          `json:"currency"`
	CreatedAt time.Time       `json:"createdAt"`
	OwnerID   int64           `json:"ownerID"`
}

type Entry struct {
	ID        int64 `json:"id"`
	AccountID int64 `json:"accountID"`
	// 可以是正或者负，表示存入或者取出
	Amount    decimal.Decimal `json:"amount"`
	CreatedAt time.Time       `json:"createdAt"`
}

type Session struct {
	ID           string    `json:"id"`
	UserID       int64     `json:"userID"`
	RefreshToken string    `json:"refreshToken"`
	UserAgent    string    `json:"userAgent"`
	ClientIp     string    `json:"clientIp"`
	IsBlocked    bool      `json:"isBlocked"`
	ExpireTime   time.Time `json:"expireTime"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Transfer struct {
	ID            int64 `json:"id"`
	FromAccountID int64 `json:"fromAccountID"`
	// Content of the post
	ToAccountID int64 `json:"toAccountID"`
	// 只能是正数
	Amount    decimal.Decimal `json:"amount"`
	CreatedAt time.Time       `json:"createdAt"`
}

type User struct {
	ID                int64     `json:"id"`
	Username          string    `json:"username"`
	HashedPassword    string    `json:"hashedPassword"`
	Email             string    `json:"email"`
	FullName          string    `json:"fullName"`
	PasswordChangedAt time.Time `json:"passwordChangedAt"`
	CreatedAt         time.Time `json:"createdAt"`
}
