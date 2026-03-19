# Challenge 6: Concurrent Worker Pool with Testing

**Difficulty:** Advanced
**Covers:** `concurrency`, `testing`, `receivers`, `parameters`, `testify`

---

## Goal

Build a generic worker pool that distributes jobs across N goroutines, started and stopped via `context.Context`. Write a full testify suite to verify correctness, and benchmark it with varying worker counts.

---

## Requirements

### 1. Core types

```go
// Job represents a unit of work. Use a value type here and benchmark
// whether pointer vs value makes a measurable difference.
type Job struct {
    ID      int
    Payload []byte
}

// Result holds the outcome of processing a Job.
type Result struct {
    JobID int
    Err   error
}

// Pool distributes jobs to N goroutines.
// Must use pointer receivers — Pool holds channels and a WaitGroup.
type Pool struct {
    workers int
    jobs    chan Job
    results chan Result
    wg      sync.WaitGroup
}

func NewPool(workers int, bufferSize int) *Pool
func (p *Pool) Start(ctx context.Context)        // launch worker goroutines
func (p *Pool) Submit(job Job) error             // send a job; error if pool is stopped
func (p *Pool) Results() <-chan Result           // receive results
func (p *Pool) Stop()                            // signal workers to stop and wait
```

### 2. Testify suite

```go
type PoolTestSuite struct {
    suite.Suite
    pool *Pool
}

func (s *PoolTestSuite) SetupSuite()    // create pool with 4 workers
func (s *PoolTestSuite) TearDownSuite() // stop pool
func (s *PoolTestSuite) TestSubmit1000Jobs()  // submit 1000, collect 1000 results
func (s *PoolTestSuite) TestNoJobsLost()      // verify all job IDs appear in results
```

### 3. Race detector

Run your tests with:
```bash
go test -race ./...
```
Fix any races before submitting.

### 4. (Optional) Benchmarks

```go
func BenchmarkPool_Workers1(b *testing.B)  { benchPool(b, 1) }
func BenchmarkPool_Workers4(b *testing.B)  { benchPool(b, 4) }
func BenchmarkPool_Workers16(b *testing.B) { benchPool(b, 16) }

// Also benchmark Job as value vs pointer
func BenchmarkPool_ValueJob(b *testing.B)   { ... }
func BenchmarkPool_PointerJob(b *testing.B) { ... }
```

---

## Skills Practiced

- Goroutine lifecycle management with `context.Context`
- Channel-based job distribution
- `sync.WaitGroup` for clean shutdown
- Testify suite: `SetupSuite` / `TearDownSuite`
- Race detector usage
- Benchmarking concurrent code

---

## Hints

- Workers should select on both `jobs` channel and `ctx.Done()` to respect cancellation
- Close the `results` channel only after all workers have finished (use `wg.Wait()` in a goroutine, then close)
- `Submit` should be non-blocking if the buffer is full — consider returning an error or using `select` with `default`
- For the benchmark, create a fresh pool per `b.N` iteration to avoid state leakage
