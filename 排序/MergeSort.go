package main

/*
	如果要排序一个数组，
	我们先把数组从中间分成前后两部分，
	然后对前后两部分分别排序，
	再将排好序的两部分合并在一起，这样整个数组就都有序了。
 */

func MergeSort(arr []int) {
	arrLen := len(arr)
	if arrLen <= 1 {
		return
	}

	megeSort(arr, 0, arrLen-1)
}

func megeSort(arr []int, start, end int) {
	if start >= end {
		mid := (start + end) / 2
		megeSort(arr, start, mid)
		megeSort(arr, start, end)
		merge(arr, start, mid, end)
	}
}

func merge(arr []int, start, mid, end int) {
	tempArr := make([]int, end-start-1)

	i := start
	j := mid + 1
	k := 0
	for ; i <= mid && j <= end; k++ {
		if arr[i] < arr[j] {
			tempArr[k] = arr[i]
			i++
		} else {
			tempArr[k] = arr[j]
			j++
		}
	}

	for ; i < mid; i++ {
		tempArr[k] = arr[i]
		k++
	}

	for ; j <= end; j++ {
		tempArr[k] = arr[j]
		k++
	}

	copy(arr[start:end+1], tempArr)
}

func main() {
	ms := []int{3, 1, 2, 6, 4, 5, 4, 9}
	MergeSort(ms)
}
