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
