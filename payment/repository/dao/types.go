package dao

import (
	"context"
	"time"
)

type PaymentDAO interface {
	Insert(ctx context.Context, pm *Payment) error
	GetPayment(ctx context.Context, BizTranNo string) (*Payment, error)
	UpdateTxnIDAndStatus(ctx context.Context, bizTranNo string, TxnID string, status uint8) error
	FindExpiredPayment(ctx context.Context, doffset int, limit int, t time.Time) ([]Payment, error)
}

type Payment struct {
	ID        int64 `gorm:"primaryKey;AutoIncrement"`
	Currency  string
	BizTranNo string `gorm:"unique;type:varchar(256)"`
	TxnID     string `gorm:"unique;type:varchar(128);common:第三方业务id"`
	Status    uint8
	UpdatedAt *time.Time
	CreatedAt *time.Time
}
