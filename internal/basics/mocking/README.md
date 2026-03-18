# Mocking Library

## Description

Interfaces in Go is used widely to loose coupling between packages. When we run unit test packages with interfaces, mocking
can be very useful to simplify the test code and focus on testing the logic within the package only.  
There are many open source mocking libraries that generates mocking code structure automatically. In this demo, we 
picked the two most popular mocking libraries in the market. 

## Mockery

https://github.com/vektra/mockery

### Pros

- It's built on top of a very popular testing library: testify
- It has a powerful CLI tool, packed with many options.
- It supports complex mocking logic and has a powerful API.

### Cons

- ~~When mocking function calls of an interface, we need to use string to match the function name.~~ ->
  you can now add a flag that will allow the same `EXPECT().` API as gomock
- The mocking structure is not type-safe.

## Go Mock

https://github.com/golang/mock

### Pros

- No dependency on any external libraries. It's a native tool.
- It has a powerful CLI tool, packed with many options.
- The mock structure created is type-safe.
- It supports complex mocking logic and has a very powerful API.

### Cons

- It doesn't provide the best error messages when error occurs.