package algorithm

import "testing"

func TestLongestCommonSubsequence(t *testing.T) {
	nums1 := "abcde"
	nums2 := "ace"
	result := LongestCommonSubsequence(nums1, nums2)
	if result != 3 {
		t.Fatalf("should be 3")
	}
	result = LongestCommonSubstring_1D(nums1, nums2)
	if result != 3 {
		t.Fatalf("should be 3")
	}
}
