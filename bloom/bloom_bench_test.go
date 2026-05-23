package bloom

import (
	"errors"
	"strconv"
	"testing"
)

var _ bool

func BenchmarkFilter_Add(b *testing.B) {
	sizes := []int{1024, 4096, 65536}

	for _, size := range sizes {
		b.Run(strconv.Itoa(size), func(b *testing.B) {
			f, err := New(Options{
				ExpectedItems:     uint64(size),
				FalsePositiveRate: 0.1,
			})
			if err != nil {
				b.Fatal(err)
			}

			items := make([][]byte, size)
			for i := 0; i < size; i++ {
				items[i] = []byte("item" + strconv.Itoa(i))
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				f.Add(items[i%size])
			}
		})
	}
}

func BenchmarkFilter_MightContainHit(b *testing.B) {
	sizes := []int{1024, 4096, 65536}

	for _, size := range sizes {
		b.Run(strconv.Itoa(size), func(b *testing.B) {
			f, err := New(Options{
				ExpectedItems:     uint64(size),
				FalsePositiveRate: 0.1,
			})
			if err != nil {
				b.Fatal(err)
			}

			items := make([][]byte, size)
			for i := 0; i < size; i++ {
				items[i] = []byte("item" + strconv.Itoa(i))
			}

			for i := range items {
				f.Add(items[i])
			}

			for i := range items {
				if !f.MightContain(items[i]) {
					b.Fatal(errors.New(string(items[i]) + " not found"))
				}
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_ = f.MightContain(items[i%size])
			}
		})
	}
}

func BenchmarkFilter_MightContainMiss(b *testing.B) {
	sizes := []int{1024, 4096, 65536}
	rates := []float64{0.1, 0.01, 0.001}

	for _, size := range sizes {
		for _, rate := range rates {
			b.Run(strconv.Itoa(size)+" "+strconv.FormatFloat(rate, 'f', 3, 64), func(b *testing.B) {
				f, err := New(Options{
					ExpectedItems:     uint64(size),
					FalsePositiveRate: rate,
				})
				if err != nil {
					b.Fatal(err)
				}

				addableItems := make([][]byte, size)
				checkableItems := make([][]byte, size)
				for i := 0; i < size; i++ {
					addableItems[i] = []byte("added:" + strconv.Itoa(i))
					checkableItems[i] = []byte("missing:" + strconv.Itoa(i))
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
	}
}
