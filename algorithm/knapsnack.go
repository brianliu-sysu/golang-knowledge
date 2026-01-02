package algorithm

func knapsack01(weights []int, values []int, capacity int) int {
	n := len(weights)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, capacity+1)
	}

	for i := 1; i <= n; i++ {
		for j := 0; j <= capacity; j++ {
			if j < weights[i-1] {
				dp[i][j] = dp[i-1][j]
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i-1][j-weights[i-1]]+values[i-1])
			}
		}
	}
	return dp[n][capacity]
}

func knapsack01_1D(weights []int, values []int, capacity int) ([]int, int) {
	n := len(weights)
	dp := make([]int, capacity+1)

	for i := 0; i < n; i++ {
		for j := capacity; j >= weights[i]; j-- {
			dp[j] = max(dp[j], dp[j-weights[i]]+values[i])
		}
	}
	return dp, dp[capacity]
}
