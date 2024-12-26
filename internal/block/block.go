package block

import "github.com/FwP-golang-lib/fifoqueue/internal/util"

type Block[T any] struct {
	elems []T

	frontIdx     int
	occupiedElem int

	next *Block[T]
	prev *Block[T]
}

func NewBlock[T any](blockSize int) (result *Block[T]) {
	result = &Block[T]{
		elems:        make([]T, blockSize),
		frontIdx:     0,
		occupiedElem: 0,
	}

	// don't let block have nil reference
	result.next = result
	result.prev = result
	return result
}

func (b *Block[T]) IsEmpty() bool {
	return b.occupiedElem == 0
}

func (b *Block[T]) IsFull() bool {
	return b.occupiedElem >= len(b.elems)
}

func (b *Block[T]) Push(elem T) {
	b.elems[(b.frontIdx+b.occupiedElem)%len(b.elems)] = elem
	b.occupiedElem++
}

func (b *Block[T]) Pop() (result T) {
	result = b.elems[b.frontIdx]
	util.Zero(&b.elems[b.frontIdx])
	b.frontIdx = (b.frontIdx + 1) % len(b.elems)
	return result
}

func (b *Block[T]) UnlinkSelf() {
	next := b.next
	prev := b.prev

	prev.next = next
	next.prev = prev
}

func (b *Block[T]) AppendNewBlock(blockSize int) *Block[T] {
	newNextBlock := NewBlock[T](blockSize)

	next2 := b.next

	b.next = newNextBlock
	newNextBlock.prev = b

	newNextBlock.next = next2
	next2.prev = newNextBlock

	return newNextBlock
}

func (b *Block[T]) GetNext() *Block[T] {
	return b.next
}
