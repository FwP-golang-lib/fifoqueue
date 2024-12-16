package queue

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type initialBlockSummary[T any] struct {
	BlockElems blockElems[T]

	MaxSize             int // constant
	UnderlyingArraySize int
}

type fixedSizeBlockSummary[T any] struct {
	BlockElems blockElems[T]
}

type blockSummary[T any] struct {
	Initial   *initialBlockSummary[T]
	FixedSize *fixedSizeBlockSummary[T]
}

func getBlockSummaries[T any](f *fifoQueue[T]) (blockSummaries []blockSummary[T]) {
	for block := f.frontBlock; block != nil; block = block.getNext() {
		var blockSummaryObj *blockSummary[T]
		switch blockImpl := block.(type) {
		case *initialBlock[T]:
			blockSummaryObj = &blockSummary[T]{
				Initial: &initialBlockSummary[T]{
					blockImpl.blockElems,
					blockImpl.maxSize,
					blockImpl.underlyingArraySize,
				},
			}

			// assertion to avoid human error on adding/removing field
			_ = initialBlock[T]{
				blockSummaryObj.Initial.BlockElems,
				blockSummaryObj.Initial.MaxSize,
				blockSummaryObj.Initial.UnderlyingArraySize,
				nil,
			}
		case *fixedSizeBlock[T]:
			blockSummaryObj = &blockSummary[T]{
				FixedSize: &fixedSizeBlockSummary[T]{
					blockImpl.blockElems,
				},
			}

			// assertion to avoid human error on adding/removing field
			_ = fixedSizeBlock[T]{
				blockSummaryObj.FixedSize.BlockElems,
				nil,
			}
		default:
			panic("unsupported types")
		}
		blockSummaries = append(blockSummaries, *blockSummaryObj)
	}
	return blockSummaries
}

func TestConfigFifoQueue(t *testing.T) {
	t.Parallel()
	t.Run("use default config", func(t *testing.T) {
		q := NewFifoQueue[string](nil)
		blockSummaries := getBlockSummaries(q.(*fifoQueue[string]))
		if diff := cmp.Diff(
			[]blockSummary[string]{
				{
					Initial: &initialBlockSummary[string]{
						BlockElems: nil,
						MaxSize:    256,
					},
				},
			},
			blockSummaries,
		); diff != "" {
			t.Fatal(diff)
		}

	})
}

func TestEndToEndFifoQueue(t *testing.T) {
	t.Parallel()
	t.Run(`
	initial block test:
	- [initial] test the value after pop
	- [initial] test underlying array size: make sure it not reduced even if there are pop operation
	- [initial] test the value of maxsize
	- [initial] test the value of blockElems
	- [initial] test the state when initial block full
	- [initial] test the value when pop-ing empty queue
	- [initial] test adding new block
	- [fixed size] test the value of blockElems
	- [fixed size] test whether it reach max capacity, it will add new block
	- [misc] if pop empty block, the front will be changed
	- [fixed size] test the value after pop
	- [fixed size] test the value when pop-ing empty queue
	`, func(t *testing.T) {
		q := NewFifoQueue[int](&FifoQueueConfig{
			BlockSize: IntPtr(10),
		})

		for i := 0; i < 9; i++ {
			q.Push(i)
		}

		firstPushed, err := q.Pop()
		if err != nil {
			t.Fatal(err)
		}
		if firstPushed != 0 {
			t.Fatalf("expected: %d, got: %d", 0, firstPushed)
		}

		blockSummaries := getBlockSummaries(q.(*fifoQueue[int]))
		if diff := cmp.Diff(
			[]blockSummary[int]{
				{
					Initial: &initialBlockSummary[int]{
						BlockElems:          blockElems[int]{1, 2, 3, 4, 5, 6, 7, 8},
						MaxSize:             10,
						UnderlyingArraySize: 9,
					},
				},
			},
			blockSummaries,
		); diff != "" {
			t.Fatal(diff)
		}

		if size := q.Size(); size != 8 {
			t.Errorf("expected: %d, got: %d", 8, size)
		}

		q.Push(9)

		blockSummaries = getBlockSummaries(q.(*fifoQueue[int]))
		if diff := cmp.Diff(
			[]blockSummary[int]{
				{
					Initial: &initialBlockSummary[int]{
						BlockElems:          blockElems[int]{1, 2, 3, 4, 5, 6, 7, 8, 9},
						MaxSize:             10,
						UnderlyingArraySize: 10,
					},
				},
			},
			blockSummaries,
		); diff != "" {
			t.Fatal(diff)
		}

		values := []int{}
		for i := 1; i < 10; i++ {
			res, err := q.Pop()
			if err != nil {
				t.Errorf("exp no error, got: %v, on iteration %d", err, i)
			}
			values = append(values, res)
		}
		if diff := cmp.Diff([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, values); diff != "" {
			t.Error(diff)
		}

		_, err = q.Pop()
		if err != ErrEmptyQueue {
			t.Errorf("exp error %v, got %v", ErrEmptyQueue, err)
		}

		q.Push(10)
		blockSummaries = getBlockSummaries(q.(*fifoQueue[int]))
		if diff := cmp.Diff(
			[]blockSummary[int]{
				{
					Initial: &initialBlockSummary[int]{
						BlockElems:          blockElems[int]{},
						MaxSize:             10,
						UnderlyingArraySize: 10,
					},
				},
				{
					FixedSize: &fixedSizeBlockSummary[int]{
						BlockElems: blockElems[int]{10},
					},
				},
			},
			blockSummaries,
		); diff != "" {
			t.Fatal(diff)
		}

		for i := 11; i < 21; i++ {
			q.Push(i)
		}
		blockSummaries = getBlockSummaries(q.(*fifoQueue[int]))
		if diff := cmp.Diff(
			[]blockSummary[int]{
				{
					Initial: &initialBlockSummary[int]{
						BlockElems:          blockElems[int]{},
						MaxSize:             10,
						UnderlyingArraySize: 10,
					},
				},
				{
					FixedSize: &fixedSizeBlockSummary[int]{
						BlockElems: blockElems[int]{10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
					},
				},
				{
					FixedSize: &fixedSizeBlockSummary[int]{
						BlockElems: blockElems[int]{20},
					},
				},
			},
			blockSummaries,
		); diff != "" {
			t.Fatal(diff)
		}
		if _, err := q.Pop(); err != nil {
			t.Fatal(err)
		}
		blockSummaries = getBlockSummaries(q.(*fifoQueue[int]))
		if diff := cmp.Diff(
			[]blockSummary[int]{
				{
					FixedSize: &fixedSizeBlockSummary[int]{
						BlockElems: blockElems[int]{11, 12, 13, 14, 15, 16, 17, 18, 19},
					},
				},
				{
					FixedSize: &fixedSizeBlockSummary[int]{
						BlockElems: blockElems[int]{20},
					},
				},
			},
			blockSummaries,
		); diff != "" {
			t.Fatal(diff)
		}

		values = []int{}
		for i := 1; i < 11; i++ {
			res, err := q.Pop()
			if err != nil {
				t.Fatalf("exp no error, got: %v, on iteration %d", err, i)
			}
			values = append(values, res)
		}
		if diff := cmp.Diff([]int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}, values); diff != "" {
			t.Fatal(diff)
		}

		if _, err := q.Pop(); err != ErrEmptyQueue {
			t.Fatalf("exp error: %v, got: %v", ErrEmptyQueue, err)
		}
		blockSummaries = getBlockSummaries(q.(*fifoQueue[int]))
		if diff := cmp.Diff(
			[]blockSummary[int]{
				{
					FixedSize: &fixedSizeBlockSummary[int]{
						BlockElems: blockElems[int]{},
					},
				},
			},
			blockSummaries,
		); diff != "" {
			t.Fatal(diff)
		}
	})
}
