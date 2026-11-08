package validate

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	gdt "github.com/D-Gine/Gestionnaire2Type"
)

func GetEntityComponents(db *sql.DB, entityId string) (gdt.ComponentMap, error) {
	var rawJSON []byte
	err := db.QueryRow(`
		SELECT jsonb_object_agg(template_id, value)
		FROM static_components
		WHERE entity_id = $1
	`, entityId).Scan(&rawJSON)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, gdt.DatabaseError(err)
	}
	componentList := make(gdt.ComponentMap)
	if len(rawJSON) > 0 {
		if err := json.Unmarshal(rawJSON, &componentList); err != nil {
			return nil, fmt.Errorf("failed to parse components: %w", err)
		}
	}
	return componentList, nil
}
func ValidateEntity(db *sql.DB, entityId string) (map[string]string, error) {
	componentList, err := GetEntityComponents(db, entityId)
	if err != nil {
		return nil, err
	}

	verifyRows, err := db.Query(`
		SELECT * FROM get_entity_verifications($1)
	`, entityId)
	if err != nil {
		return nil, gdt.DatabaseError(err)
	}
	defer func() { _ = verifyRows.Close() }()

	errorMap, _, err := ValidateRowComponents(db, &componentList, verifyRows)
	if err != nil {
		return nil, err
	}
	return errorMap, nil
}
