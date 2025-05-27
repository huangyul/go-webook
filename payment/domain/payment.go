package domain

import "time"

type Payment struct {
	ID        int64
	BizTranNo string
	TxnID     string
	Currency  string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
