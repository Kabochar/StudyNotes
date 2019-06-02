package main

type Heap struct {
	arr   []int
	nlen  int
	count int
}

//init heap
func NewHeap(capacity int) *Heap {
	heap := &Heap{}
	heap.nlen = capacity
	heap.arr = make([]int, capacity+1)
	heap.count = 0
	return heap
}

// top-max heap -> heapify from down to up
func (this *Heap) insertEle(data int) {
	if this.count == this.nlen {
		return
	}

	this.count++
	this.arr[this.count] = data

	i := this.count
	parent := i / 2
	for parent > 0 && this.arr[parent] < this.arr[i] {
		swap(this.arr, parent, i)
		i := parent
		parent = i / 2
	}
}

// heapfify from up to down
func (this *Heap) removeMax() {
	if this.count == 0 {
		return
	}

	swap(this.arr, 1, this.count)
	this.count--

	heapifyUpTodown(this.arr, this.count)
}

//heapify
func heapifyUpTodown(a []int, count int) {
	for i := 1; i <= count/2; {

		maxIndex := i
		if a[i] < a[i*2] {
			maxIndex = i * 2
		}

		if i*2+1 <= count && a[maxIndex] < a[i*2+1] {
			maxIndex = i*2 + 1
		}

		if maxIndex == i {
			break
		}

		swap(a, i, maxIndex)
		i = maxIndex
	}
}

//swap two elements
func swap(a []int, i int, j int) {
	tmp := a[i]
	a[i] = a[j]
	a[j] = tmp
}
