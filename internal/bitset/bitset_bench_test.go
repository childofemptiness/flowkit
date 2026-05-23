package bitset

import (
	"fmt"
	"testing"
)

func BenchmarkBitSet_Set(b *testing.B) {
	bs, err := New(1024)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		bs.Set(uint64(i % 1024))
	}
}

func BenchmarkBitSet_Get(b *testing.B) {
	bs, err := New(1024)
	if err != nil {
		b.Fatal(err)
	}

	bs.Set(512)

	for i := 0; i < b.N; i++ {
		_ = bs.Get(512)
	}
}

func BenchmarkBitSet_Clear(b *testing.B) {
	sizes := []uint64{
		64,
		1024,
		65536,
		1_000_000,
	}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			bs, err := New(size)
			if err != nil {
				b.Fatal(err)
			}

			for i := uint64(0); i < bs.Len(); i += 64 {
				bs.Set(i)
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				bs.Clear()
			}
		})
	}
}
