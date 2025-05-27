package repository

import (
	"context"
	"time"

	"github.com/huangyul/go-webook/payment/domain"
)

type PaymentRepository interface {
	Insert(ctx context.Context, pm *domain.Payment) error
	GetPayment(ctx context.Context, BizTranNo string) (*domain.Payment, error)
	UpdateTxnIDAndStatus(ctx context.Context, bizTranNo string, TxnID string, status uint8) error
	FindExpiredPayment(ctx context.Context, doffset int, limit int, t time.Time) ([]*domain.Payment, error)
}
