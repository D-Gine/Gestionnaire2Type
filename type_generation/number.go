package typegeneration

import (
	"encoding/json"
	"math"
	"math/rand/v2"

	gdt "github.com/D-Gine/Gestionnaire2Type"
)

func generateNumber(restriction gdt.NumberRestrictions) (json.RawMessage, error) {
	min := int64(math.MinInt32)
	if restriction.Min != nil {
		min = int64(*restriction.Min)
	}

	max := int64(math.MaxInt32)
	if restriction.Max != nil {
		max = int64(*restriction.Max)
	}

	if min > max {
		min, max = max, min
	}
	valueRange := (max - min) + 1
	value := int32(min + rand.Int64N(valueRange))

	jsonValue, err := json.Marshal(value)
	if err != nil {
		return jsonValue, err
	} else {
		return jsonValue, err
	}
}
