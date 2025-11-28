package hll

import (
	"math"
	"math/bits"
	"sync"

	"github.com/cespare/xxhash/v2"
)

const (
	MinPrecision = 4
	MaxPrecision = 18
	two32        = 4_294_967_296.0              // 2^32
	two64        = 18_446_744_073_709_551_616.0 // 2^64
)

type HLL struct {
	precision uint8 // precision will also
	// be the first p bytes which
	// we will use to index the buckets
	numberOfBuckets uint64 // the number of buckets that will be used
	// to do the final stochastic counting, this will store the number of leading
	// zeroes in the HASHED version of the unique id which is going to be used for cardinality
	registers []uint8 // registers are the buckets that will store the maximum leading zeroes from
	// different substream of the same unique id stream so that if a black swan event hits one of the buckets,
	// the maximum in another bucket pulls it back towards the average
	mu sync.RWMutex
}

func NewHll(precision uint8) *HLL {
	// find the number of buckets based on the preferred precision
	m := uint64(1) << uint64(precision) // 2^p
	registers := make([]uint8, m)
	return &HLL{
		precision:       precision,
		numberOfBuckets: m,
		registers:       registers,
	}
}

func (h *HLL) Add(item []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()

	hash := xxhash.Sum64(item)

	j := hash >> (64 - h.precision)
	rho := uint8(bits.LeadingZeros64(hash<<h.precision)) + 1
	if rho > h.registers[j] {
		h.registers[j] = rho
	}
}

func (h *HLL) Count() uint64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	m := float64(h.numberOfBuckets)
	sum := 0.0
	empty := 0

	for _, r := range h.registers {
		sum += math.Pow(2, -float64(r))
		if r == 0 {
			empty++
		}
	}

	alpha := h.getAlpha()
	estimate := alpha * m * m / sum

	if estimate <= 2.5*m && empty > 0 {
		return uint64(m*math.Log(m/float64(empty)) + 0.5)
	}

	if estimate > two32/30 {
		return uint64(-two64*math.Log(1-estimate/two64) + 0.5)
	}

	return uint64(estimate + 0.5)
}

func (h *HLL) getAlpha() float64 {
	switch h.numberOfBuckets {
	case 16:
		return 0.673
	case 32:
		return 0.697
	case 64:
		return 0.709
	default:
		return 0.7213 / (1.0 + 1.079/float64(h.numberOfBuckets))
	}
}
