package main

import "fmt"

type ListNode struct {
	val  interface{}
	next *ListNode
}

type LinkedListQueue struct {
	head   *ListNode
	tail   *ListNode
	length int
}

func NewLinkedListQueue() *LinkedListQueue {
	return &LinkedListQueue{
		nil,
		nil,
		0}
}

// 进队
func (this *LinkedListQueue) EnQueue(v interface{}) {
	node := &ListNode{v, nil}
	if nil == this.tail {
		this.tail = node
		this.head = node
	} else {
		this.tail.next = node
		this.tail = node
	}
	this.length++
}

// 出队
func (this *LinkedListQueue) DeQueue() interface{} {
	if this.head == nil {
		return nil
	}
	val := this.head.val
	this.head = this.head.next

	return val
}

// 打印
func (this *LinkedListQueue) Print() string {
	if nil == this.head {
		return "empty Queue"
	}

	result := "head "
	for cur := this.head; nil != cur; cur = cur.next {
		result += fmt.Sprintf("<- %+v ", cur.val)
	}
	result += "<- tail"

	return result
}

func InitLinkedListQueue() {
	lq := NewLinkedListQueue()
	lq.EnQueue(4)
	lq.EnQueue(5)
	fmt.Println(lq.Print())

	lq.DeQueue()
	fmt.Println(lq.Print())
}

func main() {
	InitLinkedListQueue()
}
