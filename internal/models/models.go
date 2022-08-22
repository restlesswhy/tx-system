package models

import "time"

type Action string
type Status int

const (
	ADD      Action = "add"
	SUBTRACT Action = "subtract"

	NEW_TX Status = iota + 1
	DONE_TX
	FAIL_TX
)

type TransactionRequest struct {
	Transaction *Transaction
	Res         chan error
}

type Transaction struct {
	ID       int       `json:"-"`
	UserID   int       `json:"user_id,omitempty"`
	Amount   int       `json:"amount,omitempty"`
	Action   Action    `json:"action,omitempty"`
	CreateAt time.Time `json:"create_at,omitempty"`
	Status   Status    `json:"status,omitempty"`
}

func (t *Transaction) SetNewStatus() {
	t.Status = NEW_TX
}

func (t Transaction) Validate() bool {
	if t.UserID <= 0 {
		return false
	}

	if t.Amount <= 0 {
		return false
	}

	return true
}

type User struct {
	ID      int
	Balance uint
}
