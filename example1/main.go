package main

import queue "github.com/FwP-golang-lib/fifoqueue"

func main() {
	q := queue.NewFifoQueue[int](nil)
	q.Push(12)
	q.Pop()
	q.Size()
}
