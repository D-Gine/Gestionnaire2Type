package validate

import (
	"database/sql"
	"encoding/json"
	"fmt"

	gdt "github.com/D-Gine/Gestionnaire2Type"
	"github.com/lib/pq"
)

func VerifyList(db *sql.DB, componentList *gdt.ComponentMap, value []json.RawMessage, restrictions gdt.ListRestrictions) error {
	var verification TypeVerifications
	err := db.QueryRow(`
		SELECT * FROM get_verifications($1)`, pq.Array([]string{restrictions.Type})).Scan(
		&verification.TypeId,
		&verification.BaseName,
		&verification.Restrictions,
		pq.Array(&verification.Scripts),
	)
	if err != nil {
		return fmt.Errorf("error fetching for base type, please check your restriction 'type' field, %w", err)
	}
	if restrictions.MaxLength != nil && len(value) > *restrictions.MaxLength {
		return fmt.Errorf("too long, max lenght is %d, current len is %d", *restrictions.MaxLength, len(value))
	}
	if restrictions.MinLength != nil && len(value) < *restrictions.MinLength {
		return fmt.Errorf("too short, min lenght is %d, current len is %d", *restrictions.MinLength, len(value))
	}
	for i, elem := range value {
		if err = VerifyType(db, componentList, elem, verification); err != nil {
			return fmt.Errorf("error at element : %d, %w", i, err)
		}
	}
	return nil
}
