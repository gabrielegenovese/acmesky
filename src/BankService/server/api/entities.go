package api

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID          uuid.UUID    `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at" gorm:"index"`
	User        string       `json:"user"`
	Description string       `json:"description"`
	Amount      uint32       `json:"amount"`
	Paid        bool         `json:"paid"`
}

type PaymentReq struct {
	User        string `json:"user"`
	Description string `json:"description"`
	Amount      uint32 `json:"amount"`
}

type Res struct {
	Res string `json:"res"`
}