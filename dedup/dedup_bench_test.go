package dedup

import (
	"fmt"
	"testing"
	"time"
)

var _ bool

func BenchmarkAdd(b *testing.B) {
	prepare(b, func(b *testing.B, size int, d *Dedup) {
		items := make([][]byte, size)
		for i := range items {
			items[i] = []byte(fmt.Sprintf("item-%d", i))
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			d.Add(items[i%size])
		}
	})
}

func BenchmarkSeenInEmptyFilter(b *testing.B) {
	prepare(b, func(b *testing.B, size int, d *Dedup) {
		items := make([][]byte, size)
		for i := range items {
			items[i] = []byte(fmt.Sprintf("item-%d", i))
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = d.Seen(items[i%size])
		}
	})
}

func BenchmarkSeenAfterAddHit(b *testing.B) {
	prepare(b, func(b *testing.B, size int, d *Dedup) {
		items := make([][]byte, size)
		for i := range items {
			items[i] = []byte(fmt.Sprintf("item-%d", i))
			d.Add(items[i])
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = d.Seen(items[i%size])
		}
	})
}

func BenchmarkSeenAfterAddMiss(b *testing.B) {
	prepare(b, func(b *testing.B, size int, d *Dedup) {
		addableItems := make([][]byte, size)
		checkableItems := make([][]byte, size)
		for i := 0; i < size; i++ {
			addableItems[i] = []byte(fmt.Sprintf("addable-%d", i))
			checkableItems[i] = []byte(fmt.Sprintf("checkable-%d", i))
			d.Add(addableItems[i])
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = d.Seen(checkableItems[i%size])
		}
	})
}

func BenchmarkHalfRotationHit(b *testing.B) {
	prepare(b, func(b *testing.B, size int, d *Dedup) {
		items := make([][]byte, size)
		for i := range items {
			items[i] = []byte(fmt.Sprintf("item-%d", i))
			d.Add(items[i])
		}

		d.RotateIfNeeded(d.lastRotation.Add(d.window))
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			if !d.Seen(items[i%size]) {
				b.Fatal("Should have been seen")
			}
		}
	})
}

func BenchmarkFullRotationMiss(b *testing.B) {
	prepare(b, func(b *testing.B, size int, d *Dedup) {
		items := make([][]byte, size)
		for i := range items {
			items[i] = []byte(fmt.Sprintf("item-%d", i))
			d.Add(items[i])
		}

		d.RotateIfNeeded(d.lastRotation.Add(2 * d.window))
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			if d.Seen(items[i%size]) {
				b.Fatal("should not have been seen")
			}
		}
	})
}

func prepare(b *testing.B, test func(b *testing.B, size int, d *Dedup)) {
	sizes := []int{1024, 4096, 65536}
	rates := []float64{0.1, 0.01, 0.001}
	now := time.Date(2026, 5, 26, 16, 24, 56, 0, time.UTC)

	for _, size := range sizes {
		for _, rate := range rates {
			b.Run(fmt.Sprintf("size=%d rate=%.3f", size, rate), func(b *testing.B) {
				d, err := NewAt(Options{
					Capacity:          uint64(size),
					FalsePositiveRate: rate,
					Window:            time.Nanosecond,
				}, now,
				)
				if err != nil {
					b.Fatal(err)
				}

				test(b, size, d)
			})
		}
	}
}
