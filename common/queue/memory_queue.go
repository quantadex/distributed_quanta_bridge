package queue

import (
	dll "github.com/emirpasic/gods/lists/doublylinkedlist"
	"errors"
	"sync"
)

type MemoryQueue struct {
	queues map[string]*dll.List
}


func NewMemoryQueue() Queue {
	q := &MemoryQueue{}
	return q
}

var instance Queue
var once sync.Once

func GetGlobalQueue() Queue {
	once.Do(func() {
		instance = NewMemoryQueue()
	})
	return instance
}

func (q *MemoryQueue) Connect(name string) error {
	return nil
}

func (q *MemoryQueue) CreateQueue(queueName string) error {
	q.queues[queueName] = dll.New()
	return nil
}

func (q *MemoryQueue) Put(queueName string, data []byte) error {
	if q := q.queues[queueName]; q != nil {
		q.Append(data)
	}
	return errors.New("queue not found")
}

func (q *MemoryQueue) Empty(queueName string) error {
	if q := q.queues[queueName]; q != nil {
		q.Clear()
		return nil
	}
	return errors.New("queue not found")
}

func (q *MemoryQueue) Get(queueName string) (data []byte, err error) {
	if q := q.queues[queueName]; q != nil {
		data, found := q.Get(0)
		if found == false {
			return nil, errors.New("no item")
		}
		q.Remove(0)
		return data.([]byte), nil
	}
	return nil, errors.New("queue not found")
}

