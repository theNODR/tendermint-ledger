package blockchain

import (
	"sync"
	"time"

	"github.com/pkg/errors"
)

const (
	writerAccumulatorTimeout	=	250 * time.Millisecond
)

var (
	errorWriterAlreadyStarted	= errors.New("writer already started")
	errorWriterNotStarted		= errors.New("writer not started")
)

type Writer interface {
	Start() error
	Write(tran *Tran) error
	Stop() error
}

type writer struct {
	closed		chan bool
	closing		chan bool
	isStarted	bool
	mu			sync.Mutex

	trans		[]*Tran
	client		*Client
}

func NewWriter(client *Client) Writer {
	return &writer{
		closed: nil,
		closing: nil,
		client: client,
		isStarted: false,
		trans: nil,
	}
}

func (w *writer) Start() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.isStarted {
		return errorWriterAlreadyStarted
	}

	w.closing = make(chan bool)
	w.closed = make(chan bool)
	w.trans = make([]*Tran, 0)
	writerRoutine(w)
	w.isStarted = true

	return nil
}

func (w *writer) Write(tran *Tran) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.isStarted {
		return errorWriterNotStarted
	}

	w.trans = append(w.trans, tran)
	return nil
}

func (w *writer) write() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.writeSync()
}

func (w *writer) writeSync() {
	if len(w.trans) == 0 {
		return
	}

	w.client.SendTrans(w.trans)
	w.trans = make([]*Tran, 0, 0)
}


func (w *writer) Stop() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.isStarted {
		return errorWriterNotStarted
	}

	close(w.closing)
	<-w.closed

	w.writeSync()
	w.isStarted = true
	return nil
}

func writerRoutine(w *writer) {
	go func() {
		ticker := time.NewTicker(writerAccumulatorTimeout)
		defer func() {
			ticker.Stop()
			close(w.closed)
		}()

		for {
			select {
			case <-ticker.C:
				w.write()
			case <-w.closing:
				return
			}
		}
	}()
}
