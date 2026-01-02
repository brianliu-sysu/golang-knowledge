package algorithm

func QuickSort(arr []int) {
	if len(arr) <= 1 {
		return
	}

	pivot := partition(arr, 0, len(arr)-1)
	QuickSort(arr[:pivot])
	QuickSort(arr[pivot+1:])
}

func partition(arr []int, left, right int) int {
	pivot := arr[right]
	i := left - 1
	for j := left; j < right; j++ {
		if arr[j] < pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[i+1], arr[right] = arr[right], arr[i+1]
	return i + 1
}

func partition2(arr []int, low, high int) int {
	pivot := arr[low]
	left, right := low, high

	for left < right {
		for left < right && arr[right] >= pivot {
			right--
		}

		for left < right && arr[left] <= pivot {
			left++
		}

		if left < right {
			arr[right], arr[left] = arr[left], arr[right]
		}
	}

	arr[low], arr[left] = arr[left], arr[low]
	return left
}
