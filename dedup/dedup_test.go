package dedup

import (
	"math"
	"testing"
	"time"

	"github.com/childofemptiness/flowkit/bloom"
	"github.com/stretchr/testify/require"
)

func TestNewAt(t *testing.T) {
	tests := []struct {
		name              string
		capacity          uint64
		falsePositiveRate float64
		window            time.Duration
		wantErr           error
	}{
		{name: "window value is less than zero", capacity: 1, falsePositiveRate: 0.1, window: -1 * time.Second, wantErr: ErrInvalidWindow},
		{name: "window value is zero", capacity: 1, falsePositiveRate: 0.1, window: time.Duration(0), wantErr: ErrInvalidWindow},
		{name: "capacity is zero", capacity: 0, falsePositiveRate: 0.1, window: time.Second, wantErr: bloom.ErrInvalidExpectedItems},
		{name: "false positive rate is less than zero", capacity: 1, falsePositiveRate: -1.0, window: time.Second, wantErr: bloom.ErrInvalidFalsePositiveRate},
		{name: "false positive rate is zero", capacity: 1, falsePositiveRate: 0.0, window: time.Second, wantErr: bloom.ErrInvalidFalsePositiveRate},
		{name: "false positive rate is one", capacity: 1, falsePositiveRate: 1.0, window: time.Second, wantErr: bloom.ErrInvalidFalsePositiveRate},
		{name: "false positive rate is greater than one", capacity: 1, falsePositiveRate: 1.1, window: time.Second, wantErr: bloom.ErrInvalidFalsePositiveRate},
		{name: "false positive rate is NaN", capacity: 1, falsePositiveRate: math.NaN(), window: time.Second, wantErr: bloom.ErrInvalidFalsePositiveRate},
		{name: "false positive rate is Inf", capacity: 1, falsePositiveRate: math.Inf(1), window: time.Second, wantErr: bloom.ErrInvalidFalsePositiveRate},
		{name: "valid options", capacity: 1, falsePositiveRate: 0.1, window: time.Second, wantErr: nil},
	}

	now := time.Date(2026, 5, 26, 16, 24, 56, 0, time.UTC)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			d, err := NewAt(Options{
				Capacity:          tt.capacity,
				FalsePositiveRate: tt.falsePositiveRate,
				Window:            tt.window,
			},
				now,
			)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, d)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, d)
			require.NotNil(t, d.current)
			require.NotNil(t, d.previous)
			require.Equal(t, tt.window, d.window)
			require.Equal(t, now, d.lastRotation)
		})
	}
}

func TestShouldNotBeSeenWhenEmptyFilters(t *testing.T) {
	d := newSampleDedup(t, time.Nanosecond, time.Now())

	require.False(t, d.Seen([]byte("hello")))
}

func TestShouldBeSeenAfterHalfRotation(t *testing.T) {
	now := time.Date(2026, 5, 26, 16, 24, 56, 0, time.UTC)
	window := time.Nanosecond
	d := newSampleDedup(t, window, now)

	data := []byte("hello")

	require.False(t, d.SeenOrAdd(data))
	require.True(t, d.SeenOrAdd(data))

	require.True(t, d.RotateIfNeeded(now.Add(window)))
	require.True(t, d.SeenOrAdd(data))
}

func TestShouldNotBeSeenAfterFullRotation(t *testing.T) {
	now := time.Date(2026, 5, 26, 16, 24, 56, 0, time.UTC)
	window := time.Nanosecond
	d := newSampleDedup(t, window, now)

	data := []byte("hello")

	require.False(t, d.SeenOrAdd(data))
	require.True(t, d.SeenOrAdd(data))

	require.True(t, d.RotateIfNeeded(now.Add(2*window)))

	require.False(t, d.SeenOrAdd(data))
}

func TestShouldNotBeSeenAfterBigWindow(t *testing.T) {
	now := time.Date(2026, 5, 26, 16, 24, 56, 0, time.UTC)
	window := time.Nanosecond
	d := newSampleDedup(t, window, now)

	data := []byte("hello")

	require.False(t, d.SeenOrAdd(data))
	require.True(t, d.SeenOrAdd(data))
	require.True(t, d.RotateIfNeeded(now.Add(10*window)))
	require.False(t, d.SeenOrAdd(data))
}

func TestDoesNotRotateBeforeElapsed(t *testing.T) {
	now := time.Date(2026, 5, 26, 16, 24, 56, 0, time.UTC)
	d := newSampleDedup(t, time.Nanosecond, now)

	data := []byte("hello")

	require.False(t, d.SeenOrAdd(data))
	require.False(t, d.RotateIfNeeded(now))
	require.True(t, d.SeenOrAdd(data))
}

func TestResetClearsDedup(t *testing.T) {
	d := newSampleDedup(t, time.Nanosecond, time.Now())

	data := []byte("hello")

	d.SeenOrAdd(data)
	require.True(t, d.SeenOrAdd(data))

	d.Reset()

	require.False(t, d.SeenOrAdd(data))
}

func newSampleDedup(t *testing.T, window time.Duration, now time.Time) *Dedup {
	d, err := NewAt(Options{
		Capacity:          100,
		FalsePositiveRate: 0.01,
		Window:            window,
	},
		now,
	)
	require.NoError(t, err)

	return d
}
