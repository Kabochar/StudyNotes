package main

func QuickSort(arr []int) {
	separateSort(arr, 0, len(arr)-1)
}

func separateSort(arr []int, start, end int) {
	if start >= end {
		return
	}
	mid := partition(arr, start, end)
	separateSort(arr, start, mid-1)
	separateSort(arr, mid+1, end)
}

func partition(arr []int, start, end int) int {
	// 选取最后一位当对比数字
	pivot := arr[end]
	i := start
	for j := start; j < end; j++ {
		if arr[j] < pivot {
			if !(i == j) {
				arr[i], arr[j] = arr[j], arr[i] // 交换位置
			}
			i++
		}

	}
	arr[i], arr[end] = arr[end], arr[i]

	return i
}
