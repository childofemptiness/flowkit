package bitset

import (
	"errors"
)

var InvalidBitSetSizeErr = errors.New("size must be greater than zero")

type BitSet struct {
	words []uint64
	size  uint64
}

func New(size uint64) (*BitSet, error) {
	if size == 0 {
		return nil, InvalidBitSetSizeErr
	}

	wordsCount := size / 64

	if size%64 != 0 {
		wordsCount++
	}

	return &BitSet{words: make([]uint64, wordsCount), size: size}, nil
}

func (b *BitSet) Set(i uint64) {
	if i >= b.size {
		panic("bitset: index out of range")
	}

	wordIdx := i / 64
	bitIdx := i % 64

	b.words[wordIdx] |= uint64(1) << bitIdx
}

func (b *BitSet) Get(i uint64) bool {
	if i >= b.size {
		panic("bitset: index out of range")
	}

	wordIdx := i / 64
	bitIdx := i % 64

	return b.words[wordIdx]&(uint64(1)<<bitIdx) != 0
}

func (b *BitSet) Clear() {
	for i := range b.words {
		b.words[i] = 0
	}
}

func (b *BitSet) Len() uint64 {
	return b.size
}
