package validate

import (
	"database/sql"
	"errors"

	gdt "github.com/D-Gine/Gestionnaire2Type"
	"github.com/lib/pq"
)

func VerifyEnum(db *sql.DB, componentList *gdt.ComponentMap, value string, restrictions gdt.EnumRestrictions) error {
	if restrictions.Tags != nil {
		rows := db.QueryRow(`
			SELECT COUNT(*) FROM link_enums_tags WHERE enum_id = $1 AND tag_id = ANY($2::uuid[])`, value, pq.Array(restrictions.Tags))
		var count int
		rows.Scan(&count)
		if count != len(restrictions.Tags) {
			return errors.New("does not respect tag restrictions")
		}
	}
	return nil
}
