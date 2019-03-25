package main

import (
	"fmt"
)

type ArrayStack struct {
	data []interface{}
	top int
}

func NewArrayStack() *ArrayStack {
	return &ArrayStack{
		data: make([]interface{}, 0, 20),
		top: -1,
	}
}

func (this *ArrayStack) IsEmpty() bool {
	if this.top < 0 {
		return true
	}

	return false
}

func (this *ArrayStack) Push(v interface{}) {
	if this.top < 0 {
		this.top = 0
	} else {
		this.top += 1
	}

	if this.top > len(this.data) - 1 {
		this.data = append(this.data, v)
	} else {
		this.data[this.top] = v
	}
}

func (this *ArrayStack) Pop() interface{} {
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
		fmt.Println("Empty Array Stack")
	} else {
		for i := 0; i < len(this.data); i++ {
			fmt.Println(this.data[i])
		}
	}
}