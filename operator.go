package queue

func (f *fifoQueue[T]) Push(elem T) {
	f.lock.Lock()
	defer f.lock.Unlock()

	if f.backBlock.isFull() {
		newBlock := &fixedSizeBlock[T]{
			blockElems: make([]T, 0, f.blockSize),
		}
		f.backBlock.setNext(newBlock)
		f.backBlock = newBlock
	}
	f.backBlock.push(elem)
	f.size++
}

func (f *fifoQueue[T]) Pop() (result T, err error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	if f.size == 0 {
		return result, ErrEmptyQueue
	}

	if f.frontBlock.isEmpty() {
		f.frontBlock = f.frontBlock.getNext()
	}

	result = f.frontBlock.pop()
	f.size--
	return result, nil
}

func (f *fifoQueue[T]) Size() int {
	f.lock.RLock()
	defer f.lock.RUnlock()

	return f.size
}
