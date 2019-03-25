package main

import (
	"math"
)

// 计数排序
func CountingSort(arr []int, n int) {
	if n <= 1{
		return
	}

	var max int = math.MaxInt32
	for i := range arr {
		if arr[i] > max {
			max = arr[i]
		}
	}

	c := make([]int, max+1)
	for i := range arr {
		c[arr[i]]++
	}
	for i := 1; i <= max; i++ {
		c[i] += c[i-1]
	}

	r := make([]int, n)
	for i := range arr {
		index := c[arr[i]] - 1
		r[index] = a[i]
		c[arr[i]]--
	}

	copy(arr, r)
}
