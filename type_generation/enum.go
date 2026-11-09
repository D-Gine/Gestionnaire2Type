package typegeneration

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	gdt "github.com/D-Gine/Gestionnaire2Type"
	"github.com/lib/pq"
)

func generateEnum(db *sql.DB, restriction gdt.EnumRestrictions) (json.RawMessage, error) {
	var enum string

	if len(restriction.Tags) == 0 {
		//TODO
		return nil, errors.New("cannot generate enum with no tag")
	}
	err := db.QueryRow(`
			SELECT enum_id
			FROM link_enums_tags
			WHERE tag_id = ANY($1::UUID[])
			GROUP BY enum_id
			HAVING COUNT(DISTINCT tag_id) = array_length($1, 1)
			ORDER BY RANDOM()
			LIMIT 1
		`, pq.Array(restriction.Tags)).Scan(&enum)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no enum found for tag requirement %q", restriction.Tags)
		}
		return nil, gdt.DatabaseError(fmt.Errorf("error fetching enum: %w", err))
	}
	return json.Marshal(enum)
}
