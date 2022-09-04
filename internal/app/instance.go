package app

import (
	"fmt"
	"sync"
	"txsystem/internal/models"

	"github.com/sirupsen/logrus"
)

type instance struct {
	id    int
	store Store
	wg    *sync.WaitGroup
	queue []*models.Transaction
	done  chan<- int

	recv chan *models.Transaction

	close chan struct{}
}

func newInstance(id int, store Store, wg *sync.WaitGroup, done chan<- int, recv chan *models.Transaction) *instance {
	i := &instance{
		id:    id,
		store: store,
		wg:    wg,
		done:  done,
		queue: make([]*models.Transaction, 0),
		recv:  recv,
		close: make(chan struct{}),
	}

	wg.Add(2)
	go i.recieve()
	go i.do()

	return i
}

func (i *instance) do() {
	defer i.wg.Done()

	for {
		if len(i.queue) > 0 {
			fmt.Println(len(i.queue))
			logrus.Info("hello")
			// t.Reset(3 * time.Second)
			tx := i.queue[0]
			i.queue = i.queue[1:]

			switch tx.Action {
			case models.ADD:
				if err := i.store.AddBalanceByID(tx.UserID, tx.Amount); err != nil {
					if err := i.store.UpdateTxStatusByID(models.FAIL_TX, tx.ID); err != nil {
						logrus.Errorf("change tx status error: %v", err)
					}

					logrus.Errorf("add balance error: %v", err)
					continue
				}

				if err := i.store.UpdateTxStatusByID(models.DONE_TX, tx.ID); err != nil {
					logrus.Errorf("change tx status error: %v", err)
				}
				logrus.Info("added")

			case models.SUBTRACT:
				if err := i.store.SubtractBalanceByID(tx.UserID, tx.Amount); err != nil {
					if err := i.store.UpdateTxStatusByID(models.FAIL_TX, tx.ID); err != nil {
						logrus.Errorf("change tx status error: %v", err)
					}

					logrus.Errorf("subtract balance error: %v", err)
					continue
				}

				if err := i.store.UpdateTxStatusByID(models.DONE_TX, tx.ID); err != nil {
					logrus.Errorf("change tx status error: %v", err)
				}
			}
		} else {
			i.done <- i.id
			close(i.close)
			break
		}
	}
}

func (i *instance) recieve() {
	logrus.Info("start reciever")
	defer i.wg.Done()

main:
	for {
		select {
		case <-i.close:
			break main

		case tx, ok := <-i.recv:
			if !ok {
				continue
			}
			i.queue = append(i.queue, tx)

		}
	}

	logrus.Debug("reciever closed")
}
