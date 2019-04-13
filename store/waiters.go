package store

import (
	"sync"
	"time"

	"github.com/pkg/errors"

	"common"
	"datastruct/queue"
)

const (
	waiterLifeTime		= time.Minute
	waiterClearInterval	= time.Minute
)

var (
	errorWaiterAddIdAlreadyExist	= errors.New("waiter for id already exist")
	errorWaiterAlreadyStarted		= errors.New("waiter has already started")
	errorWaiterNotStarted			= errors.New("waiter hasn't started yet")
	errorWaiterTryCloseIdNotExist	= errors.New("waiter for id not found")
)

type Waiters interface {
	Add(id string, ch chan *LedgerTransaction) error
	Start() error
	Stop() error
	TryClose(tran *LedgerTransaction) error
	TryCloseAnother(id string, tran *LedgerTransaction) error
	TryInterrupt(id string) error
}

type waiterDictItem struct {
	ch		chan *LedgerTransaction
	node	queue.Node
}

type waiterQueueItem struct {
	id			string
	timestamp	time.Time
}

type waiters struct {
	closed		chan bool
	closing		chan bool
	isStarted	bool

	dict		map[string]*waiterDictItem
	queue		*queue.Queue
	mu			sync.RWMutex
}

func NewWaiters() Waiters {
	return &waiters{
		closed: nil,
		closing: nil,
		isStarted: false,

		dict: nil,
		queue: nil,
	}
}

func (w *waiters) Add(id string, ch chan *LedgerTransaction) error {
	if ok, err := w.exist(id); err != nil {
		return err
	} else if ok {
		return errorWaiterAddIdAlreadyExist
	}

	return w.add(id, ch)
}

func (w *waiters) add(id string, ch chan *LedgerTransaction) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if err := w.checkStartedSync(); err != nil {
		return err
	}
	if _, ok := w.dict[id]; ok {
		return errorWaiterAddIdAlreadyExist
	}

	node := w.queue.Push(&waiterQueueItem{
		id: id,
		timestamp: time.Now(),
	})
	w.dict[id] = &waiterDictItem{
		ch: ch,
		node: node,
	}
	return nil
}

func (w *waiters) Start() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.isStarted {
		return errorWaiterAlreadyStarted
	}

	w.dict = make(map[string]*waiterDictItem)
	w.queue = queue.New()
	w.closing = make(chan bool)
	w.closed = make(chan bool)
	waiterRoutine(w)
	w.isStarted = true
	return nil
}

func (w *waiters) Stop() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.isStarted {
		return errorWaiterNotStarted
	}

	close(w.closing)
	<-w.closed

	for _, item := range w.dict {
		item.ch <- nil
	}
	w.dict = nil
	w.queue = nil
	w.isStarted = false
	return nil
}

func (w *waiters) TryClose(tran *LedgerTransaction) error {
	return w.TryCloseAnother(tran.Id, tran)
}

func (w *waiters) TryCloseAnother(id string, tran *LedgerTransaction) error {
	if ok, err := w.exist(id); err != nil {
		common.Log.PrintFull(
			common.Printf(
				"can't close tran err=%v: id=%v",
				err.Error(),
				id,
			),
		)
		return err
	} else if !ok {
		common.Log.PrintFull(
			common.Printf(
				"can't close tran - not find: id=%v",
				id,
			),
		)
		return errorWaiterTryCloseIdNotExist
	}

	return w.tryClose(id, tran)
}

func (w *waiters) TryInterrupt(id string) error {
	if ok, err := w.exist(id); err != nil {
		return err
	} else if !ok {
		return errorWaiterTryCloseIdNotExist
	}

	return w.tryInterrupt(id)
}

func (w *waiters) tryClose(id string, tran *LedgerTransaction) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if err := w.checkStartedSync(); err != nil {
		return err
	}

	if _, ok := w.dict[id]; !ok {
		return errorWaiterTryCloseIdNotExist
	}

	w.removeSync(id, tran)
	return nil
}

func (w *waiters) tryInterrupt(id string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if err := w.checkStartedSync(); err != nil {
		return err
	}

	if _, ok := w.dict[id]; !ok {
		return errorWaiterTryCloseIdNotExist
	}

	w.removeSync(id, nil)
	return nil
}

func (w *waiters) clear() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if err := w.checkStartedSync(); err != nil {
		return err
	}

	node := w.queue.FirstNode()
	lockTimeout := time.Now().Add(-waiterLifeTime)
	for {
		if node == nil {
			break
		}
		next := node.Next()
		data := node.GetData().(*waiterQueueItem)
		if data.timestamp.Before(lockTimeout) {
			w.removeSync(data.id, nil)
		} else {
			break
		}

		node = next
	}
	return nil
}

func (w *waiters) exist(id string) (bool, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if err := w.checkStartedSync(); err != nil {
		return false, err
	}

	_, ok := w.dict[id]
	return ok, nil
}

func (w *waiters) removeSync(id string, res *LedgerTransaction) {
	item := w.dict[id]
	delete(w.dict, id)
	w.queue.Remove(item.node)
	if res != nil && res.Status != InvalidTransactionStatus {
		item.ch <- res
	} else {
		item.ch <- nil
	}
}

func (w *waiters) checkStartedSync() error {
	if w.isStarted {
		return nil
	} else {
		return errorWaiterNotStarted
	}
}

func waiterRoutine(w *waiters) {
	go func() {
		t := time.NewTicker(waiterClearInterval)
		defer func() {
			t.Stop()
			close(w.closed)
		}()
		for {
			select {
				case <-w.closing:
					return
				case <-t.C:
					w.clear()
			}
		}
	}()
}
