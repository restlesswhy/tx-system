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

	mu   sync.Mutex
	recv chan *models.Transaction

	close chan struct{}
	stop  chan struct{}
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
		stop:  make(chan struct{}),
	}

	wg.Add(2)
	go i.reciever()
	go i.worker()

	return i
}

func (i *instance) worker() {
	logrus.Info("start worker")
	// defer logrus.Info("worker closed")
	defer i.wg.Done()

main:
	for {
		select {
		case <-i.stop:
			break main

		case <-i.close:
			break main

		default:
			// fmt.Println(2)
			if len(i.queue) > 0 {
				logrus.Info("got some in queue")
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

	logrus.Info("worker closed")
}

func (i *instance) reciever() {
	logrus.Info("start reciever")
	defer logrus.Info("reciever closed")
	defer i.wg.Done()

	// t := time.NewTimer(3 * time.Second)
	// main2:
	for {
		select {
		// TODO: safeclose?
		case <-i.close:
			logrus.Info("close ch recv")
			return

		case <-time.After(3 * time.Second):
			// fmt.Println(tt)
			logrus.Info("closing instance....")
			i.done <- i.id
			// TODO: check if worker do job
			close(i.stop)
			return

		case tx, ok := <-i.recv:
			// logrus.Info("start tx")
			// fmt.Println("start tx")
			if !ok {
				continue
			}

			i.mu.Lock()
			i.queue = append(i.queue, tx)
			i.mu.Unlock()
			// fmt.Println("end tx")
			// logrus.Info("end tx")
			// t.Stop()
			// default:
		}
	}

	// logrus.Info("reciever closed")
}

func (i *instance) GetRecvCh() chan<- *models.Transaction {
	return i.recv
}

func (i *instance) Close() {
	close(i.close)
}
