package dedup

import (
	"errors"
	"time"

	"github.com/childofemptiness/flowkit/bloom"
)

var ErrInvalidWindow = errors.New("invalid window")

type Options struct {
	Capacity          uint64
	FalsePositiveRate float64
	Window            time.Duration
}

type Dedup struct {
	current  *bloom.Filter
	previous *bloom.Filter

	window       time.Duration
	lastRotation time.Time
}

func New(opt Options) (*Dedup, error) {
	return NewAt(opt, time.Now())
}

func NewAt(opt Options, now time.Time) (*Dedup, error) {
	if opt.Window <= 0 {
		return nil, ErrInvalidWindow
	}

	current, err := bloom.New(bloom.Options{
		ExpectedItems:     opt.Capacity,
		FalsePositiveRate: opt.FalsePositiveRate,
	})
	if err != nil {
		return nil, err
	}

	previous, err := bloom.New(bloom.Options{
		ExpectedItems:     opt.Capacity,
		FalsePositiveRate: opt.FalsePositiveRate,
	})
	if err != nil {
		return nil, err
	}

	return &Dedup{
		current:      current,
		previous:     previous,
		window:       opt.Window,
		lastRotation: now,
	}, nil
}

func (d *Dedup) Seen(key []byte) bool {
	return d.current.MightContain(key) || d.previous.MightContain(key)
}

func (d *Dedup) Add(key []byte) {
	d.current.Add(key)
}

func (d *Dedup) SeenOrAdd(key []byte) bool {
	seen := d.Seen(key)
	if seen {
		return true
	}

	d.Add(key)
	return false

}

func (d *Dedup) RotateIfNeeded(now time.Time) bool {
	elapsed := now.Sub(d.lastRotation)

	if elapsed < d.window {
		return false
	}

	if elapsed >= 2*d.window {
		d.current.Reset()
		d.previous.Reset()
		d.lastRotation = now
		return true
	}

	d.previous, d.current = d.current, d.previous
	d.current.Reset()
	d.lastRotation = now

	return true
}

func (d *Dedup) Reset() {
	d.current.Reset()
	d.previous.Reset()
}
