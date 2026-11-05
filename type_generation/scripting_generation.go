package typegeneration

import (
	"encoding/json"
	"fmt"

	gdt "github.com/D-Gine/Gestionnaire2Type"
	"github.com/D-Gine/Gestionnaire2Type/scripting"
)

func ScriptingGeneration(staticValues gdt.ComponentMap, infos GenerationInfos) (json.RawMessage, error) {
	value, err := scripting.ScriptEngine.Eval(
		*infos.Script,
		"generate_type",
		[]any{infos.Restrictions, staticValues},
		nil,
	)
	if err != nil {
		return nil, err
	}
	rawValue, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("unable to parse return value for '%s': %w", infos.TypeName, err)
	}
	return rawValue, nil
}
