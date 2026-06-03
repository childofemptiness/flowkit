# Count-Min Sketch

A count-min sketch is a memory-efficient data structure used to count occurrences of each item.

## Guarantees and limitations

- `Estimate` never underestimates the count if only positive increments are used and counters do not overflow.
- `Estimate` may overestimate because of hash collisions.
- Items cannot be removed from the sketch.
- The sketch is not safe for concurrent use without external synchronization.

## Example

```go
package main

import (
	"fmt"

	"github.com/childofemptiness/flowkit/cms"
)

func main() {
	s, err := cms.New(cms.Options{
		Width: 100,
		Depth: 4,
	})
	if err != nil {
		panic(err)
	}

	data := []byte("hello")

	s.Add(data)
	fmt.Println(s.Estimate(data)) // 1

	s.AddN(data, 3)
	fmt.Println(s.Estimate(data)) // 4

	s.Reset()
	fmt.Println(s.Estimate(data)) // 0
}
```

## Parameters

- `Width` - number of table columns. Larger values reduce overcount caused by hash collisions.
- `Depth` - number of table rows and hash probes used for each item. Larger values reduce the probability of a large overcount.

