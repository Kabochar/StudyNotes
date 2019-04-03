package main

import (
	"fmt"
)

type BinaryTree struct {
	root *Node
}

func NewBinaryTree(rootV interface{}) *BinaryTree {
	return &BinaryTree{NewNode(rootV)}
}

func (this *BinaryTree) InOrderTraverse() {
	pRoot := this.root
	stack := NewArrayStack()

	for !stack.IsEmpty() || nil != pRoot {
		if nil != pRoot {
			stack.Push(pRoot)
			pRoot = pRoot.left
		} else {
			temp := stack.Pop().(*Node)
			fmt.Printf("%+v ", temp.data)
			pRoot = temp.right
		}
	}
	fmt.Println()
}

func (this *BinaryTree) PreOrderTraverse() {
	pRoot := this.root
	stack := NewArrayStack()

	for !stack.IsEmpty() || nil != pRoot {
		if nil != pRoot {
			fmt.Printf("%+v ", pRoot.data)
			stack.Push(pRoot)
			pRoot = pRoot.left
		} else {
			pRoot = stack.Pop().(*Node).right
		}
	}
	fmt.Println()
}


func (this *BinaryTree) PostOrderTraverse() {
	stack_1 := NewArrayStack()
	stack_2 := NewArrayStack()

	stack_1.Push(this.root)

	for !stack_1.IsEmpty() {
		p := stack_1.Pop().(*Node)
		stack_2.Push(p)
		if nil != p.left {
			stack_1.Push(p.left)
		}
		if nil != p.right {
			stack_1.Push(p.right)
		}
	}

	for !stack_2.IsEmpty() {
		fmt.Printf("%+v", stack_2.Pop().(*Node).data)
	}
}

// 使用一个堆栈，前置光标从后期顺序遍历
func (this *BinaryTree) PostOrderTraverse2() {
	root := this.root
	stack := NewArrayStack()

	// 指向上次访问节点
	var pre *Node

	stack.Push(root)

	for !stack.IsEmpty() {
		root = stack.Top().(*Node)
		if (root.left == nil && root.right == nil) ||
			(pre != nil && (pre == root.left || pre == root.right)) {
			fmt.Printf("%+v ", root.data)
			stack.Pop()
			pre = root
		} else {
			if root.right != nil {
				stack.Push(root.right)
			}

			if root.left != nil {
				stack.Push(root.left)
			}
		}
	}
}
