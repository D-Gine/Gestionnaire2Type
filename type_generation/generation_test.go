package typegeneration

import (
	"encoding/json"
	"math"
	"regexp"
	"testing"
	"unicode/utf8"

	gdt "github.com/D-Gine/Gestionnaire2Type"
)

func ptr[T any](v T) *T { return &v }

func TestGenerateNumber(t *testing.T) {
	for _, tc := range []struct {
		name     string
		min, max *int
		lo, hi   int64
	}{
		{"bounded", ptr(1), ptr(6), 1, 6},
		{"swapped", ptr(6), ptr(1), 1, 6},
		{"single", ptr(4), ptr(4), 4, 4},
		{"default", nil, nil, math.MinInt32, math.MaxInt32},
		{"above int32", ptr(math.MaxInt32 + 1), ptr(math.MaxInt32 + 10), math.MaxInt32 + 1, math.MaxInt32 + 10},
		{"full range", ptr(math.MinInt), ptr(math.MaxInt), math.MinInt64, math.MaxInt64},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for range 200 {
				raw, err := generateNumber(gdt.NumberRestrictions{Min: tc.min, Max: tc.max})
				if err != nil {
					t.Fatal(err)
				}
				var v int64
				if err := json.Unmarshal(raw, &v); err != nil {
					t.Fatal(err)
				}
				if v < tc.lo || v > tc.hi {
					t.Fatalf("value %d outside [%d, %d]", v, tc.lo, tc.hi)
				}
			}
		})
	}
}

func TestGenerateString(t *testing.T) {
	for _, tc := range []struct {
		name     string
		min, max *int
		lo, hi   int
	}{
		{"bounded", ptr(3), ptr(5), 3, 5},
		{"swapped", ptr(5), ptr(3), 3, 5},
		{"negative max", nil, ptr(-5), 0, 0},
		{"negative min", ptr(-5), ptr(2), 0, 2},
		{"min only", ptr(10), nil, 10, 10 + 254},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for range 200 {
				raw, err := generateString(gdt.StringRestrictions{MinLength: tc.min, MaxLength: tc.max})
				if err != nil {
					t.Fatal(err)
				}
				var s string
				if err := json.Unmarshal(raw, &s); err != nil {
					t.Fatal(err)
				}
				if n := utf8.RuneCountInString(s); n < tc.lo || n > tc.hi {
					t.Fatalf("len %d outside [%d, %d]", n, tc.lo, tc.hi)
				}
			}
		})
	}

	raw, err := generateString(gdt.StringRestrictions{RegexExpression: ptr(`^[A-Z]{3}-[0-9]{2}$`)})
	if err != nil {
		t.Fatal(err)
	}
	var s string
	_ = json.Unmarshal(raw, &s)
	if !regexp.MustCompile(`^[A-Z]{3}-[0-9]{2}$`).MatchString(s) {
		t.Errorf("generated %q does not match regex", s)
	}
}

func TestBaseTypeGeneration(t *testing.T) {
	_, err := GenerateKnownedType(nil, nil, GenerationInfos{TypeName: "number", Restrictions: json.RawMessage(`{"min":1,"max":3}`)})
	if err != nil {
		t.Errorf("number: %v", err)
	}
	if _, err := GenerateKnownedType(nil, nil, GenerationInfos{TypeName: "potato", Restrictions: json.RawMessage(`{}`)}); err == nil {
		t.Error("unknown type should fail")
	}
	if _, err := GenerateKnownedType(nil, nil, GenerationInfos{TypeName: "enum", Restrictions: json.RawMessage(`{}`)}); err == nil {
		t.Error("enum without tags should fail before touching the db")
	}
	// "dictionary" must be routed to the dictionary generator (still TODO).
	if _, err := GenerateKnownedType(nil, nil, GenerationInfos{TypeName: "dictionary", Restrictions: json.RawMessage(`{}`)}); err == nil || err.Error() != "TODO" {
		t.Errorf("dictionary should reach generateDictionary, got %v", err)
	}
}

func TestScriptingGeneration(t *testing.T) {
	script := `
def generate_type(restrictions, static_values):
    return static_values["str"] * 2 + restrictions["bonus"]
`
	raw, err := GenerateKnownedType(nil,
		gdt.ComponentMap{"str": json.RawMessage(`7`)},
		GenerationInfos{Script: &script, TypeName: "number", Restrictions: json.RawMessage(`{"bonus":1}`)},
	)
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "15" {
		t.Errorf("got %s, want 15", raw)
	}
}
