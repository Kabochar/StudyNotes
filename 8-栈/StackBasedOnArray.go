package main

import "fmt"


type ArrayStack struct {
	data []interface{} // 数据
	top  int           // 栈顶指针
}

// 创建 顺序栈
func NewArrayStack() *ArrayStack {
	return &ArrayStack{
		data: make([]interface{}, 0, 20),
		top:  -1,
	}
}

// 检测 空栈
func (this *ArrayStack) IsEmpty() bool {
	if this.top < 0 {
		return true
	}

	return false
}

/**
关键点，
1，检查 栈顶
2，检查 容量
*/
// 压栈
func (this *ArrayStack) Push(v interface{}) {
	if this.top < 0 {
		this.top = 0
	} else {
		this.top += 1
	}
	// 判断是否容量充足
	if this.top > len(this.data)-1 {
		this.data = append(this.data, v)
	} else {
		this.data[this.top] = v
	}
}

// 退栈
func (this *ArrayStack) Pop() interface{} {
	if this.IsEmpty() {
		return nil
	}
	value := this.data[this.top]
	this.top -= 1

	return value
}

// 获取 栈顶元素
func (this *ArrayStack) Top() interface{} {
	if this.IsEmpty() {
		return nil
	}

	return this.data[this.top]
}

// 清空栈
func (this *ArrayStack) Flush() {
	this.top = -1
}

func (this *ArrayStack) Print() {
	if this.IsEmpty() {
		fmt.Println("Empty Stack")
	} else {
		for i := this.top; i >= 0; i-- {
			fmt.Println(this.data[i])
		}
	}
}



func InitArrayStack() {
	as := NewArrayStack()
	as.Push(5)
	as.Print()

	as.Push(1)
	as.Push(2)

	as.Pop()

	as.Top()

	as.Flush()

	as.Print()

	// ----------
	ab := NewArrayStack()
	for i := 0; i < 23; i++ {
		ab.Push(i)
		fmt.Println(len(ab.data), cap(ab.data))
	}

}

func main() {
	InitArrayStack()
}
