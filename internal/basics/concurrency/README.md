# 🏎️ Concurrency in Go

Go is famous for its first-class support for concurrency. Unlike many other languages that use OS threads, Go uses **Goroutines**—lightweight threads managed by the Go runtime.

---

## 1. Goroutines

A goroutine is a function that is capable of running concurrently with other functions. They are extremely cheap (starting at ~2KB of stack memory).

### 🖼️ Conceptual View
```text
  Main Goroutine           New Goroutine
  +------------+           +------------+
  |  Step 1    |           |            |
  |  Step 2    | ----+     |            |
  |  Step 3    |     |     |  Parallel  |
  |  Step 4    | <---+     |  Task      |
  |  ...       |           |            |
  +------------+           +------------+
```

### 📝 Example
```go
go doSomething() // Starts a new goroutine
```

---

## 2. Channels (Communication)

"Don't communicate by sharing memory; share memory by communicating."

Channels are the pipes that connect concurrent goroutines. You can send values into channels from one goroutine and receive those values into another goroutine.

### 🖼️ Unbuffered vs Buffered
```text
  Unbuffered (Blocking)           Buffered (Non-blocking until full)
  +-------+       +-------+       +-------+       +-------+
  |  G1   | --|-->|  G2   |       |  G1   | --[###]-->|  G2   |
  +-------+       +-------+       +-------+       +-------+
    Wait for Handshake               Fill the Queue
```

### 📝 Example
```go
ch := make(chan string)    // Unbuffered
ch := make(chan int, 100)  // Buffered (capacity 100)

ch <- "hello" // Send
msg := <-ch   // Receive
```

---

## 3. The `select` Statement

The `select` statement lets a goroutine wait on multiple communication operations. It's like a `switch` but for channels.

```go
select {
case msg1 := <-ch1:
    fmt.Println("Received", msg1)
case ch2 <- "hi":
    fmt.Println("Sent hi")
case <-time.After(time.Second):
    fmt.Println("Timed out")
}
```

---

## 4. Sync Primitives (Sharing Memory Safely)

While channels are preferred, sometimes you need traditional locking.

| Primitive | Purpose |
|-----------|---------|
| `sync.WaitGroup` | Wait for a collection of goroutines to finish. |
| `sync.Mutex` | Mutual exclusion lock (only one goroutine at a time). |
| `sync.RWMutex` | Allows multiple readers OR one writer. |
| `sync.Once` | Ensures a function runs exactly once. |

---

## 5. Common Patterns

### Worker Pool
Distribute tasks across a fixed number of workers to limit resource usage.

### Fan-out / Fan-in
- **Fan-out**: Multiple functions reading from the same channel until it's closed.
- **Fan-in**: A function reads from multiple inputs and multiplexes them onto a single channel.

---

## ⚠️ Critical Pitfalls

1.  **Goroutine Leaks**: Starting a goroutine that never finishes or gets cleaned up.
2.  **Race Conditions**: Two goroutines accessing shared memory without synchronization. **Always run tests with `-race`!**
3.  **Deadlocks**: All goroutines are asleep/blocked, waiting for each other.
4.  **Closing Channels**:
    - Never close from the receiver side.
    - Never close if there are multiple concurrent senders.
    - Closing a closed channel causes a panic.

---

## 🧪 Running the Examples

Explore the unit tests for runnable patterns:
- `concurrency_test.go`: Basics, Mutexes, and Channel behaviors.

```bash
# Run with race detector (Highly Recommended)
go test -v -race ./internal/basics/concurrency/...
```
