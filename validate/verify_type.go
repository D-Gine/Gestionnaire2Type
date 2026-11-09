package validate

import (
	"database/sql"
	"encoding/json"
	"fmt"

	gdt "github.com/D-Gine/Gestionnaire2Type"
	"github.com/lib/pq"
)

func VerifyType(db *sql.DB, componentList *gdt.ComponentMap, valueJson json.RawMessage, baseVerification Verification) error {
	verification := baseVerification.GetBase()
	if err := ValidateBaseType(db, componentList, valueJson, verification.BaseName, verification.Restrictions); err != nil {
		return fmt.Errorf("error validation base type: %w", err)
	}

	for i := 0; i < len(verification.Scripts); i++ {
		if err := runValidateScript(
			db,
			valueJson,
			verification.Scripts[i],
			verification.Restrictions,
			componentList,
		); err != nil {
			return err
		}
	}
	return nil
}

func ValidateRowTypes(db *sql.DB, toVerify *map[string]json.RawMessage, components *gdt.ComponentMap, rows *sql.Rows) (errorsMap map[string]string, count int, err error) {
	errorsMap = make(map[string]string)
	for count = 0; rows.Next(); count++ {
		var result TypeVerifications
		if err = rows.Scan(&result.TypeId, &result.BaseName, pq.Array(&result.Scripts), &result.Restrictions); err != nil {
			return errorsMap, count, gdt.DatabaseError(err)
		}
		if typeErr := VerifyType(db, components, (*toVerify)[result.TypeId], result); typeErr != nil {
			errorsMap[result.TypeId] = typeErr.Error()
		}
	}
	return errorsMap, count, gdt.DatabaseError(rows.Err())
}

func ValidateRowComponents(db *sql.DB, components *gdt.ComponentMap, rows *sql.Rows) (map[string]string, int, error) {
	return ValidateRowComponentsWithValues(db, components, components, rows)
}

func ValidateRowComponentsWithValues(db *sql.DB, toVerify *gdt.ComponentMap, components *gdt.ComponentMap, rows *sql.Rows) (errorsMap map[string]string, count int, err error) {
	errorsMap = make(map[string]string)
	for count = 0; rows.Next(); count++ {
		var result ComponentVerification
		if err = rows.Scan(&result.ComponentId, &result.BaseName, pq.Array(&result.Scripts), &result.Restrictions); err != nil {
			return errorsMap, count, gdt.DatabaseError(err)
		}
		if typeErr := VerifyType(db, components, (*toVerify)[result.ComponentId], result); typeErr != nil {
			errorsMap[result.ComponentId] = typeErr.Error()
		}
	}
	return errorsMap, count, gdt.DatabaseError(rows.Err())
}
