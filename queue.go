package queue

import "sync"

const (
	defaultBlockSize = 256
)

func IntPtr(v int) *int {
	return &v
}

type FifoQueueConfig struct {
	BlockSize *int
}

type Queue[T any] interface {
	Push(T)
	Pop() (T, error)
	Size() int
}

type fifoQueue[T any] struct {
	lock      sync.RWMutex
	size      int
	blockSize int

	frontBlock block[T]
	backBlock  block[T]
}

func (c *FifoQueueConfig) GetBlockSize() int {
	if c == nil || c.BlockSize == nil {
		return defaultBlockSize
	}
	return *c.BlockSize
}

type fifoQueueConfigInternal struct {
	BlockSize int
}

func (c *FifoQueueConfig) toInternal() fifoQueueConfigInternal {
	return fifoQueueConfigInternal{
		BlockSize: c.GetBlockSize(),
	}
}

func NewFifoQueue[T any](c *FifoQueueConfig) Queue[T] {
	conf := c.toInternal()

	initialBlock := &initialBlock[T]{
		maxSize: conf.BlockSize,
	}
	return &fifoQueue[T]{
		blockSize:  conf.BlockSize,
		frontBlock: initialBlock,
		backBlock:  initialBlock,
	}
}
