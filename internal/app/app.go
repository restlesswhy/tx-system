package app

import (
	"io"
	"sync"
	"txsystem/internal/models"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
	close chan struct{}

	store Store
	req   chan *models.Transaction
	done  chan int
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

	clients := make(map[int]chan<- *models.Transaction)

main:
	for {
		select {
		case <-t.close:
			break main

		case req := <-t.req:
			recv, ok := clients[req.UserID]
			if !ok {
				ch := make(chan *models.Transaction)
				newInstance(req.UserID, t.store, t.wg, t.done, ch)

				clients[req.UserID] = ch

				ch <- req

				continue
			}

			recv <- req

		case id := <-t.done:
			logrus.Infof("del %d inst", id)
			close(clients[id])
			delete(clients, id)
		}
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

	t.req <- tx

	return nil
}

func (t *transactionService) CreateUser(user *models.User) error {
	if err := t.store.CreateUser(user); err != nil {
		return errors.Wrap(err, "create user error")
	}

	return nil
}

func (t *transactionService) Close() error {
	close(t.close)
	t.wg.Wait()
	return nil
}
