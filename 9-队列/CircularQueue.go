package main

import "fmt"

type CircularQueue struct {
	queue []interface{}
	capacity int
	head int
	tail int
}

// 创建
func NewCircularQueue(n int) *CircularQueue {
	if n == 0 {
		return nil
	}

	return &CircularQueue{make([]interface{}, n), n, 0, 0}
}

// 是否为空
func (this *CircularQueue) IsEmpty() bool{
	if this.head == this.tail {
		return true
	}

	return false
}

// 是否已满
func (this *CircularQueue) IsFull() bool {
	if this.head == (this.tail+1) % this.capacity {
		return true
	}

	return false
}

// 进队
func (this *CircularQueue) EnQueue(v interface{}) bool {
	if this.IsFull() {
		return false
	}
	this.queue[this.tail] = v
	this.tail = (this.tail+1) % this.capacity

	return true
}

// 出队
func (this *CircularQueue) DeQueue() interface{} {
	if this.IsEmpty() {
		return nil
	}

	val := this.queue[this.head]
	this.head = (this.tail+1) % this.capacity

	return val
}

// 打印
func (this *CircularQueue) Print() string {
	if this.IsEmpty() {
		return "Empty queue"
	}
	result := "head"
	var i = this.head
	for true {
		result += fmt.Sprintf(" <- %+v", this.queue[i])
		i = (i + 1) % this.capacity
		if i == this.tail {
			break
		}
	}
	result += " <- tail"

	return result
}

func InitCircularQueue() {
	cq := NewCircularQueue(3)
	cq.EnQueue(2)
	fmt.Println("cq.IsFull(): ", cq.IsFull())
	fmt.Println("cq.IsEmpty(): ", cq.IsEmpty())
	fmt.Println(cq.Print())
}

func main() {
	InitCircularQueue()
}