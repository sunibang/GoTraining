# Challenge 02: Interfaces and Receivers

## Goal
Implement an interface and two different structs that satisfy it. This challenge tests your understanding of Go interfaces and method receivers.

## Tasks
1. Define a `Shape` interface with an `Area() float64` method.
2. Implement a `Circle` struct with a `Radius` (float64) and a `Rectangle` struct with `Width` and `Height` (float64).
3. Implement the `Area()` method for both structs.
4. Implement a function `PrintArea(s Shape) float64` that returns the area of any shape.

## Run Tests
```bash
go test ./internal/challenges/basics/02-interfaces-and-receivers/...
```
