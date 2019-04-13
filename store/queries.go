package store

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"datastruct/queue"
)

const (
	clearInterval = 30 * time.Second
	queryLifetime = 5 * time.Minute
)

type QueriesStopChan = chan interface{}

type queryData struct {
	filters			interface{}
	id				string
	lastActivity	time.Time
}

type Queries struct {
	dict				map[string]queue.Node
	mutex				sync.RWMutex
	queue				*queue.Queue
	stopChan			QueriesStopChan
	stoppedChan			QueriesStopChan
}

func worker(q *Queries) {
	ticker := time.NewTicker(clearInterval)

	defer func() {
		ticker.Stop()
		close(q.stoppedChan)
	}()


	for {
		select {
		case <-q.stopChan:
			return
		case <-ticker.C:
			q.clear()
			break
		}
	}
}

func NewQueries() *Queries {
	q := &Queries{
		dict: make(map[string]queue.Node),
		queue: queue.New(),
		stopChan: make(QueriesStopChan),
		stoppedChan: make(QueriesStopChan),
	}

	go worker(q)

	return q
}

func (q *Queries) Add(filters interface{}) (string, error) {
	uid, err := uuid.NewRandom()

	if err != nil {
		return "", err
	}

	suid := uid.String()

	data := &queryData{
		filters: filters,
		id: suid,
		lastActivity: time.Now(),
	}

	q.mutex.Lock()
	defer q.mutex.Unlock()

	node := q.queue.Push(data)
	q.dict[suid] = node

	return suid, nil
}

func (q *Queries) Get(id string) (interface{}, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	item, ok := q.dict[id]

	if !ok {
		return nil, errors.New("query not find")
	}

	err := q.queue.MoveToEnd(item)
	if err != nil {
		return nil, err
	}

	rawData := item.GetData()
	data := rawData.(*queryData)
	data.lastActivity = time.Now()

	return data.filters, nil
}

func (q *Queries) clear() {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for {
		rawData := q.queue.First()

		if rawData == nil {
			return
		}

		data := rawData.(*queryData)

		diff := time.Since(data.lastActivity)
		if  diff < queryLifetime {
			return
		}

		q.queue.Pop()
	}
}

func (q *Queries) Stop() {
	close(q.stopChan)
	<-q.stoppedChan
}
