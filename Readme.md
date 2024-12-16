# FIFO Queue Implementation

This repository contains a simple custom implementation of a FIFO (First-In-First-Out).

## Overview

The implementation is designed to handle queue operations efficiently, using a combination of initial blocks and fixed-size blocks. This ensures optimal performance and flexibility for various use cases.

## Key Features

- **Dynamic Block Management**:
  - Uses an initial block with configurable maximum size.
  - Dynamically adds fixed-size blocks when the queue exceeds the initial block's capacity.
- **Type-Safe Implementation**: Utilizes Go generics for type safety.
- **Extensive Test Coverage**: Includes unit tests to validate all core operations and edge cases.

## Installation

To use this package, add it to your project:

```bash
go get github.com/your-repo-name/fifoqueue
```

## Usage

### Creating a FIFO Queue

```go
q := NewFifoQueue[int](&FifoQueueConfig{
    BlockSize: IntPtr(10),
})
```

### Performing Queue Operations

- **Push an Element**:

```go
q.Push(42)
```

- **Pop an Element**:

```go
value, err := q.Pop()
if err != nil {
    fmt.Println("Queue is empty")
} else {
    fmt.Printf("Popped value: %d\n", value)
}
```

- **Check the Queue Size**:

```go
size := q.Size()
fmt.Printf("Queue size: %d\n", size)
```

## Error Handling

The package includes error handling for scenarios such as:

- Attempting to pop from an empty queue (returns `ErrEmptyQueue`).

## Contributing

1. Fork the repository.
2. Create a new branch for your feature/fix.
3. Commit your changes.
4. Submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).

