package validate

import (
	"database/sql"
	"encoding/json"
	"errors"
	"maps"
	"slices"

	gdt "github.com/D-Gine/Gestionnaire2Type"
	"github.com/lib/pq"
)

// ValidateNodesTypes validates the values sent for some nodes of a ruleset.
//
// nodeValues maps node_id -> value, nodeIds lists the nodes that must be
// validated (usually the keys of nodeValues).
//
// The returned map is node_id -> error message (empty when everything is
// valid). err is set for database errors or when some of nodeIds do not exist
// in the ruleset.
func ValidateNodesTypes(db *sql.DB, nodeValues map[string]json.RawMessage, nodeIds []string, rulesetId string) (map[string]string, error) {
	nodeKeys := slices.Collect(maps.Keys(nodeValues))
	nodeRows, err := db.Query(`
		SELECT id, template_id FROM nodes WHERE id = ANY($1)
	`, pq.Array(nodeKeys))
	if err != nil {
		return nil, gdt.DatabaseError(err)
	}
	defer func() { _ = nodeRows.Close() }()
	componentList := make(gdt.ComponentMap)
	nodeRecoverMap := make(map[string]string)
	for nodeRows.Next() {
		var nodeId, compId string
		if err := nodeRows.Scan(&nodeId, &compId); err != nil {
			return nil, gdt.DatabaseError(err)
		}
		if val, exists := nodeValues[nodeId]; exists {
			componentList[compId] = val
			nodeRecoverMap[compId] = nodeId
		}
	}
	if err := nodeRows.Err(); err != nil {
		return nil, gdt.DatabaseError(err)
	}

	verifyRows, err := db.Query(`
		SELECT * FROM get_nodes_verifications($1, $2::uuid)
	`, pq.Array(nodeIds), rulesetId)
	if err != nil {
		return nil, gdt.DatabaseError(err)
	}
	defer func() { _ = verifyRows.Close() }()

	errorMap, count, err := ValidateRowComponents(db, &componentList, verifyRows)
	nodeErrorMap := make(map[string]string)
	for id, v := range errorMap {
		nodeErrorMap[nodeRecoverMap[id]] = v
	}
	if err != nil {
		return nodeErrorMap, err
	}

	if count != len(nodeIds) {
		return nil, errors.New("some sent nodes are invalid, cannot compute")
	}

	return nodeErrorMap, nil
}
