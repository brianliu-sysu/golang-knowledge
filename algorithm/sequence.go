package algorithm

func LongestCommonSubsequence(nums1 string, nums2 string) int {
	// dp[i][j] 表示 nums1[0:i] 和 nums2[0:j] 的最长公共子序列长度
	dp := make([][]int, len(nums1)+1)
	for i := range dp {
		dp[i] = make([]int, len(nums2)+1)
	}
	for i := 1; i <= len(nums1); i++ {
		for j := 1; j <= len(nums2); j++ {
			if nums1[i-1] == nums2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}
	return dp[len(nums1)][len(nums2)]
}

func LongestCommonSubstring_1D(nums1 string, nums2 string) int {
	dp := make([]int, len(nums2)+1)
	for i := 1; i <= len(nums1); i++ {
		prev := 0
		for j := 1; j <= len(nums2); j++ {
			temp := dp[j]
			if nums1[i-1] == nums2[j-1] {
				dp[j] = prev + 1
			} else {
				dp[j] = max(dp[j], dp[j-1])
			}
			prev = temp
		}
	}
	return dp[len(nums2)]
}
