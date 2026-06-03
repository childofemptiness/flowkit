package cms

import (
	"errors"
	"hash/maphash"
	"math"
)

var (
	ErrInvalidWidth = errors.New("invalid width")
	ErrInvalidDepth = errors.New("invalid depth")
	ErrTooLarge     = errors.New("too large")
)

type Options struct {
	Width uint64
	Depth uint64
}

type Sketch struct {
	width uint64
	depth uint64
	table []uint64

	seed1 maphash.Seed
	seed2 maphash.Seed
}

func New(opt Options) (*Sketch, error) {
	if opt.Width == 0 {
		return nil, ErrInvalidWidth
	}

	if opt.Depth == 0 {
		return nil, ErrInvalidDepth
	}

	if opt.Width > math.MaxUint64/opt.Depth {
		return nil, ErrTooLarge
	}

	cells := opt.Width * opt.Depth
	if cells > uint64(math.MaxInt) {
		return nil, ErrTooLarge
	}

	return &Sketch{
		width: opt.Width,
		depth: opt.Depth,
		table: make([]uint64, int(cells)),

		seed1: maphash.MakeSeed(),
		seed2: maphash.MakeSeed(),
	}, nil
}

func (s *Sketch) Add(data []byte) {
	s.AddN(data, 1)
}

func (s *Sketch) AddN(data []byte, n uint64) {
	h1 := s.hash(data, s.seed1)
	h2 := s.hash(data, s.seed2)
	if h2 == 0 {
		h2 = 1
	}

	for row := uint64(0); row < s.depth; row++ {
		col := (h1 + row*h2) % s.width
		s.table[s.index(row, col)] += n
	}
}

func (s *Sketch) Estimate(data []byte) uint64 {
	minCount := uint64(math.MaxUint64)

	h1 := s.hash(data, s.seed1)
	h2 := s.hash(data, s.seed2)
	if h2 == 0 {
		h2 = 1
	}

	for row := uint64(0); row < s.depth; row++ {
		col := (h1 + row*h2) % s.width
		count := s.table[s.index(row, col)]
		if count < minCount {
			minCount = count
		}
	}

	return minCount
}

func (s *Sketch) Reset() {
	clear(s.table)
}

func (s *Sketch) hash(data []byte, seed maphash.Seed) uint64 {
	var h maphash.Hash
	h.SetSeed(seed)
	_, _ = h.Write(data)

	return h.Sum64()
}

func (s *Sketch) index(row, col uint64) uint64 {
	return row*s.width + col
}
