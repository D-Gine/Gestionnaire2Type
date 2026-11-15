package validate

import (
	"database/sql"
	"fmt"

	gdt "github.com/D-Gine/Gestionnaire2Type"
)

func VerifyNumber(db *sql.DB, componentList *gdt.ComponentMap, value int, restrictions gdt.NumberRestrictions) error {
	if restrictions.Max != nil && value > *restrictions.Max {
		return fmt.Errorf("too hight, max value is %d, current value is %d", *restrictions.Max, value)
	}
	if restrictions.Min != nil && value < *restrictions.Min {
		return fmt.Errorf("too low, min value is %d, current value is %d", *restrictions.Min, value)
	}
	return nil
}
