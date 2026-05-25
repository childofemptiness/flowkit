package bloom

import (
	"errors"
	"hash/maphash"
	"math"

	"github.com/childofemptiness/flowkit/internal/bitset"
)

var (
	ErrInvalidExpectedItems     = errors.New("expected items must be greater than 0")
	ErrInvalidFalsePositiveRate = errors.New("false positive rate must be between 0 and 1")
)

type Options struct {
	ExpectedItems     uint64
	FalsePositiveRate float64
}

type Filter struct {
	bits bitset.BitSet

	m uint64 // number of bits
	k uint64 // number of hash probes

	seed1 maphash.Seed
	seed2 maphash.Seed
}

func New(opt Options) (*Filter, error) {
	if opt.ExpectedItems == 0 {
		return nil, ErrInvalidExpectedItems
	}

	if opt.FalsePositiveRate <= 0 || opt.FalsePositiveRate >= 1 || math.IsNaN(opt.FalsePositiveRate) ||
		math.IsInf(opt.FalsePositiveRate, 0) {
		return nil, ErrInvalidFalsePositiveRate
	}

	n := float64(opt.ExpectedItems)
	ln2 := math.Log(2)
	m := math.Ceil((-n * math.Log(opt.FalsePositiveRate)) / (ln2 * ln2))
	k := math.Max(math.Round((m/n)*ln2), 1)

	seed1 := maphash.MakeSeed()
	seed2 := maphash.MakeSeed()

	bs, err := bitset.New(uint64(m))
	if err != nil {
		return nil, err
	}

	return &Filter{
		bits:  *bs,
		m:     uint64(m),
		k:     uint64(k),
		seed1: seed1,
		seed2: seed2,
	}, nil
}

func (f *Filter) Add(data []byte) {
	h1 := hash(data, f.seed1)
	h2 := hash(data, f.seed2)

	var i uint64
	for i = 0; i < f.k; i++ {
		f.bits.Set(f.index(h1, h2, i))
	}
}

func (f *Filter) MightContain(data []byte) bool {
	h1 := hash(data, f.seed1)
	h2 := hash(data, f.seed2)

	var i uint64
	for i = 0; i < f.k; i++ {
		if !f.bits.Get(f.index(h1, h2, i)) {
			return false
		}
	}

	return true
}

func (f *Filter) Reset() {
	f.bits.Clear()
}

func (f *Filter) BitSize() uint64 {
	return f.m
}

func (f *Filter) HashCount() uint64 {
	return f.k
}

func hash(data []byte, seed maphash.Seed) uint64 {
	var h maphash.Hash

	h.SetSeed(seed)

	_, _ = h.Write(data)

	return h.Sum64()
}

func (f *Filter) index(h1, h2, i uint64) uint64 {
	if h2 == 0 {
		h2 = 1
	}

	return (h1 + i*h2) % f.m
}
