package queue

/**
 * Queue
 *
 * This is a straight-forward thread safe queue.
 * The only twist is it is actually a collection of named queues
 * each of which is handled separately.
 */
type Queue interface {
    Connect(name string) error
    CreateQueue(queueName string) error
    Put(queueName string, data []byte) error
    Empty(queueName string) error
    Get(queueName string) (data []byte, err error)
}
