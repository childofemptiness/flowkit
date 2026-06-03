package cms

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		width   uint64
		depth   uint64
		wantErr error
	}{
		{name: "width is zero", width: 0, depth: 10, wantErr: ErrInvalidWidth},
		{name: "depth is zero", width: 100, depth: 0, wantErr: ErrInvalidDepth},
		{name: "width is too large", width: math.MaxUint64/10 + 1, depth: 10, wantErr: ErrTooLarge},
		{name: "depth is too large", width: 100, depth: math.MaxUint64/100 + 1, wantErr: ErrTooLarge},
		{name: "valid options", width: 100, depth: 10, wantErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := New(Options{
				Width: tt.width,
				Depth: tt.depth,
			})

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, s)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, s)
			require.Equal(t, uint64(0), s.Estimate([]byte("unknown")))
		})
	}
}

func TestAddOne(t *testing.T) {
	s := newSampleSketch(t, 100, 10)

	data := []byte("hello")

	s.Add(data)
	require.Equal(t, uint64(1), s.Estimate(data))
}

func TestAddN(t *testing.T) {
	s := newSampleSketch(t, 100, 10)

	data := []byte("hello")

	s.AddN(data, 10)
	require.Equal(t, uint64(10), s.Estimate(data))
}

func TestAddNZero(t *testing.T) {
	s := newSampleSketch(t, 100, 10)

	data := []byte("hello")

	s.AddN(data, 0)
	require.Equal(t, uint64(0), s.Estimate(data))
}

func TestAddOneAndN(t *testing.T) {
	s := newSampleSketch(t, 100, 10)

	data := []byte("hello")

	s.Add(data)
	s.AddN(data, 4)

	require.Equal(t, uint64(5), s.Estimate(data))
}

func TestAddDifferentItems(t *testing.T) {
	s := newSampleSketch(t, 100, 10)

	first := []byte("hello")
	second := []byte("world")

	s.AddN(first, 5)
	s.AddN(second, 10)

	require.True(t, s.Estimate(first) >= 5)
	require.True(t, s.Estimate(second) >= 10)
}

func TestAddDifferentItemsWithCollisions(t *testing.T) {
	s := newSampleSketch(t, 1, 10)

	first := []byte("hello")
	second := []byte("world")

	s.AddN(first, 5)
	s.AddN(second, 10)

	require.Equal(t, uint64(15), s.Estimate(first))
	require.Equal(t, uint64(15), s.Estimate(second))
}

func TestEstimateInEmptySketch(t *testing.T) {
	s := newSampleSketch(t, 100, 10)

	require.Equal(t, uint64(0), s.Estimate([]byte("unknown")))
}

func TestReset(t *testing.T) {
	s := newSampleSketch(t, 100, 10)

	data := []byte("hello")

	s.Add(data)
	require.Equal(t, uint64(1), s.Estimate(data))

	s.Reset()

	require.Equal(t, uint64(0), s.Estimate(data))
}

func newSampleSketch(t *testing.T, width, depth uint64) *Sketch {
	t.Helper()

	s, err := New(Options{
		Width: width,
		Depth: depth,
	})
	require.NoError(t, err)

	return s
}
