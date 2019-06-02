package main

import "fmt"

type ArrayStack struct {
	data []interface{} // data
	top int // stack top point
}

func NewArrayStack() *ArrayStack {
	return &ArrayStack{
		data: make([]interface{}, 0, 32),
		top: -1,
	}
}

func (this *ArrayStack) IsEmpty() bool {
	if this.top < 0 {
		return true
	}

	return false
}

func (this *ArrayStack) Push(val interface{}) {
	if this.top < 0 {
		this.top = 0
	} else {
		this.top += 1
	}

	if this.top > len(this.data)-1 {
		this.data = append(this.data, val)
	} else {
		this.data[this.top] = val
	}
}

func (this *ArrayStack) Pop() interface{} {
	if this.IsEmpty(){
		return nil
	}
	val := this.data[this.top]
	this.top -= 1

	return val
}

func (this *ArrayStack) Top() interface{} {
	if this.IsEmpty() {
		return nil
	}

	return this.data[this.top]
}

func (this *ArrayStack) Flush() {
	this.top = -1
}

func (this *ArrayStack) Print() {
	if this.IsEmpty() {
		fmt.Println("Empty stack")
	} else {
		for i := this.top; i > 0 ; i-- {
			fmt.Println(this.data[i])
		}
	}
}
