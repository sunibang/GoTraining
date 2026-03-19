package benchmark

import "math/big"

// RecursiveFactorial calculates the factorial of n using recursion.
func RecursiveFactorial(n int) *big.Int {
	if n <= 1 {
		return big.NewInt(1)
	}

	res := RecursiveFactorial(n - 1)
	return res.Mul(big.NewInt(int64(n)), res)
}

// IterativeFactorial calculates the factorial of n using a loop.
func IterativeFactorial(n int) *big.Int {
	if n <= 1 {
		return big.NewInt(1)
	}

	result := big.NewInt(1)
	for i := 1; i <= n; i++ {
		result.Mul(result, big.NewInt(int64(i)))
	}
	return result
}
