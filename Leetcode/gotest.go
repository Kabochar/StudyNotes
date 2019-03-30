package main

import (
	"fmt"
)

func isValid(s string) bool {
	size := len(s)

	stack := make([]byte, size)
	top := 0

	for i := 0; i < size; i++ {
		cs := s[i]
		switch cs {
		case '(':
			stack[top] = cs + 1
			top++
		case '[', '{':
			stack[top] = cs + 2
			top++
		case ')', ']', '}':
			if top > 0 && stack[top-1] == cs {
				top--
			} else {
				return false
			}
		}
	}
	return top == 0
}
func main() {
	fmt.Println(isValid("()"))
	fmt.Println(isValid("()[]{}"))
	fmt.Println(isValid("(]"))
	fmt.Println(isValid("([)]"))
	fmt.Println(isValid("{[]}"))
}
