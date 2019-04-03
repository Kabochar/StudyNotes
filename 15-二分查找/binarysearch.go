package main

import (
	"fmt"
)

func BinarySearch(arr []int, val int) int {
	aLen := len(arr)
	if aLen == 0 {
		return -1
	}

	low := 0
	high := aLen - 1
	for low <= high {
		mid := low + ((high - low) >> 1)
		if arr[mid] == val {
			return mid
		} else if arr[mid] > val {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}

	return -1
}

func BinarySearchRecursive(arr []int, val int) int {
	aLen := len(arr)
	if aLen == 0 {
		return -1
	}

	return Bs(arr, val, 0, aLen-1)
}

func Bs(arr []int, val int, low, high int) int {
	if low > high {
		return -1
	}

	mid := low + ((high - low) >> 1)
	if arr[mid] == val {
		return mid
	} else if arr[mid] > val {
		return Bs(arr, val, low, mid-1)
	} else {
		return Bs(arr, val, mid+1, high)
	}

}

func BinarySearchFirstEle(arr []int, val int) int {
	aLen := len(arr)
	if aLen == 0 {
		return -1
	}

	low := 0
	high := aLen - 1
	for low <= high {
		mid := low + ((high - low) >> 1)
		if arr[mid] > val {
			high = mid - 1
		} else if arr[mid] < val {
			low = mid + 1
		} else {
			if mid == 0 || arr[mid-1] != val {
				return mid
			} else {
				high = mid - 1
			}
		}
	}

	return -1
}

func BinarySearchLastEle(arr []int, val int) int {
	aLen := len(arr)
	if aLen == 0 {
		return -1
	}

	low := 0
	high := aLen - 1
	for low <= high {
		mid := low + ((high - low) >> 1)
		if arr[mid] > val {
			high = mid - 1
		} else if arr[mid] < val {
			low = mid + 1
		} else {
			if mid == aLen-1 || arr[mid+1] != val {
				return mid
			} else {
				low = mid + 1
			}
		}
	}
	return -1
}

func BinarySearchFirstBigEle(arr []int, val int) int {
	aLen := len(arr)
	if aLen == 0 {
		return -1
	}

	low := 0
	high := aLen - 1
	for low <= high {
		mid := low + ((high - low) >> 1)
		if arr[mid] >= val {
			if mid == 0 || arr[mid-1] < val {
				return mid
			} else {
				high = mid - 1
			}
		} else {
			low = mid + 1
		}
	}
	return -1
}

func BinarySearchLastSmallEle(arr []int, val int) int {
	aLen := len(arr)
	if aLen == 0 {
		return -1
	}

	low := 0
	high := aLen - 1
	for low <= high {
		mid := low + ((high - low) >> 1)
		if arr[mid] > val {
			high = mid - 1
		} else {
			if mid == aLen-1 || arr[mid+1] > val {
				return mid - 1
			} else {
				low = mid + 1
			}
		}
	}
	return -1
}

func main() {
	arr := []int{2, 3, 5, 7, 8, 9}
	rlt := BinarySearchLastSmallEle(arr, 8)
	fmt.Println(arr[rlt])
}
