package bloom

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name              string
		expectedItems     uint64
		falsePositiveRate float64
		wantErr           error
		wantM             uint64
		wantK             uint64
	}{
		{name: "expected items is zero", expectedItems: 0, falsePositiveRate: 0.1, wantErr: ErrInvalidExpectedItems},
		{name: "valid options", expectedItems: 1, falsePositiveRate: 0.1, wantM: 5, wantK: 3},
		{name: "false positive rate is negative", expectedItems: 1, falsePositiveRate: -0.1, wantErr: ErrInvalidFalsePositiveRate},
		{name: "false positive rate is zero", expectedItems: 1, falsePositiveRate: 0.0, wantErr: ErrInvalidFalsePositiveRate},
		{name: "false positive rate is one", expectedItems: 1, falsePositiveRate: 1, wantErr: ErrInvalidFalsePositiveRate},
		{name: "false positive rate greater than one", expectedItems: 1, falsePositiveRate: 1.1, wantErr: ErrInvalidFalsePositiveRate},
		{name: "false positive rate is NaN", expectedItems: 1, falsePositiveRate: math.NaN(), wantErr: ErrInvalidFalsePositiveRate},
		{name: "false positive rate is Inf", expectedItems: 1, falsePositiveRate: math.Inf(1), wantErr: ErrInvalidFalsePositiveRate},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := New(Options{
				ExpectedItems:     tt.expectedItems,
				FalsePositiveRate: tt.falsePositiveRate,
			})

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantM, f.BitSize())
			require.Equal(t, tt.wantK, f.HashCount())
		})
	}
}

func TestEmptyFilterDoesNotContainItem(t *testing.T) {
	f := newSampleFilter(t, 100, 0.1)

	require.False(t, f.MightContain([]byte("hello")))
}

func TestAddedItemMightBeContained(t *testing.T) {
	f := newSampleFilter(t, 100, 0.1)

	data := []byte("hello")
	f.Add(data)

	require.True(t, f.MightContain(data))
}

func TestMultipleItemsHaveNoFalseNegatives(t *testing.T) {
	itemsCount := 1000
	f := newSampleFilter(t, uint64(itemsCount), 0.1)

	items := make([][]byte, itemsCount)

	for i := range items {
		items[i] = []byte("item-" + strconv.Itoa(i))
		f.Add(items[i])
	}

	for i := range items {
		require.True(t, f.MightContain(items[i]), "false negative for item %d", i)
	}
}

func TestResetClearsFilter(t *testing.T) {
	f := newSampleFilter(t, 100, 0.1)

	data := []byte("hello")
	f.Add(data)

	require.True(t, f.MightContain(data))

	f.Reset()

	require.False(t, f.MightContain(data))
}

func newSampleFilter(t *testing.T, items uint64, rate float64) *Filter {
	f, err := New(Options{
		ExpectedItems:     items,
		FalsePositiveRate: rate,
	})
	require.NoError(t, err)

	return f
}
