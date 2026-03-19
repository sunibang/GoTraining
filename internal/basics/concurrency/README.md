# 🏎️ Concurrency in Go

Go is famous for its first-class support for concurrency. Unlike many other languages that use OS threads, Go uses **Goroutines**—lightweight threads managed by the Go runtime—and **Channels** to coordinate between them.

---

## 1. Core Concepts

| Concept | Description / Purpose |
| :--- | :--- |
| **Goroutine** | A lightweight thread (~2KB) started with the `go` keyword. |
| **Channel** | A thread-safe pipe for communication (`ch <- val`, `<-ch`). |
| **`select`** | A control structure that lets you wait on multiple channel operations. |
| **`sync` Package** | Traditional primitives like `WaitGroup` and `Mutex` for shared-memory sync. |

---

## 2. 🖼️ Visual Representation

### Communicating by Sharing Memory vs. Sharing Memory by Communicating
Go's philosophy: "Don't communicate by sharing memory; share memory by communicating."

```text
  Unbuffered Channel (Handshake)           Buffered Channel (Queue)
  +-------+       +-------+               +-------+       +-------+
  |  G1   | --|-->|  G2   |               |  G1   | --[###]-->|  G2   |
  +-------+       +-------+               +-------+       +-------+
    Blocks until both are ready             Blocks only when full
```

---

## 3. 📝 Implementation Examples

### Goroutines and Channels

```go
func main() {
    ch := make(chan string)

    // Start a new goroutine
    go func() {
        ch <- "Hello from concurrent world!"
    }()

    // Receive from channel (blocks until data arrives)
    msg := <-ch
    fmt.Println(msg)
}
```

### The `select` statement

```go
select {
case msg := <-ch1:
    fmt.Println("Received:", msg)
case <-time.After(time.Second):
    fmt.Println("Timed out")
}
```

---

## 4. 🚀 Common Patterns & Use Cases

- **Worker Pool**: Distributing tasks across a fixed number of goroutines to limit resource usage.
- **Fan-out / Fan-in**: Parallelising work across multiple goroutines and then aggregating the results.
- **Cancellation**: Using a `done` channel to stop long-running background tasks when no longer needed.

---

## 5. ⚠️ Critical Pitfalls & Best Practices

> [!WARNING]
> Accessing shared variables from multiple goroutines without synchronization causes **Race Conditions**. Always run your tests with the `-race` flag.

1.  **Goroutine Leaks**: Starting a goroutine that never finishes or gets cleaned up.
2.  **Race Conditions**: Two goroutines accessing shared memory without synchronization. **Always run tests with `-race`!**
3.  **Deadlocks**: All goroutines are asleep/blocked, waiting for each other.
4.  **Closing Channels**:
    - Never close from the receiver side.
    - Never close if there are multiple concurrent senders.
    - Closing a closed channel causes a panic.

---

## 🧪 Running the Examples

Explore the unit tests for runnable patterns covering basic channels, mutexes, and the race detector.

```bash
# Run with the race detector (Highly Recommended)
go test -v -race ./internal/basics/concurrency/...
```

---

## 📚 Further Reading

- [Effective Go: Concurrency](https://go.dev/doc/effective_go#concurrency)
- [Go Blog: Pipelines and Cancellation](https://go.dev/blog/pipelines)
