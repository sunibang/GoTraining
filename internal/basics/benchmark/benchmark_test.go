package benchmark

import (
	"math/big"
	"testing"
)

// run: go test -bench=. -benchmem
func BenchmarkTest2(b *testing.B)   { benchmarkTest(2, b) }
func BenchmarkTest4(b *testing.B)   { benchmarkTest(4, b) }
func BenchmarkTest8(b *testing.B)   { benchmarkTest(8, b) }
func BenchmarkTest16(b *testing.B)  { benchmarkTest(16, b) }
func BenchmarkTest32(b *testing.B)  { benchmarkTest(32, b) }
func BenchmarkTest64(b *testing.B)  { benchmarkTest(64, b) }
func BenchmarkTest128(b *testing.B) { benchmarkTest(128, b) }

func BenchmarkTestRec2(b *testing.B)   { benchmarkTestRec(2, b) }
func BenchmarkTestRec4(b *testing.B)   { benchmarkTestRec(4, b) }
func BenchmarkTestRec8(b *testing.B)   { benchmarkTestRec(8, b) }
func BenchmarkTestRec16(b *testing.B)  { benchmarkTestRec(16, b) }
func BenchmarkTestRec32(b *testing.B)  { benchmarkTestRec(32, b) }
func BenchmarkTestRec64(b *testing.B)  { benchmarkTestRec(64, b) }
func BenchmarkTestRec128(b *testing.B) { benchmarkTestRec(128, b) }

func benchmarkTestRec(n int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		RecursiveFactorial(n)
	}
}

func benchmarkTest(n int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		IterativeFactorial(n)
	}
}

type testCase struct {
	name string
	n    int
	want *big.Int
}

var commonTestCases = []testCase{
	{name: "-1", n: -1, want: big.NewInt(1)},
	{name: "0", n: 0, want: big.NewInt(1)},
	{name: "1", n: 1, want: big.NewInt(1)},
	{name: "2", n: 2, want: big.NewInt(2)},
	{name: "3", n: 3, want: big.NewInt(6)},
	{name: "4", n: 4, want: big.NewInt(24)},
	{name: "70", n: 70, want: fromString("11978571669969891796072783721689098736458938142546425857555362864628009582789845319680000000000000000")},
}

func runFactorialTest(t *testing.T, fn func(int) *big.Int) {
	for _, tt := range commonTestCases {
		t.Run(tt.name, func(t *testing.T) {
			if got := fn(tt.n); got.Cmp(tt.want) != 0 {
				t.Errorf("factorial(%d): got %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}

func Test_RecursiveFactorial(t *testing.T) {
	runFactorialTest(t, RecursiveFactorial)
}

func Test_IterativeFactorial(t *testing.T) {
	runFactorialTest(t, IterativeFactorial)
}

func fromString(s string) *big.Int {
	i, _ := big.NewInt(1).SetString(s, 10)
	return i
}
