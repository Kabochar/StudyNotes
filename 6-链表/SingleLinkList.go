package main

import "fmt"

type ListNode struct {
	value interface{}
	next *ListNode
}

type LinkedList struct {
	head *ListNode
	length uint
}

func (this *LinkedList) GetLength() uint {
	return this.length
}

func (this *ListNode) GetValue() interface{} {
	return this.value
}

func (this *ListNode) GetNext() *ListNode {
	return this.next
}

func NewListNode(v interface{}) *ListNode {
	return &ListNode{v, nil}
}

func NewLinkedList() *LinkedList {
	return &LinkedList{NewListNode(0),0}
}

func (this *LinkedList) InsertAfter(p *ListNode, v interface{}) bool {
	if nil == p {
		return false
	}
	newNode := NewListNode(v)
	oldNode := p.next
	p.next = newNode
	newNode.next = oldNode
	this.length++

	return true
}

func (this *LinkedList) InsertBefore(p *ListNode, v interface{}) bool   {
	if nil == p || this.head == p {
		return false
	}
	cur := this.head.next
	pre := this.head
	for nil != cur {
		if cur == p {
			break
		}
		cur = cur.next // 寻找 p 节点位置
		pre = cur
	}
	if nil == cur { // 经过循环，cur 找到了 tail.next，错误
		return false
	}
	// 经过 判断，cur 此时为能 insert 的 节点
	newNode := NewListNode(v)
	pre.next = newNode
	newNode.next = cur
	this.length++

	return false
}

func (this *LinkedList) InsertToHead(v interface{}) bool {
	return this.InsertAfter(this.head, v) // LinkedList 存在哨兵
}

func (this *LinkedList) InsertToTail(v interface{}) bool {
	cur := this.head
	for nil != cur.next {
		cur = cur.next
	}

	return this.InsertAfter(cur, v)
}

func (this *LinkedList) FindByIndex(index uint) *ListNode {
	if index >= this.length {
		return nil
	}
	cur := this.head.next
	var i uint = 0
	for ; i < index; i++ {
		cur = cur.next
	}

	return cur
}

func (this *LinkedList) DeleteNode(p *ListNode) bool {
	if nil == p {
		return false
	}
	cur := this.head.next
	pre := this.head
	for nil != cur {
		if cur == p {
			break
		}
		pre = cur
		cur = cur.next
	}
	if nil == cur {
		return false
	}
	pre.next = p.next
	p = nil
	this.length--

	return true
}

func (this *LinkedList) Print() {
	cur := this.head.next
	format := ""
	for nil != cur {
		format += fmt.Sprintf("%+v", cur.value)
		cur = cur.next
		if nil != cur {
			format += " -> "
		}
	}
	fmt.Println(format)
}

// 反转
func (this *LinkedList) Reverse() bool {
	// 最起码保证有两个节点
	if nil == this.head || nil == this.head.next || nil == this.head.next.next {
		return false
	}
	var pre *ListNode = nil
	cur := this.head.next
	for nil != cur {
		temp := cur.next
		cur.next = pre
		pre = cur
		cur = temp
	}

	this.head.next = pre

	return true
}

// 是否有环
func (this *LinkedList) HasCircle() bool {
	if nil != this.head {
		slow := this.head
		fast := this.head

		for nil != fast && nil != fast.next {
			slow = slow.next
			fast = fast.next.next
			if slow == fast {
				return true
			}
		}
	}

	return false
}

// 寻找中间节点
func(this *LinkedList) FindMiddleNode() *ListNode {
	// 只有哨兵节点
	if nil == this.head || nil == this.head.next{
		return nil
	}
	// 如果只有 2 个节点
	if nil == this.head.next.next {
		return this.head.next
	}

	slow, fast := this.head, this.head
	for nil != fast && nil != fast.next {
		slow = slow.next
		fast = fast.next.next
	}

	return slow
}

// 删除倒数第 N 个节点
func (this *LinkedList) DeleteBottomN(n int) bool {
	if n < 0 || nil == this.head || nil == this.head.next {
		return false
	}

	fast := this.head
	for i := 1; i < n && nil != fast; i++ {
		fast = fast.next
	}

	if nil == fast {
		return false
	}

	slow := this.head
	for nil != fast.next {
		slow = slow.next
		fast = fast.next
	}

	slow.next = slow.next.next

	return true
}

func InitLinkedList() {
	list := NewLinkedList()

	list.InsertToTail(5)
	list.InsertToTail(1)
	list.Print()


	list.Reverse()
	list.Print()

	ok := list.DeleteBottomN(3)
	fmt.Println(ok)
}

func main() {
	InitLinkedList()
}