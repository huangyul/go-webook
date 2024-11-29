package sms

import "context"

type SMSService interface {
	Send(ctx context.Context, tplId string, args []string, number ...string) error
}
