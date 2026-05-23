package bitset

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBitSet_NewInvalidSize(t *testing.T) {
	_, err := New(0)
	require.ErrorIs(t, err, InvalidBitSetSizeErr)
}

func TestBitSet_NewSuccess(t *testing.T) {
	tests := []struct {
		name               string
		size               uint64
		expectedWordsCount int
	}{
		{name: "one bit", size: 1, expectedWordsCount: 1},
		{name: "less than one word", size: 63, expectedWordsCount: 1},
		{name: "exactly one word", size: 64, expectedWordsCount: 1},
		{name: "one bit over word", size: 65, expectedWordsCount: 2},
		{name: "exactly two words", size: 128, expectedWordsCount: 2},
		{name: "one bit over two words", size: 129, expectedWordsCount: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs, err := New(tt.size)
			require.NoError(t, err)
			require.Equal(t, tt.size, bs.size)
			require.Len(t, bs.words, tt.expectedWordsCount)
		})
	}
}

func TestBitSet_SetLogicalRangePanics(t *testing.T) {
	bs, err := New(10)
	require.NoError(t, err)

	require.Panics(t, func() {
		bs.Set(63)
	})
}

func TestBitSet_Set(t *testing.T) {
	size := uint64(129)

	bs, err := New(size)
	require.NoError(t, err)

	bs.Set(128)

	require.Equal(t, uint64(1), bs.words[2])
	require.True(t, bs.Get(128))
}

func TestBitSet_GetLogicalRangePanics(t *testing.T) {
	bs, err := New(10)
	require.NoError(t, err)

	require.Panics(t, func() {
		bs.Get(63)
	})
}

func TestBitSet_Get(t *testing.T) {
	bs, err := New(uint64(130))
	require.NoError(t, err)

	bs.Set(128)
	require.True(t, bs.Get(128))
	require.False(t, bs.Get(129))
}

func TestBitSet_Boundaries(t *testing.T) {
	bs, err := New(130)
	require.NoError(t, err)

	indexes := []uint64{0, 63, 64, 127, 128, 129}

	for _, idx := range indexes {
		bs.Set(idx)
	}

	for _, idx := range indexes {
		require.True(t, bs.Get(idx))
	}
}

func TestBitSet_Len(t *testing.T) {
	size := uint64(10)

	bs, err := New(size)
	require.NoError(t, err)
	require.Equal(t, size, bs.Len())
}

func TestBitSet_Clear(t *testing.T) {
	size := uint64(129)

	bs, err := New(size)
	require.NoError(t, err)

	bs.Set(0)
	bs.Set(64)
	bs.Set(128)

	bs.Clear()

	require.False(t, bs.Get(0))
	require.False(t, bs.Get(64))
	require.False(t, bs.Get(128))
}
