package validate

import (
	"database/sql"
	"fmt"
	"regexp"
	"unicode/utf8"

	gdt "github.com/D-Gine/Gestionnaire2Type"
)

func VerifyString(db *sql.DB, componentList *gdt.ComponentMap, value string, restrictions gdt.StringRestrictions) error {
	length := utf8.RuneCountInString(value)
	if restrictions.MaxLength != nil && length > *restrictions.MaxLength {
		return fmt.Errorf("too long, max lenght is %d, current len is %d", *restrictions.MaxLength, length)
	}
	if restrictions.MinLength != nil && length < *restrictions.MinLength {
		return fmt.Errorf("too short, min lenght is %d, current len is %d", *restrictions.MinLength, length)
	}
	if restrictions.RegexExpression != nil {
		match, err := regexp.MatchString(*restrictions.RegexExpression, value)
		if err != nil {
			return fmt.Errorf("invalid regex expression %q: %w", *restrictions.RegexExpression, err)
		}
		if !match {
			return fmt.Errorf("regex expression %q not matching", *restrictions.RegexExpression)
		}
	}
	return nil
}
