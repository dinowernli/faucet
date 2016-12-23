package worker

import (
	pb_worker "dinowernli.me/faucet/proto/service/worker"
)

// TODO(dino): Find a suitable library and avoid hand-rolling the queue.

type queue struct {
	front *queueNode // The next element to be dequeued.
	back  *queueNode // The last element to be dequeued.
	size  int
}

func newQueue() *queue {
	return &queue{front: nil, back: nil, size: 0}
}

func (q *queue) enqueue(request *pb_worker.ExecutionRequest) {
	if q.empty() {
		q.back = &queueNode{next: nil, prev: nil, request: request}
		q.front = q.back
		q.size = 1
		return
	}

	newBack := &queueNode{next: q.back, prev: nil, request: request}
	q.back.prev = newBack
	q.back = newBack
	q.size++
}

// dequeue removes the next element from the queue and returns it. The queue
// must not be empty.
func (q *queue) dequeue() *pb_worker.ExecutionRequest {
	result := q.front
	newFront := q.front.prev
	if newFront != nil {
		newFront.next = nil
	}
	q.front = newFront
	q.size--
	return result.request
}

func (q *queue) empty() bool {
	return q.size == 0
}

type queueNode struct {
	next    *queueNode
	prev    *queueNode
	request *pb_worker.ExecutionRequest
}
