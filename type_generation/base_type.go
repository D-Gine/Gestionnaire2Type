package typegeneration

import (
	"database/sql"
	"encoding/json"
	"fmt"

	gdt "github.com/D-Gine/Gestionnaire2Type"
)

func BaseTypeGeneration(db *sql.DB, staticValues gdt.ComponentMap, infos GenerationInfos) (json.RawMessage, error) {
	switch infos.TypeName {
	case "number":
		r, err := gdt.ParseRestriction[gdt.NumberRestrictions](infos.TypeName, infos.Restrictions)
		if err != nil {
			return nil, err
		}
		return generateNumber(r)
	case "string":
		r, err := gdt.ParseRestriction[gdt.StringRestrictions](infos.TypeName, infos.Restrictions)
		if err != nil {
			return nil, err
		}
		return generateString(r)
	case "enum":
		r, err := gdt.ParseRestriction[gdt.EnumRestrictions](infos.TypeName, infos.Restrictions)
		if err != nil {
			return nil, err
		}
		return generateEnum(db, r)
	case "list":
		r, err := gdt.ParseRestriction[gdt.ListRestrictions](infos.TypeName, infos.Restrictions)
		if err != nil {
			return nil, err
		}
		return generateList(r, staticValues)
	case "dictionary":
		r, err := gdt.ParseRestriction[gdt.DictionaryRestrictions](infos.TypeName, infos.Restrictions)
		if err != nil {
			return nil, err
		}
		return generateDictionary(r, staticValues)
	}
	return nil, fmt.Errorf("unknown type %q", infos.TypeName)
}
