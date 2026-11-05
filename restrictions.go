package gdt

import (
	"encoding/json"
)

type ComponentMap = map[string]json.RawMessage

type BaseRestrictions struct {
	OptionalBool *bool `json:"optional,omitempty"`
}

func (b BaseRestrictions) GetBaseRestrictions() BaseRestrictions {
	return b
}

type NumberRestrictions struct {
	BaseRestrictions
	Min *int `json:"min,omitempty"`
	Max *int `json:"max,omitempty"`
}

type StringRestrictions struct {
	BaseRestrictions
	MinLength       *int    `json:"min_length,omitempty"`
	MaxLength       *int    `json:"max_length,omitempty"`
	RegexExpression *string `json:"regex_expression,omitempty"`
}

type EnumRestrictions struct {
	BaseRestrictions
	Tags []string `json:"tags,omitempty"`
}

type ListRestrictions struct {
	BaseRestrictions
	Type      string `json:"type"`
	MinLength *int   `json:"min_length,omitempty"`
	MaxLength *int   `json:"max_length,omitempty"`
}

type DictionaryRestrictions struct {
	BaseRestrictions
	Interface *map[string]string `json:"interface,omitempty"`
}
