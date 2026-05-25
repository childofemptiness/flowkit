# Bloom Filter

A simple Bloom filter implementation based on an internal bitset.

A bloom filter is a probabilistic data structure used to test whether an item is a member of a set.
It is memory-efficient, but it may return false positives.

## Guarantees and limitations

- If `MightContain` returns `false`, the item is definitely not in the filter.
- If `MightContain` returns `true`, the item may be in the filter.
- False positives are possible.
- False negatives are not expected if the filter is used correctly.
- Items cannot be removed from the filter.
- The filter is not safe for concurrent use unless external synchronization is used.

## Example

```go
package main

import (
	"github.com/childofemtiness/flowkit/bloom"
	"fmt"
)

func main() {
	f, err := bloom.New(bloom.Options{
		ExpectedItems:     1000,
		FalsePositiveRate: 0.01,
	})
	if err != nil {
		panic(err)
	}

	data := []byte("hello")

	fmt.Println(f.MightContain(data)) // false

	f.Add(data)

	fmt.Println(f.MightContain(data)) // true

	f.Reset()

	fmt.Println(f.MightContain(data)) // false
}
```

## Parameters

The filter size and number of hash probes are calculated from two input parameters:
 - `ExpectedItems` - expected number of items to be added to the filter.
 - `FalsePositiveRate` - desired false positive probability.

Internally, the filter calculates:
 - `m` - number of bits in the bitset.
 - `k` - number of hash probes used for each item.

## Formulas
$$m = \lceil\frac{-n \cdot \ln(p)}{\ln(2)^2}\rceil$$
$$k = \max(round(\frac{m}{n} \cdot \ln(2)), 1)$$

### Where:
 - `n` is the expected number of items.
 - `p` is the desired false positive rate.
 - `m` is the number of bits.
 - `k` is the number of hash probes.

## Notes

The actual false positive rate depends on how many items are added.
Adding significantly more items than ExpectedItems will increase the false positive rate.