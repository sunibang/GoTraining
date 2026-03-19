# 🚀 Benchmarking in Go

Benchmarking is the process of measuring the performance of your code. Go provides a powerful, built-in benchmarking tool as part of the `testing` package.

## 📌 Why Benchmark?
- **Performance Baseline**: Know how fast your code is today.
- **Optimization**: Prove that your "faster" version actually is faster.
- **Regression Testing**: Ensure new changes don't slow down critical paths.
- **Resource Usage**: Track memory allocations and CPU efficiency.

---

## 🏗️ How Benchmarks Work

Go's benchmark runner calls your function repeatedly until it can provide a statistically significant result.

```text
  +-------------------------------------------------------+
  |                   Benchmark Workflow                  |
  +-------------------------------------------------------+
  |  1. Start with b.N = 1                                |
  |  2. Run the loop b.N times                            |
  |  3. If too fast, increase b.N (e.g., 1, 2, 5, 10...)   |
  |  4. Repeat until time limit (default 1s) is reached   |
  |  5. Calculate averages: Total Time / b.N              |
  +-------------------------------------------------------+
            |
            v
  +-------------------------------------------------------+
  | Result: BenchmarkSomeFunc  1000000  1050 ns/op         |
  +-------------------------------------------------------+
```

---

## ✍️ Anatomy of a Benchmark Function

Benchmarks must follow these rules:
1. Reside in a `_test.go` file.
2. Function name starts with `Benchmark`.
3. Take exactly one argument: `b *testing.B`.

```go
func BenchmarkMyFunction(b *testing.B) {
    // 1. Expensive Setup (optional)
    data := prepareHugeDataSet()
    
    // 2. Reset the timer so setup time isn't included
    b.ResetTimer()

    // 3. The Core Loop
    for i := 0; i < b.N; i++ {
        MyFunction(data)
    }
}
```

---

## 🏃 Running Benchmarks

Use the `go test` command with the `-bench` flag.

| Command | Description |
|---------|-------------|
| `go test -bench=.` | Run all benchmarks in current directory |
| `go test -bench=MyFunction` | Run specific benchmark |
| `go test -bench=. -benchmem` | **(Recommended)** Show memory allocations |
| `go test -bench=. -benchtime=5s` | Run for 5 seconds instead of 1 |
| `go test -bench=. -count=5` | Run 5 times (helps find variance) |

---

## 📊 Understanding the Output

When you run `go test -bench=. -benchmem`, you'll see something like this:

`BenchmarkTest128-11    965350    1312 ns/op    312 B/op    5 allocs/op`

1.  **`BenchmarkTest128-11`**: The name of the benchmark. The `-11` is the number of CPUs used (GOMAXPROCS).
2.  **`965350`**: The value of `b.N`. The function was executed ~1 million times.
3.  **`1312 ns/op`**: Average time per operation (Nanoseconds). **Lower is better.**
4.  **`312 B/op`**: Average memory allocated per operation (Bytes). **Lower is better.**
5.  **`5 allocs/op`**: Average number of heap allocations per operation. **Lower is better.**

---

## 💡 Pro Tips for Starters

### 1. `b.ResetTimer()` is your friend
Use it if you have a long setup before the loop starts.

### 2. Don't be fooled by the Compiler
If your function is too simple and its result isn't used, the Go compiler might "optimise" it away entirely, giving you 0.01 ns/op results. To prevent this, assign the result to a package-level variable:

```go
var result int

func BenchmarkAdd(b *testing.B) {
    var r int
    for i := 0; i < b.N; i++ {
        r = Add(1, 2)
    }
    result = r // Prevent compiler optimization
}
```

### 3. Iterative vs. Recursive
In the `benchmark.go` file in this directory, we compare two ways of calculating factorials:
- `IterativeFactorial` (Iterative)
- `RecursiveFactorial` (Recursive)

**Run it yourself to see which one scales better!**
```bash
go test -bench=. -benchmem ./internal/basics/benchmark/...
```

---

## 🛠️ Comparison Example

From our own benchmarks on an Apple M3 Pro (ARM64):

| Input Size | Iterative (ns/op) | Recursive (ns/op) | Performance Gap |
|------------|-------------------|-------------------|-----------------|
| 2          | ~46               | ~39               | Recursive is faster (!) |
| 16         | ~159              | ~147              | Recursive is still faster |
| 128        | ~1289             | ~1822             | **Iterative is ~40% faster** |

*Note: The optimized recursive implementation now uses in-place multiplication (`res.Mul`), significantly reducing memory allocations (5 for n=128, matching the iterative version). Iterative only pulls ahead as the input size grows and the overhead of recursive function calls starts to outweigh the loop logic.*
