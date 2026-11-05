package typegeneration

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"strings"

	gdt "github.com/D-Gine/Gestionnaire2Type"
	"github.com/lucasjones/reggen"
)

func generateString(restriction gdt.StringRestrictions) (json.RawMessage, error) {
	if restriction.RegexExpression != nil && *restriction.RegexExpression != "" {
		generated, err := reggen.Generate(*restriction.RegexExpression, 10)
		if err != nil {
			return nil, fmt.Errorf("error generating regex: %w", err)
		}
		jsonString, err := json.Marshal(generated)
		if err != nil {
			return nil, err
		}
		return jsonString, nil
	}

	minLen := 0
	if restriction.MinLength != nil {
		minLen = max(*restriction.MinLength, 0)
	}

	maxLen := minLen + (rand.Int() % 254) + 1
	if restriction.MaxLength != nil {
		maxLen = *restriction.MaxLength
	}
	if minLen > maxLen {
		minLen, maxLen = maxLen, minLen
	}
	span := int64(maxLen - minLen + 1)
	length := minLen + int(rand.Int64N(span))

	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var sb strings.Builder
	sb.Grow(length)
	for range length {
		idx := rand.IntN(len(letterBytes))
		sb.WriteByte(letterBytes[idx])
	}

	jsonString, err := json.Marshal(sb.String())
	if err != nil {
		return nil, err
	}
	return jsonString, nil
}
