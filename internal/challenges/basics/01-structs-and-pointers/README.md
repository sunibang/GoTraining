# Challenge 01: Structs and Pointers

## Goal
Implement a basic `Person` struct and a function that updates the person's age. This challenge will test your understanding of structs and when to use pointers to mutate state.

## Tasks
1. Define a `Person` struct with `Name` (string) and `Age` (int) fields.
2. Implement a function `UpdateAge(p *Person, newAge int)` that updates the person's age.
3. Observe what happens if you pass a value instead of a pointer.

## Run Tests
```bash
go test ./internal/challenges/basics/01-structs-and-pointers/...
```
