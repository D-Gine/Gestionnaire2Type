package validate

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	gdt "github.com/D-Gine/Gestionnaire2Type"
	"github.com/D-Gine/Gestionnaire2Type/scripting"
)

func runValidateScript(db *sql.DB, value json.RawMessage, script string, restrictions json.RawMessage, componentList *gdt.ComponentMap) error {
	rawResult, err := scripting.ScriptEngine.Eval(
		script,
		"validate_type",
		[]any{value, restrictions, *componentList},
		nil,
	)
	if err != nil {
		return fmt.Errorf("script execution error: %w", err)
	}
	stringError, err := scripting.AsOptionalString(rawResult)
	if err != nil {
		return fmt.Errorf("script return type error: %w", err)
	}
	if stringError == nil {
		return nil
	}
	return errors.New(*stringError)
}
