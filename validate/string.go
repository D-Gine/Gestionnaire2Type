package validate

import (
	"database/sql"
	"fmt"
	"regexp"

	gdt "github.com/D-Gine/Gestionnaire2Type"
)

func VerifyString(db *sql.DB, componentList *gdt.ComponentMap, value string, restrictions gdt.StringRestrictions) error {
	if restrictions.MaxLength != nil && len(value) > *restrictions.MaxLength {
		return fmt.Errorf("too long, max lenght is %d, current len is %d", *restrictions.MaxLength, len(value))
	}
	if restrictions.MinLength != nil && len(value) < *restrictions.MinLength {
		return fmt.Errorf("too short, min lenght is %d, current len is %d", *restrictions.MinLength, len(value))
	}
	if restrictions.RegexExpression != nil {
		match, err := regexp.MatchString(*restrictions.RegexExpression, value)
		if !match {
			return fmt.Errorf("regex expression not matching, %w", err)
		}
	}
	return nil
}
