package validate

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"

	gdt "github.com/D-Gine/Gestionnaire2Type"
	"github.com/lib/pq"
)

func VerifyDictionary(db *sql.DB, componentList *gdt.ComponentMap, value map[string]json.RawMessage, restrictions gdt.DictionaryRestrictions) error {
	if restrictions.Interface == nil {
		return nil
	}
	iface := *restrictions.Interface

	typeSet := make(map[string]struct{})
	for _, typeId := range iface {
		typeSet[typeId] = struct{}{}
	}
	types := slices.Collect(maps.Keys(typeSet))
	rows, err := db.Query(`
		SELECT * FROM get_verifications($1)`, pq.Array(types))
	if err != nil {
		return fmt.Errorf("error fetching interfaces types, please recheck validity of you interface, %w", err)
	}
	defer func() { _ = rows.Close() }()
	verificationMap := make(map[string]TypeVerifications)
	for rows.Next() {
		var verification TypeVerifications
		var baseName sql.NullString
		if err := rows.Scan(
			&verification.TypeId,
			&baseName,
			&verification.Restrictions,
			pq.Array(&verification.Scripts),
		); err != nil {
			return fmt.Errorf("error scaning db row, %w", err)
		}
		if !baseName.Valid {
			return fmt.Errorf("unknown type %q in interface", verification.TypeId)
		}
		verification.BaseName = baseName.String
		verificationMap[verification.TypeId] = verification
	}
	if err = rows.Err(); err != nil {
		return fmt.Errorf("error reading db rows, %w", err)
	}

	missing := maps.Clone(iface)
	for _, key := range slices.Sorted(maps.Keys(value)) {
		typeId, ok := iface[key]
		if !ok {
			return fmt.Errorf(`unknow key "%s"`, key)
		}
		verification, ok := verificationMap[typeId]
		if !ok {
			return fmt.Errorf(`unknown type %q for key "%s"`, typeId, key)
		}
		if err := VerifyType(db, componentList, value[key], verification); err != nil {
			return fmt.Errorf(`error at key "%s", %w`, key, err)
		}
		delete(missing, key)
	}
	if len(missing) > 0 {
		return fmt.Errorf(
			"missing keys: %s",
			strings.Join(slices.Sorted(maps.Keys(missing)), ", "),
		)
	}
	return nil
}
