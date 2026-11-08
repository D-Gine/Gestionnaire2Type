package typegeneration

import (
	"database/sql"
	"encoding/json"

	gdt "github.com/D-Gine/Gestionnaire2Type"
)

type GenerationInfos struct {
	Script       *string         `json:"script,omitempty"`
	TypeName     string          `json:"type_name"`
	Restrictions json.RawMessage `json:"restrictions"`
}

func GenerateKnownedType(db *sql.DB, staticValues gdt.ComponentMap, infos GenerationInfos) (json.RawMessage, error) {
	if infos.Script != nil {
		return ScriptingGeneration(staticValues, infos)
	}
	return BaseTypeGeneration(db, staticValues, infos)
}
func GenerateRowTypes(db *sql.DB, components gdt.ComponentMap, rows *sql.Rows) (map[string]json.RawMessage, error) {
	generated := make(map[string]json.RawMessage)
	for rows.Next() {
		var infos GenerationInfos
		var id string
		if err := rows.Scan(&id, &infos.Script, &infos.TypeName, &infos.Restrictions); err != nil {
			return nil, gdt.DatabaseError(err)
		}
		value, err := GenerateKnownedType(db, components, infos)
		if err != nil {
			return nil, err
		}
		generated[id] = value
	}
	if err := rows.Err(); err != nil {
		return nil, gdt.DatabaseError(err)
	}
	return generated, nil
}
