package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"

	"github.com/kshitijmohan15/hll/internal/hll"
)

func main() {
	const realCount = 1_000_000_000
	const precision = 16

	fmt.Printf("Testing HLL with %d truly random unique items (p=%d)\n", realCount, precision)

	h := hll.NewHll(precision)

	rand.Seed(time.Now().UnixNano()) // different seed every run

	start := time.Now()

	for range realCount {
		randomID := rand.Uint64()

		var buf [8]byte
		binary.LittleEndian.PutUint64(buf[:], randomID)
		h.Add(buf[:])
	}

	elapsed := time.Since(start)
	estimate := h.Count()

	errorPct := 100.0 * (float64(estimate) - float64(realCount)) / float64(realCount)

	fmt.Printf("\nResults (random inputs):\n")
	fmt.Printf("   Real uniques      : %d\n", realCount)
	fmt.Printf("   HLL estimate      : %d\n", estimate)
	fmt.Printf("   Error             : %+.4f%% \n", errorPct)
	fmt.Printf("   Time              : %v\n", elapsed)
	fmt.Printf("   Memory            : ~%d KB\n", (1<<precision)/1024)

	if abs(errorPct) < 1.0 {
		fmt.Printf("   PASSED: Error within ±1.0%% (theoretical bound: ~±0.81%%)\n")
	} else {
		fmt.Printf("   FAILED: Error too high!\n")
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
