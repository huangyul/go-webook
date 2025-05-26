package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var _ PaymentDAO = (*GORMPaymentDAO)(nil)

type GORMPaymentDAO struct {
	db *gorm.DB
}

// FindExpiredPayment
func (dao *GORMPaymentDAO) FindExpiredPayment(ctx context.Context, offset int, limit int, t time.Time) ([]Payment, error) {
	var res []Payment
	er := dao.db.WithContext(ctx).Model(&Payment{}).Where("updated_at < ?", t).Offset(offset).Limit(limit).Find(&res).Error
	return res, er
}

// GetPayment
func (dao *GORMPaymentDAO) GetPayment(ctx context.Context, BizTranNo string) (*Payment, error) {
	var p *Payment
	err := dao.db.WithContext(ctx).Model(&Payment{}).Where("biz_tran_no = ?", BizTranNo).First(p).Error
	return p, err
}

// Insert
func (dao *GORMPaymentDAO) Insert(ctx context.Context, pm *Payment) error {
	now := time.Now()
	pm.UpdatedAt = &now
	pm.CreatedAt = &now
	return dao.db.WithContext(ctx).Create(pm).Error
}

// UpdateTxnIDAndStatus
func (dao *GORMPaymentDAO) UpdateTxnIDAndStatus(ctx context.Context, bizTranNo string, TxnID string, status uint8) error {
	now := time.Now()
	return dao.db.WithContext(ctx).Model(&Payment{}).Where("biz_tran_no = ?", bizTranNo).Updates(map[string]any{
		"updated_at": now,
		"status":     status,
		"txn_id":     TxnID,
	}).Error
}

func NewGORMPaymentDAO(db *gorm.DB) PaymentDAO {
	return &GORMPaymentDAO{
		db: db,
	}
}
