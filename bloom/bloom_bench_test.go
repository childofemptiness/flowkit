package bloom

import (
	"fmt"
	"testing"
)

var _ bool

func BenchmarkAdd(b *testing.B) {
	prepare(b, func(b *testing.B, f *Filter, size int) {
		items := make([][]byte, size)
		for i := 0; i < size; i++ {
			items[i] = []byte(fmt.Sprintf("item-%d", i))
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			f.Add(items[i%size])
		}
	})
}

func BenchmarkMightContainHit(b *testing.B) {
	prepare(b, func(b *testing.B, f *Filter, size int) {
		items := make([][]byte, size)
		for i := 0; i < size; i++ {
			items[i] = []byte(fmt.Sprintf("item-%d", i))
		}

		for i := range items {
			f.Add(items[i])
		}

		for i := range items {
			if !f.MightContain(items[i]) {
				b.Fatalf("should have been seen")
			}
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = f.MightContain(items[i%size])
		}
	})
}

func BenchmarkMightContainMiss(b *testing.B) {
	prepare(b, func(b *testing.B, f *Filter, size int) {
		addableItems := make([][]byte, size)
		checkableItems := make([][]byte, size)
		for i := 0; i < size; i++ {
			addableItems[i] = []byte(fmt.Sprintf("added-%d", i))
			checkableItems[i] = []byte(fmt.Sprintf("missing-%d", i))
		}

		for i := range addableItems {
			f.Add(addableItems[i])
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = f.MightContain(checkableItems[i%size])
		}
	})
}

func prepare(b *testing.B, test func(b *testing.B, f *Filter, size int)) {
	sizes := []int{1024, 4096, 65536}
	rates := []float64{0.1, 0.01, 0.001}

	for _, size := range sizes {
		for _, rate := range rates {
			b.Run(fmt.Sprintf("size %d rate %.3f", size, rate), func(b *testing.B) {
				f, err := New(Options{
					ExpectedItems:     uint64(size),
					FalsePositiveRate: rate,
				})
				if err != nil {
					b.Fatal(err)
				}

				test(b, f, size)
			})
		}
	}
}
