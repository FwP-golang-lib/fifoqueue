package queue

type block[T any] interface {
	setNext(block[T])
	getNext() block[T]
	isFull() bool
	isEmpty() bool
	push(T)
	pop() T
}

var (
	_ block[any] = new(initialBlock[any])
	_ block[any] = new(fixedSizeBlock[any])
)

type blockElems[T any] []T

func (e *blockElems[T]) isEmpty() bool {
	return len(*e) == 0
}

func (e *blockElems[T]) push(elem T) {
	*e = append(*e, elem)
}

func (e *blockElems[T]) pop() (result T) {
	result = (*e)[0]
	*e = (*e)[1:]
	return result
}

type initialBlock[T any] struct {
	blockElems[T]

	maxSize             int // constant
	underlyingArraySize int
	next                block[T]
}

func (b *initialBlock[T]) setNext(nb block[T]) {
	b.next = nb
}

func (b *initialBlock[T]) getNext() block[T] {
	return b.next
}

func (b *initialBlock[T]) isFull() bool {
	return b.underlyingArraySize >= b.maxSize
}

func (b *initialBlock[T]) push(elem T) {
	b.blockElems.push(elem)
	b.underlyingArraySize++
}

type fixedSizeBlock[T any] struct {
	blockElems[T]

	next block[T]
}

func (b *fixedSizeBlock[T]) setNext(nb block[T]) {
	b.next = nb
}

func (b *fixedSizeBlock[T]) getNext() block[T] {
	return b.next
}

func (b *fixedSizeBlock[T]) isFull() bool {
	return len(b.blockElems) == cap(b.blockElems)
}
