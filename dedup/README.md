# Dedup

A simple probabilistic deduplicator built on two Bloom filters.

The deduplicator keeps two filters: current and previous. New items are added to the current filter. After the window time has passed and `RotateIfNeeded` is called, the current filter becomes previous and the old previous filter is cleared and reused as current. If more than two windows have passed, both filters are cleared.

## Guarantees and limitations

- If `Seen` returns `false`, the item is definitely not present.
- If `SeenOrAdd` returns `false`, the item was not present before the call and has now been added.
- If `SeenOrAdd` or `Seen` returns `true`, the item may be in the filter.
- After one rotation, items added during the previous window may still be reported as seen.
- After more than two windows without rotation, both filters are cleared on the next rotation check.
- False positives are possible.
- False negatives are not expected if the deduplicator is used correctly.
- The deduplicator is not safe for concurrent use without external synchronization.

## Example

```go
package main

import (
	"fmt"
	"time"

	"github.com/childofemptiness/flowkit/dedup"
)

func main() {
	now := time.Now()
	window := time.Minute
	d, err := dedup.NewAt(dedup.Options{
		Capacity:          1000,
		FalsePositiveRate: 0.01,
		Window:            window,
	},
		now,
	)
	if err != nil {
		panic(err)
	}

	data := []byte("hello")
	
	fmt.Println(d.SeenOrAdd(data)) // false
	fmt.Println(d.SeenOrAdd(data)) // true

	d.RotateIfNeeded(now.Add(window))
	fmt.Println(d.Seen(data)) // true
	
	d.RotateIfNeeded(now.Add(2 * window))
	fmt.Println(d.Seen(data)) // false
	
	d.Add(data)
	fmt.Println(d.Seen(data)) // true
	
	d.Reset()
	fmt.Println(d.Seen(data)) // false
}
```

## Parameters

- `Capacity` - expected number of items within the deduplication window.
- `FalsePositiveRate` - target false positive probability.

## Notes

The actual false positive rate depends on how many items are added. Adding significantly more items than `Capacity` will increase the false positive rate.
