package queue

import (
	"reflect"
	"sync"

	"github.com/FwP-golang-lib/fifoqueue/internal/block"
)

const (
	defaultMaxBlockSizeInByte = 2048
	defaultIntiialBlockSize   = 1
)

func IntPtr(v int) *int {
	return &v
}

type FifoQueueConfig[T any] struct {
	MaxBlockSize *int
}

type Queue[T any] interface {
	Push(T)
	Pop() (T, error)
	Size() int
}

type fifoQueue[T any] struct {
	lock          sync.RWMutex
	size          int
	maxBlockSize  int
	lastBlockSize int

	frontBlock *block.Block[T]
	backBlock  *block.Block[T]
}

func (c *FifoQueueConfig[T]) GetMaxBlockSize() int {
	if c == nil || c.MaxBlockSize == nil {
		typeSizeInByte := reflect.TypeOf((*T)(nil)).Elem().Size()
		if typeSizeInByte > defaultMaxBlockSizeInByte {
			return 1
		}
		return defaultMaxBlockSizeInByte / int(typeSizeInByte)
	}
	return *c.MaxBlockSize
}

type fifoQueueConfigInternal struct {
	MaxBlockSize int
}

func (c *FifoQueueConfig[T]) toInternal() fifoQueueConfigInternal {
	return fifoQueueConfigInternal{
		MaxBlockSize: c.GetMaxBlockSize(),
	}
}

func NewFifoQueue[T any](c *FifoQueueConfig[T]) Queue[T] {
	conf := c.toInternal()

	initialBlock := block.NewBlock[T](defaultIntiialBlockSize)
	return &fifoQueue[T]{
		maxBlockSize:  conf.MaxBlockSize,
		lastBlockSize: defaultIntiialBlockSize,
		frontBlock:    initialBlock,
		backBlock:     initialBlock,
	}
}
