package main

// 冒泡排序。a 表示数组，n 表示数组大小
func BubbleSort(a []int, n int) {
	if n < 1 {
		return
	}
	for i := 0; i < n; i++ {
		flag := false // 提前退出标志
		for j := 0; j < n-i-1; j++ {
			if a[j] > a[j+1] {
				a[j], a[j+1] = a[j+1], a[j]
				flag = true
			}
		}
		if !flag {
			break // 如果没有数据交换，提前退出
		}
	}

}

// 插入排序。a 表示数组，n 表示数组大小
// 分为 已排，和 未排
func InsertionSort(a []int, n int) {
	if n < 1 {
		return
	}
	for i := 1; i < n; i++ {
		// 查找最下值
		value := a[i]
		j := i - 1
		// 查找要插入的位置并移动数据
		for ; j >= 0; j-- {
			if a[j] > value {
				a[j+1] = a[j]
			} else {
				break
			}
		}
		a[j+1] = value
	}
}

// 选择排序。a 表示数组，n 表示数组大小
func SelectionSort(a []int, n int) {
	if n <= 1 {
		return
	}
	for i := 0; i < n; i++ {
		minIndex := i
		for j := i + 1; j < n; j++ {
			if a[j] < a[minIndex] {
				minIndex = j
			}
		}
		a[i], a[minIndex] = a[minIndex], a[i]
	}
}

func main() {
	ab := []int{1, 5, 4, 6, 12, 4, 3}
	BubbleSort(ab, len(ab))
	InsertionSort(ab, len(ab))
	SelectionSort(ab, len(ab))
}
