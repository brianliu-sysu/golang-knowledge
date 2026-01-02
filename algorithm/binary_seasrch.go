package algorithm

func BinarySearch(arr []int, target int) int {
	left, right := 0, len(arr)-1
	for left <= right {
		mid := (left + right) / 2
		if arr[mid] == target {
			return mid
		}

		if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return binarySearchRecursive(arr, target, left, right)
}

func binarySearchRecursive(arr []int, target int, left, right int) int {
	if left > right {
		return -1
	}
	mid := (left + right) / 2
	if arr[mid] == target {
		return mid
	}
	if arr[mid] < target {
		return binarySearchRecursive(arr, target, mid+1, right)
	}
	return binarySearchRecursive(arr, target, left, mid-1)
}
