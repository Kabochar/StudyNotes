package main

import "fmt"

// 栈节点
type Node struct {
	next *Node
	val interface{}
}

// 栈链表，链式栈
type LinkedListStack struct {
	topNode *Node
}

// 创建
func NewLinkedListStack() *LinkedListStack {
	return &LinkedListStack{nil}
}

// 空栈
func (this *LinkedListStack) IsEmpty() bool {
	if this.topNode == nil {
		return true
	}

	return false
}

// 压栈
func (this *LinkedListStack) Push(v interface{}) {
	this.topNode = &Node{
		next: this.topNode,
		val: v,
	}
}

// 退栈
func (this *LinkedListStack) Pop() interface{} {
	if this.IsEmpty() {
		return nil
	}

	val := this.topNode.val
	this.topNode = this.topNode.next

	return val
}

// 获取 栈顶元素
func (this *LinkedListStack) Top() interface{} {
	if this.IsEmpty() {
		return nil
	}

	return this.topNode.val
}

// 清空栈
func (this *LinkedListStack) Flush() {
	this.topNode = nil
}

func (this *LinkedListStack) Print() {
	if this.IsEmpty() {
		fmt.Println("Empty stack")
	} else {
		cur := this.topNode
		for nil != cur {
			fmt.Println(cur.val)
			cur = cur.next
		}
	}
}

func InitLinkedListStack() {
	ls := NewLinkedListStack()
	ls.Push(2)
	ls.Push(3)
	ls.Push(5)
	ls.Print()

	fmt.Println("-----")

	val := ls.Top()
	fmt.Println(val)

	fmt.Println("-----")

	ls.Pop()
	ls.Print()

	ls.Flush()
	ls.Print()

}

func main() {
	InitLinkedListStack()
}
