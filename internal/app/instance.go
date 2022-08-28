package app

import (
	"sync"
	"txsystem/internal/models"

	"github.com/sirupsen/logrus"
)

type instance struct {
	id        int
	store     Store
	wg        *sync.WaitGroup
	queue     []*models.TransactionRequest
	txService TransactionService
	mu        sync.Mutex

	alive bool
}

func (i *instance) run() {
	i.alive = true
	defer func() {
		i.mu.Lock()
		i.txService.Done(i.id)
		i.alive = false
		i.wg.Done()
		i.mu.Unlock()
	}()

	for len(i.queue) > 0 {
		tx := i.queue[0]
		i.queue = i.queue[1:]
		switch tx.Transaction.Action {
		case models.ADD:
			if err := i.store.AddBalanceByID(tx.Transaction.UserID, tx.Transaction.Amount); err != nil {
				if err := i.store.UpdateTxStatusByID(models.FAIL_TX, tx.Transaction.ID); err != nil {
					logrus.Errorf("change tx status error: %v", err)
				}

				logrus.Errorf("add balance error: %v", err)
				continue
			}

			if err := i.store.UpdateTxStatusByID(models.DONE_TX, tx.Transaction.ID); err != nil {
				logrus.Errorf("change tx status error: %v", err)
			}

		case models.SUBTRACT:
			if err := i.store.SubtractBalanceByID(tx.Transaction.UserID, tx.Transaction.Amount); err != nil {
				if err := i.store.UpdateTxStatusByID(models.FAIL_TX, tx.Transaction.ID); err != nil {
					logrus.Errorf("change tx status error: %v", err)
				}

				logrus.Errorf("subtract balance error: %v", err)
				continue
			}

			if err := i.store.UpdateTxStatusByID(models.DONE_TX, tx.Transaction.ID); err != nil {
				logrus.Errorf("change tx status error: %v", err)
			}
		}
	}
}

func (i *instance) AddTx(tx *models.TransactionRequest) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if !i.alive {
		i.wg.Add(1)
		go i.run()
	}
	i.queue = append(i.queue, tx)

	return nil
}
