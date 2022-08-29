package app

import (
	"io"
	"sync"
	"txsystem/internal/models"

	"github.com/pkg/errors"
)

type Store interface {
	CreateUser(u *models.User) error
	CreateTx(tx *models.Transaction) error
	UpdateTxStatusByID(status models.Status, id int) error
	GetBalanceByUserID(id int) (int, error)
	AddBalanceByID(id int, amount uint) error
	SubtractBalanceByID(id int, amount uint) error
}

type TransactionService interface {
	io.Closer
	ChangeBalance(tx *models.Transaction) error
	Done(id int)
}

type transactionService struct {
	wg    *sync.WaitGroup
	mu    sync.Mutex
	store Store
	close chan struct{}

	req  chan *models.Transaction
	done chan int
}

func New(store Store) *transactionService {
	t := &transactionService{
		wg:    &sync.WaitGroup{},
		store: store,
		close: make(chan struct{}),
		req:   make(chan *models.Transaction),
		done:  make(chan int),
	}

	t.wg.Add(1)
	go t.run()

	return t
}

func (t *transactionService) run() {
	defer t.wg.Done()

	clients := make(map[int]*instance)

main:
	for {
		select {
		case <-t.close:
			break main

		case req := <-t.req:
			inst, ok := clients[req.UserID]
			if !ok {
				inst = NewInstance(req.UserID, inst.store, inst.wg, inst.done)

				clients[req.UserID] = inst
			}

			inst.GetRecvCh() <- req

		case id := <-t.done:
			t.mu.Lock()
			delete(clients, id)
			t.mu.Unlock()
		}
	}

}

func (t *transactionService) Done(id int) {
	select {
	case <-t.close:
	case t.done <- id:
	}
}

func (t *transactionService) ChangeBalance(tx *models.Transaction) error {
	tx.SetNewStatus()

	if err := t.store.CreateTx(tx); err != nil {
		return errors.Wrap(err, "create tx error")
	}

	if tx.Action == models.SUBTRACT {
		balance, err := t.store.GetBalanceByUserID(tx.UserID)
		if err != nil {
			return errors.Wrap(err, "get balance error")
		}

		if !tx.CheckSubtract(balance) {
			return errors.New("low balance")
		}
	}
	// res := make(chan error, 1)

	// req := &models.TransactionRequest{
	// 	Transaction: tx,
	// 	// Res:         res,
	// }

	t.req <- tx
	// err := <-req.Res
	// if err != nil {
	// 	return errors.Wrap(err, "change balance error")
	// }

	return nil
}

func (t *transactionService) CreateUser(user *models.User) error {
	if err := t.store.CreateUser(user); err != nil {
		return errors.Wrap(err, "create user error")
	}

	return nil
}

func (t *transactionService) Close() error {
	t.close <- struct{}{}
	t.wg.Wait()
	return nil
}