package app

import (
	"sync"
	"time"
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

func NewInstance(id int, store Store, wg *sync.WaitGroup, done chan<- int) *instance {
	i := &instance{
		id:    id,
		store: store,
		wg:    wg,
		done:  done,
		queue: make([]*models.Transaction, 0),
		recv:  make(chan *models.Transaction),
		close: make(chan struct{}),
	}

	wg.Add(1)
	go i.recieve()

	return i
}

func (i *instance) recieve() {
	logrus.Info("start reciever")
	defer i.wg.Done()

	t := time.NewTimer(3 * time.Second)

main:
	for {
		select {
		case <-i.close:
			break main

		case <-t.C:
			logrus.Info("stoping recieve")
			// TODO: check if done closed
			i.done <- i.id
			break main

		case tx, ok := <-i.recv:
			if !ok {
				continue
			}
			i.queue = append(i.queue, tx)

		default:
			if len(i.queue) > 0 {
				logrus.Info("hello")
				t.Reset(3 * time.Second)
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
			}
		}
	}

	close(i.recv)
	logrus.Info("reciever closed")
}

func (i *instance) GetRecvCh() chan<- *models.Transaction {
	return i.recv
}

func (i *instance) Close() {
	close(i.close)
}
