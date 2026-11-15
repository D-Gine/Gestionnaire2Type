package gdt

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func GetAllowedFields(v any) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return ""
	}

	var fields []string
	for field := range t.Fields() {
		field := field
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			// Extract tag name before any options like omitempty
			tagName, _, _ := strings.Cut(jsonTag, ",")
			fields = append(fields, fmt.Sprintf("- %s: %s (optional)", tagName, field.Type.Kind()))
		}
	}
	return strings.Join(fields, "\n  ")
}

func ParseRestriction[T interface {
	GetBaseRestrictions() BaseRestrictions
}](
	name string,
	restrictionsJson json.RawMessage,
) (T, error) {

	var restrictions T
	if err := json.Unmarshal(restrictionsJson, &restrictions); err != nil {
		return restrictions, fmt.Errorf(`invalid restrictions for "%s": %w
			Allowed fields:
			%s`, name, err, GetAllowedFields(restrictions))
	}
	return restrictions, nil
}
