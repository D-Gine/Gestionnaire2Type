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
