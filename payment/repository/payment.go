package repository

import (
	"context"
	"time"

	"github.com/huangyul/go-webook/payment/domain"
	"github.com/huangyul/go-webook/payment/repository/dao"
)

type PaymentRepositoryImp struct {
	dao dao.PaymentDAO
}

// FindExpiredPayment
func (p *PaymentRepositoryImp) FindExpiredPayment(ctx context.Context, doffset int, limit int, t time.Time) ([]*domain.Payment, error) {
	res, err := p.dao.FindExpiredPayment(ctx, doffset, limit, t)
	if err != nil {
		return nil, err
	}
	payments := make([]*domain.Payment, 0, len(res))
	for _, r := range res {
		payments = append(payments, p.toDomain(&r))
	}
	return payments, nil
}

// GetPayment
func (p *PaymentRepositoryImp) GetPayment(ctx context.Context, BizTranNo string) (*domain.Payment, error) {
	res, err := p.dao.GetPayment(ctx, BizTranNo)
	if err != nil {
		return nil, err
	}
	return p.toDomain(res), nil
}

// Insert
func (p *PaymentRepositoryImp) Insert(ctx context.Context, pm *domain.Payment) error {
	return p.dao.Insert(ctx, p.toEntity(pm))
}

// UpdateTxnIDAndStatus
func (p *PaymentRepositoryImp) UpdateTxnIDAndStatus(ctx context.Context, bizTranNo string, TxnID string, status uint8) error {
	return p.dao.UpdateTxnIDAndStatus(ctx, bizTranNo, TxnID, status)
}

func NewPaymentRepositoryImp(dao dao.PaymentDAO) PaymentRepository {
	return &PaymentRepositoryImp{
		dao: dao,
	}
}

func (p *PaymentRepositoryImp) toDomain(payment *dao.Payment) *domain.Payment {
	return &domain.Payment{
		ID:        payment.ID,
		BizTranNo: payment.BizTranNo,
		TxnID:     payment.TxnID,
		Currency:  payment.Currency,
		CreatedAt: payment.CreatedAt,
		UpdatedAt: payment.UpdatedAt,
	}
}

func (p *PaymentRepositoryImp) toEntity(payment *domain.Payment) *dao.Payment {
	return &dao.Payment{
		ID:        payment.ID,
		BizTranNo: payment.BizTranNo,
		TxnID:     payment.TxnID,
		Currency:  payment.Currency,
	}
}
