package app

import (
	"sync"
	"txsystem/internal/models"

	"github.com/sirupsen/logrus"
)

type instance struct {
	id        int
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
		// i.mu.Lock()
		// logrus.Info(1111, len(i.queue))
		tx := i.queue[0]
		i.queue = i.queue[1:]
		// time.Sleep(1 * time.Second)
		logrus.Info(tx.Transaction)
		// i.mu.Unlock()
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
