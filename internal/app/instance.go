package app

import (
	"sync"
	"txsystem/internal/models"
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
			
		case models.SUBTRACT:

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
