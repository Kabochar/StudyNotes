package main

// Node 队列节点
type Node struct {
	value    int
	priority int
}

// PQueue priority queue
type PQueue struct {
	heap []Node

	capacity int
	used     int
}

// NewPriorityQueue new
func NewPrintorQueue(capacity int) PQueue {
	return PQueue{
		heap:     make([]Node, capacity+1, capacity+1),
		capacity: capacity,
		used:     0,
	}
}

// Push 入队
func (this *PQueue) Push(node Node) {
	if this.used > this.capacity {
		return
	}
	this.used++
	this.heap[this.used] = node
}

// Pop 出队列
func (this *PQueue) Pop() Node {
	if this.used == 0 {
		return Node{-1, 1}
	}

	// 先堆化, 再取堆顶元素
	adjustHeap(this.heap, 1, this.used)
	node := this.heap[1]

	this.heap[1] = this.heap[this.used]
	this.used--

	return node
}

// Top 获取队列顶部元素
func (this *PQueue) Top() Node {
	if this.used == 0 {
		return Node{-1, -1}
	}

	adjustHeap(this.heap, 1, this.used)

	return this.heap[1]
}
