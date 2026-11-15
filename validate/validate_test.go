package validate

import (
	"encoding/json"
	"strings"
	"testing"

	gdt "github.com/D-Gine/Gestionnaire2Type"
)

func ptr[T any](v T) *T { return &v }

func TestVerifyNumber(t *testing.T) {
	r := gdt.NumberRestrictions{Min: ptr(1), Max: ptr(20)}
	for _, tc := range []struct {
		value int
		ok    bool
	}{{1, true}, {20, true}, {10, true}, {0, false}, {21, false}} {
		err := VerifyNumber(nil, nil, tc.value, r)
		if (err == nil) != tc.ok {
			t.Errorf("VerifyNumber(%d) err=%v, want ok=%v", tc.value, err, tc.ok)
		}
	}
}

func TestVerifyString(t *testing.T) {
	r := gdt.StringRestrictions{MinLength: ptr(2), MaxLength: ptr(5), RegexExpression: ptr(`^[a-zéè]+$`)}
	for _, tc := range []struct {
		value string
		ok    bool
	}{
		{"ab", true},
		{"élève", true},
		{"a", false},
		{"abcdef", false},
		{"AB", false},
	} {
		err := VerifyString(nil, nil, tc.value, r)
		if (err == nil) != tc.ok {
			t.Errorf("VerifyString(%q) err=%v, want ok=%v", tc.value, err, tc.ok)
		}
		if err != nil && strings.Contains(err.Error(), "%!") {
			t.Errorf("VerifyString(%q) malformed error message: %v", tc.value, err)
		}
	}

	err := VerifyString(nil, nil, "x", gdt.StringRestrictions{RegexExpression: ptr(`(`)})
	if err == nil || !strings.Contains(err.Error(), "invalid regex") {
		t.Errorf("invalid regex should be reported, got %v", err)
	}
}

func TestValidateBaseType(t *testing.T) {
	for _, tc := range []struct {
		name, base, value, restrictions string
		ok                              bool
	}{
		{"valid number", "number", `5`, `{"min":1,"max":10}`, true},
		{"number too high", "number", `50`, `{"max":10}`, false},
		{"number is string", "number", `"5"`, `{}`, false},
		{"float is not number", "number", `5.5`, `{}`, false},
		{"null not optional", "number", `null`, `{}`, false},
		{"null optional", "number", `null`, `{"optional":true}`, true},
		{"valid string", "string", `"hello"`, `{"max_length":10}`, true},
		{"bad restriction field type", "string", `"hello"`, `{"max_length":"10"}`, false},
		{"unknown base type", "potato", `1`, `{}`, false},
		{"empty base type", "", `1`, `{}`, false},
		{"dictionary without interface", "dictionary", `{"a":1}`, `{}`, true},
		{"dictionary wrong value", "dictionary", `[1]`, `{}`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateBaseType(nil, &gdt.ComponentMap{}, json.RawMessage(tc.value), tc.base, json.RawMessage(tc.restrictions))
			if (err == nil) != tc.ok {
				t.Errorf("err=%v, want ok=%v", err, tc.ok)
			}
		})
	}
}

const validateScript = `
def validate_type(value, restrictions, components):
    if value % 2 != 0:
        return "must be even"
    if components.get("other") != None and value <= components["other"]:
        return "must be greater than other"
    return None
`

func TestVerifyTypeWithScript(t *testing.T) {
	verif := TypeVerifications{
		TypeId: "t",
		BaseVerifications: BaseVerifications{
			BaseName:     "number",
			Restrictions: json.RawMessage(`{"max":100}`),
			Scripts:      []string{validateScript},
		},
	}
	components := gdt.ComponentMap{"other": json.RawMessage(`10`)}
	for _, tc := range []struct {
		value   string
		wantErr string
	}{
		{`12`, ""},
		{`13`, "must be even"},
		{`8`, "must be greater than other"},
		{`102`, "too hight"},
	} {
		err := VerifyType(nil, &components, json.RawMessage(tc.value), verif)
		switch {
		case tc.wantErr == "" && err != nil:
			t.Errorf("value %s: unexpected error %v", tc.value, err)
		case tc.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tc.wantErr)):
			t.Errorf("value %s: err=%v, want containing %q", tc.value, err, tc.wantErr)
		}
	}
}

func TestVerifyTypeBrokenScript(t *testing.T) {
	verif := ComponentVerification{
		ComponentId: "c",
		BaseVerifications: BaseVerifications{
			BaseName:     "string",
			Restrictions: json.RawMessage(`{}`),
			Scripts:      []string{"def nope(:\n"},
		},
	}
	if err := VerifyType(nil, &gdt.ComponentMap{}, json.RawMessage(`"x"`), verif); err == nil {
		t.Error("expected script error")
	}
}
