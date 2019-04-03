package main

func buildHeap(a []int, n int) {
	//heapify from the last parent node
	for i := n / 2; i >= 1; i-- {
		heapifyUpToDown(a, i, n)
	}
}

//heapify from up to down , node index = top
func heapifyUpToDown(arr []int, top int, count int) {

	for i := top; i <= count/2; {

		maxIndex := i
		if arr[i] < a[i*2] {
			maxIndex = i * 2
		}

		if i*2+1 <= count && arr[maxIndex] < arr[i*2+1] {
			maxIndex = i*2 + 1
		}

		if maxIndex == i {
			break
		}

		swap(arr, i, maxIndex)
		i = maxIndex
	}

}

//sort by ascend, a index begin from 1, has n elements
func sort(a []int, n int) {
	buildHeap(a, n)

	k := n
	for k >= 1 {
		swap(a, 1, k)
		heapifyUpToDown(a, 1, k-1)
		k--
	}
}
