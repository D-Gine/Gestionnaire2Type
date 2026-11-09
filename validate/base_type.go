package validate

import (
	"database/sql"
	"encoding/json"
	"fmt"

	gdt "github.com/D-Gine/Gestionnaire2Type"
)

func InitBaseTypeVerification[Value any, Restrictions interface {
	GetBaseRestrictions() gdt.BaseRestrictions
}](
	db *sql.DB,
	nodeList *gdt.ComponentMap,
	name string,
	valueJson json.RawMessage,
	restrictionsJson json.RawMessage,
	verifyFunc func(*sql.DB, *gdt.ComponentMap, Value, Restrictions) error) error {

	var value *Value

	if err := json.Unmarshal(valueJson, &value); err != nil {
		return fmt.Errorf(`invalid value for "%s": %w`, name, err)
	}
	restrictions, err := gdt.ParseRestriction[Restrictions](name, restrictionsJson)
	if err != nil {
		return err
	}
	if value == nil {
		opt := restrictions.GetBaseRestrictions().OptionalBool
		if opt != nil && *opt {
			return nil
		} else {
			return fmt.Errorf("invalid %s: can't be null", name)
		}
	}

	if err := verifyFunc(db, nodeList, *value, restrictions); err != nil {
		return fmt.Errorf("invalid %s: %w", name, err)
	}
	return nil
}

func ValidateBaseType(db *sql.DB, nodeList *gdt.ComponentMap, valueJson json.RawMessage, name string, restrictionsJson json.RawMessage) error {
	switch name {
	case "number":
		return InitBaseTypeVerification(db, nodeList, name, valueJson, restrictionsJson, VerifyNumber)
	case "string":
		return InitBaseTypeVerification(db, nodeList, name, valueJson, restrictionsJson, VerifyString)
	case "enum":
		return InitBaseTypeVerification(db, nodeList, name, valueJson, restrictionsJson, VerifyEnum)
	case "list":
		return InitBaseTypeVerification(db, nodeList, name, valueJson, restrictionsJson, VerifyList)
	case "dictionary":
		return InitBaseTypeVerification(db, nodeList, name, valueJson, restrictionsJson, VerifyDictionary)
	}

	return fmt.Errorf("unknown base type %q", name)
}
