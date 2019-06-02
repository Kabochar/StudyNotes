package main

import "fmt"

// 通过两个栈来实现
type Browser struct {
	forwardStack Stack
	backStack Stack
}

func NewBrowser() *Browser {
	return &Browser{
		forwardStack: NewArrayStack(),
		backStack: NewLinkedListStack(),
	}
}

// 前驱
func (this *Browser) CanForward() bool {
	if this.forwardStack.IsEmpty() {
		return false
	}

	return true
}
// 后驱
func (this *Browser) CanBack() bool {
	if this.backStack.IsEmpty() {
		return false
	}

	return true
}

// 打开页面
func (this *Browser) Open(addr string) {
	fmt.Printf("Open new addr %+v\n", addr)
	this.forwardStack.Flush()
}

// 返回
func (this *Browser) PushBack(addr string) {
	this.backStack.Push(addr)
}

// 前进
func (this *Browser) Forward() {
	if this.forwardStack.IsEmpty() {
		return
	}

	top := this.forwardStack.Pop()
	this.backStack.Push(top)
	fmt.Println("forward to %+v\n", top)
}

// 返回
func (this *Browser) Back() {
	if this.backStack.IsEmpty() {
		return
	}

	top := this.backStack.Pop()
	this.forwardStack.Push(top)
	fmt.Printf("Back to %+v\n", top)
}