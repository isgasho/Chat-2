package models

import (
	"errors"
)

const (
	defaultQueueSize = 10
)

type MyQueue struct {
	front        int
	rear         int
	currentCount int
	queueSize    int
	elements     []interface{}
}

/**
  指定大小的初始化
*/
func NewMyQueueBySize(size int) *MyQueue {
	return &MyQueue{0, size - 1, 0, size, make([]interface{}, size)}
}

/**
  按默认大小进行初始化
*/
func NewMyQueue() *MyQueue {
	return NewMyQueueBySize(defaultQueueSize)
}

/**
  向下一个位置做探测
*/
func (queue *MyQueue) ProbeNext(i int) int {
	return (i + 1) % queue.queueSize
}

func (queue *MyQueue) Probe(i int) int {
	return i % queue.queueSize
}

/**
  清空队列
*/
func (queue *MyQueue) ClearQueue() {
	queue.front = 0
	queue.rear = queue.queueSize - 1
	queue.currentCount = 0
}

/**
  是否为空队列
*/
func (queue *MyQueue) IsEmpty() bool {
	if queue.currentCount == 0 {
		return true
	}
	return false
}

/**
  队列是否满了
*/
func (queue *MyQueue) IsFull() bool {
	if queue.currentCount == queue.queueSize {
		return true
	}
	return false
}

/**
  入队
*/
func (queue *MyQueue) Offer(e interface{}) error {
	if queue.IsFull() == true {
		return errors.New("the queue is full.")
	}
	queue.rear = queue.ProbeNext(queue.rear)
	queue.elements[queue.rear] = e
	queue.currentCount = queue.currentCount + 1
	return nil
}

/**
  出队一个元素
*/
func (queue *MyQueue) Poll() (interface{}, error) {
	if queue.IsEmpty() == true {
		return nil, errors.New("the queue is empty.")
	}
	tmp := queue.front
	queue.front = queue.ProbeNext(queue.front)
	queue.currentCount = queue.currentCount - 1
	return queue.elements[tmp], nil
}

//返回列表
func (queue *MyQueue) GetArray(uidmax int32) (*[]interface{}, error) {
	result := make([]interface{}, 0, queue.queueSize)
	for index := 0; index < queue.currentCount; index++ {
		tmp := queue.elements[queue.Probe(queue.front+index)].(*ChatModel)
		if tmp.UID > uidmax {
			result = append(result, tmp)
		}
	}

	// if queue.front < queue.rear {
	// 	var front = queue.front
	// 	for front = queue.front; front < queue.rear; front++ {
	// 		tmp := queue.elements[front].(*ChatModel)
	// 		if tmp.UID > uidmax {
	// 			break
	// 		}
	// 	}
	// 	if front < queue.rear {
	// 		result = append(result, queue.elements[front:queue.rear])

	// 	}
	// } else if queue.front > queue.rear {
	// 	var front = queue.front
	// 	for front = queue.front; front < queue.queueSize; front++ {
	// 		tmp := queue.elements[front].(*ChatModel)
	// 		if tmp.UID > uidmax {
	// 			break
	// 		}
	// 	}
	// 	if front < queue.queueSize {
	// 		result = append(result, queue.elements[front:], queue.elements[:queue.rear])
	// 	} else {
	// 		front = 0
	// 		for front = 0; front < queue.rear; front++ {
	// 			tmp := queue.elements[front].(*ChatModel)
	// 			if tmp.UID > uidmax {
	// 				break
	// 			}
	// 		}
	// 		if front < queue.rear {
	// 			result = append(result, queue.elements[front:queue.rear])
	// 		}
	// 	}
	// } else {
	// 	//等于就是没有数据
	// }
	return &result, nil
}
