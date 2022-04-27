// Thread-safe queue implementation.
package DataStructures

import "sync"

type ConcurrentQueue struct {
	head *node
	tail *node
	size int
	*sync.Mutex
}

type node struct {
	data interface{}
	next *node
	prev *node
}

func (q ConcurrentQueue) Size() int {
	return q.size
}

func (q ConcurrentQueue) IsEmpty() bool {
	return q.size == 0
}

func (q *ConcurrentQueue) Enqueue(data interface{}) {
	q.Lock()
	defer q.Unlock()

	newNode := &node{data: data}

	if q.IsEmpty() {
		q.head = newNode
		q.tail = newNode
	} else {
		q.tail.next = newNode
		newNode.prev = q.tail
		q.tail = newNode
	}

	q.size++
}

func (q *ConcurrentQueue) Dequeue() interface{} {
	q.Lock()
	defer q.Unlock()

	if q.IsEmpty() {
		return nil
	}

	data := q.head.data

	if q.head == q.tail {
		q.head = nil
		q.tail = nil
	} else {
		q.head = q.head.next
		q.head.prev = nil
	}

	q.size--

	return data
}

func (q *ConcurrentQueue) Peek() interface{} {
	q.Lock()
	defer q.Unlock()

	if q.IsEmpty() {
		return nil
	}

	return q.head.data
}

func (q *ConcurrentQueue) Clear() {
	q.Lock()
	defer q.Unlock()

	q.head = nil
	q.tail = nil
	q.size = 0
}
