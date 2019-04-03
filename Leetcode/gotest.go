package main

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

func NewListNode(v int) *ListNode {
	return &ListNode{v, nil}
}

func mergeTwoLists(l1 *ListNode, l2 *ListNode) *ListNode {
	if nil == l1 {
		return l2
	} else if nil == l2 {
		return l1
	} else if nil == l1 && nil == l2 {
		return nil
	}

	var head, cur *ListNode
	if l1.Val < l2.Val {
		head = l1
		cur = l1
		l1 = l1.Next
	} else {
		head = l2
		cur = l2
		l2 = l2.Next
	}

	for nil != l1 && nil != l2 {
		if l1.Val < l2.Val {
			cur.Next = l1
			l1 = l1.Next
		} else {
			cur.Next = l2
			l2 = l2.Next
		}
		cur = cur.Next
	}

	if nil != l1 {
		cur.Next = l1
	}

	if nil != l2 {
		cur.Next = l2
	}

	return head
}

func removeNthFromEnd(head *ListNode, n int) *ListNode {
	d, headIsNthFromEnd := getDaddy(head, n)

	if headIsNthFromEnd {
		// 删除head节点
		return head.Next
	}

	d.Next = d.Next.Next

	return head
}

// 获取倒数第N个节点的父节点
func getDaddy(head *ListNode, n int) (daddy *ListNode, headIsNthFromEnd bool) {
	daddy = head

	for head != nil {
		if n < 0 {
			daddy = daddy.Next
		}

		n--
		head = head.Next
	}

	// n==0，说明链的长度等于n
	headIsNthFromEnd = n == 0

	return
}

func main() {
	listf := NewListNode(1)
	listf.Next = NewListNode(2)
	listf.Next.Next = NewListNode(3)

	rst := removeNthFromEnd(listf, 4)
	fmt.Println(rst)
	fmt.Println(rst.Next)
	fmt.Println(rst.Next.Next)
}
