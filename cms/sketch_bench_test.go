package cms

import (
	"fmt"
	"strconv"
	"testing"
)

var res uint64

const keysCount = 1024

func BenchmarkAddOneKey(b *testing.B) {
	prepare(
		b,
		func(width, depth uint64) int {
			return keysCount
		},
		func(b *testing.B, keysCount int, s *Sketch) {
			data := []byte("item")

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				s.Add(data)
			}
		})
}

func BenchmarkAddManyKeys(b *testing.B) {
	prepare(
		b,
		func(width, depth uint64) int {
			return keysCount
		},
		func(b *testing.B, keysCount int, s *Sketch) {
			items := make([][]byte, keysCount)

			for i := range items {
				items[i] = []byte("item-" + strconv.Itoa(i))
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				s.Add(items[i%keysCount])
			}
		})
}

func BenchmarkAddManyScaled(b *testing.B) {
	prepare(
		b,
		func(width, depth uint64) int {
			return int(width)
		},
		func(b *testing.B, keysCount int, s *Sketch) {
			items := make([][]byte, keysCount)

			for i := range items {
				items[i] = []byte("item-" + strconv.Itoa(i))
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				s.Add(items[i%keysCount])
			}
		})
}

func BenchmarkEstimateHit(b *testing.B) {
	prepare(
		b,
		func(width, depth uint64) int {
			return keysCount
		},
		func(b *testing.B, keysCount int, s *Sketch) {
			items := make([][]byte, keysCount)

			for i := range items {
				items[i] = []byte("item-" + strconv.Itoa(i))
				s.Add(items[i])
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				res = s.Estimate(items[i%keysCount])
			}
		})
}

func BenchmarkEstimateMiss(b *testing.B) {
	prepare(
		b,
		func(width, depth uint64) int {
			return keysCount
		},
		func(b *testing.B, keysCount int, s *Sketch) {
			addedItems := make([][]byte, keysCount)
			missingItems := make([][]byte, keysCount)

			for i := range addedItems {
				addedItems[i] = []byte("added-" + strconv.Itoa(i))
				missingItems[i] = []byte("missing-" + strconv.Itoa(i))
				s.Add(addedItems[i])
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				res = s.Estimate(missingItems[i%keysCount])
			}
		})
}

func prepare(
	b *testing.B,
	keysCountForCase func(width, depth uint64) int,
	test func(b *testing.B, keysCount int, s *Sketch),
) {
	b.Helper()
	b.ReportAllocs()

	widths := []uint64{1024, 4096, 65536}
	depths := []uint64{4, 8, 16}

	for _, width := range widths {
		for _, depth := range depths {
			b.Run(fmt.Sprintf("width=%d depth=%d", width, depth), func(b *testing.B) {
				s, err := New(Options{
					Width: width,
					Depth: depth,
				})
				if err != nil {
					b.Fatal(err)
				}

				test(b, keysCountForCase(s.width, s.depth), s)
			})
		}
	}
}
