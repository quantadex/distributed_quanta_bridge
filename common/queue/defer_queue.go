package queue

import (
	dll "github.com/emirpasic/gods/lists/doublylinkedlist"
	"sync"
	"github.com/go-errors/errors"
)
type deferQ map[int]*dll.List

type DeferQ struct {
	sync.Mutex
	data map[string]deferQ
	delay int
	curEpoch int
}

func NewDeferQ(delay int) *DeferQ {
	return &DeferQ{
		data: make(map[string]deferQ),
		delay: delay,
		curEpoch: 0,
	}
}

func (q *DeferQ) Connect(name string) error {
	return nil
}

func (q *DeferQ) CreateQueue(queueName string) error {
	if _, ok := q.data[queueName]; ok {
		return nil
	}
	q.data[queueName] = make(map[int]*dll.List, 0)
	return nil
}

func (q *DeferQ) Put(queueName string, msg interface{}) error {
	queueData, found := q.data[queueName]
	if !found {
		return errors.New("No queue found")
	}

	expires := q.curEpoch + q.delay
	//r.log.Info(fmt.Sprintf("Added msg to DeferQ at expiration=%d epoch=%d", expires, q.curEpoch))

	var deferList *dll.List

	deferList, found = queueData[expires]
	if !found {
		deferList = dll.New()
		queueData[expires] = deferList
	}
	deferList.Append(msg)
	return nil
}

func (q *DeferQ) Empty(queueName string) error {
	return nil
}

func (q *DeferQ) Get(queueName string) (data interface{}, err error) {
	queueData, found := q.data[queueName]
	if !found {
		return nil, errors.New("No queue found")
	}

	for k, v := range queueData {
		if k <= q.curEpoch {
			it := v.Iterator()
			for it.Next() {
				result := it.Value()
				delete(queueData, k)
				return result, nil
			}
		}
	}
	return nil, nil
}

func (q *DeferQ) AddTick() {
	q.curEpoch +=1
}

func (q *DeferQ) Epoch() int {
	return q.curEpoch
}