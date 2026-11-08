# Gestionnaire2Type

Shared Go package for D-Gine dynamic types: **validate** values against types
stored in the database and **generate** values for them. Used by both the API
and the EngineWorker so they produce the same results without calling each
other.

```
github.com/D-Gine/Gestionnaire2Type                 // gdt: restriction structs, ComponentMap
github.com/D-Gine/Gestionnaire2Type/validate        // validation
github.com/D-Gine/Gestionnaire2Type/type_generation // generation (package typegeneration)
github.com/D-Gine/Gestionnaire2Type/scripting       // Starlark script engine
```

## Install

The repo is private:

```sh
go env -w GOPRIVATE=github.com/D-Gine/*
go get github.com/D-Gine/Gestionnaire2Type@latest
```

To use it as a git submodule (like the API, in `api/src/gdt`), point the module
to the submodule folder in the consumer's `go.mod`:

```
replace github.com/D-Gine/Gestionnaire2Type => ./src/gdt
```

## Usage

Every function takes the `*sql.DB` explicitly (no global, no gin).

```go
import (
	"github.com/D-Gine/Gestionnaire2Type/validate"
	typegeneration "github.com/D-Gine/Gestionnaire2Type/type_generation"
)

// Validate values sent for nodes of a ruleset.
// errs: node_id -> error message (empty = valid). err: db error / unknown nodes.
errs, err := validate.ValidateNodesTypes(db, nodeValues, nodeIds, rulesetId)

// Validate every static component stored for an entity (character...).
// errs: template_id -> error message.
errs, err := validate.ValidateEntity(db, entityId)

// Lower level: validate rows from get_nodes_verifications / get_entity_verifications.
errs, count, err := validate.ValidateRowComponents(db, &components, rows)

// Generate a value for a type (base type or with a generation script).
value, err := typegeneration.GenerateKnownedType(db, staticValues, typegeneration.GenerationInfos{
	TypeName:     "number",
	Restrictions: json.RawMessage(`{"min": 1, "max": 20}`),
})
```

### Scripts

Scripts are Starlark (`while` loops allowed). Each evaluation is bounded by
`scripting.DefaultMaxSteps` so an infinite loop returns an error instead of
hanging. Expected functions:

- `validate_type(value, restrictions, component_list)` → `None` / `""` if valid, an error string otherwise
- `generate_type(restrictions, static_values)` → the generated value

## Tests

```sh
go test ./...
```

DB tests run only when `GDT_TEST_DSN` is set, against a database filled with
`Database/scripts/1_fill_test`:

```sh
GDT_TEST_DSN="host=localhost port=5432 user=... password=... dbname=... sslmode=disable" go test ./...
```
