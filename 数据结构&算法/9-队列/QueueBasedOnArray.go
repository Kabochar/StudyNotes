package main

import "fmt"

type ArrayQueue struct {
	queue    []interface{}
	capacity int
	head     int
	tail     int
}

// 创建
func NewArrayQueue(n int) *ArrayQueue {
	return &ArrayQueue{
		make([]interface{}, n),
		n,
		0,
		0,
	}
}

// 进队
func (this *ArrayQueue) EnQueue(v interface{}) bool {
	if this.tail == this.capacity {
		return false
	}

	this.queue[this.tail] = v
	this.tail++

	return true
}

// 出队
func (this *ArrayQueue) DeQueue() interface{} {
	if this.head == this.tail {
		return nil
	}
	val := this.queue[this.head]
	this.head++

	return val
}

// 打印队列
func (this *ArrayQueue) Print() string {
	if this.head == this.tail {
		return "empty queue"
	}
	result := "head"
	for i := this.head; i <= this.tail-1; i++ {
		result += fmt.Sprintf(" <- %+v", this.queue[i])
	}
	result += " <- tail"

	return result
}

func InitArrayQueue() {
	aq := NewArrayQueue(5)
	aq.EnQueue(1)
	aq.EnQueue(2)
	fmt.Println(aq.Print())

	aq.DeQueue()
	fmt.Println(aq.Print())

}

func main() {
	InitArrayQueue()
}
