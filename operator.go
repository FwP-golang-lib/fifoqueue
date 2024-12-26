package queue

func (f *fifoQueue[T]) Push(elem T) {
	f.lock.Lock()
	defer f.lock.Unlock()

	if f.backBlock.IsFull() {
		f.lastBlockSize = min(f.lastBlockSize*2, f.maxBlockSize)
		if f.backBlock.GetNext().IsFull() {
			f.backBlock.AppendNewBlock(f.lastBlockSize)
		}
		f.backBlock = f.backBlock.GetNext()
	}
	f.backBlock.Push(elem)
	f.size++
}

func (f *fifoQueue[T]) Pop() (result T, err error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	if f.size == 0 {
		return result, ErrEmptyQueue
	}

	if f.frontBlock.IsEmpty() {
		nextFrontBlock := f.frontBlock.GetNext()
		f.frontBlock.UnlinkSelf()
		f.frontBlock = nextFrontBlock
	}

	result = f.frontBlock.Pop()
	f.size--
	return result, nil
}

func (f *fifoQueue[T]) Size() int {
	f.lock.RLock()
	defer f.lock.RUnlock()

	return f.size
}
