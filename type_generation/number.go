package typegeneration

import (
	"encoding/json"
	"math"
	"math/rand/v2"

	gdt "github.com/D-Gine/Gestionnaire2Type"
)

func generateNumber(restriction gdt.NumberRestrictions) (json.RawMessage, error) {
	lo := int64(math.MinInt32)
	if restriction.Min != nil {
		lo = int64(*restriction.Min)
	}

	hi := int64(math.MaxInt32)
	if restriction.Max != nil {
		hi = int64(*restriction.Max)
	}

	if lo > hi {
		lo, hi = hi, lo
	}
	span := uint64(hi) - uint64(lo)
	var offset uint64
	if span == math.MaxUint64 {
		offset = rand.Uint64()
	} else {
		offset = rand.Uint64N(span + 1)
	}
	value := int64(uint64(lo) + offset)

	return json.Marshal(value)
}
