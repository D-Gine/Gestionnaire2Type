package typegeneration

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"slices"
	"testing"

	gdt "github.com/D-Gine/Gestionnaire2Type"
	"github.com/lib/pq"
)

func TestIntegrationGenerateEnum(t *testing.T) {
	dsn := os.Getenv("GDT_TEST_DSN")
	if dsn == "" {
		t.Skip("GDT_TEST_DSN not set")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	want := []string{"650d303e-fa01-499d-ac3b-319d7317496e", "d82a1464-7b0c-4d2a-83e7-77b659164fec"}
	for range 20 {
		raw, err := GenerateKnownedType(db, nil, GenerationInfos{
			TypeName:     "enum",
			Restrictions: json.RawMessage(`{"tags":["080bae99-829d-4fc9-91c7-6b43b06b24b0","110ee2a7-d94a-416f-b66d-4d26927fc268"]}`),
		})
		if err != nil {
			t.Fatal(err)
		}
		var id string
		if err := json.Unmarshal(raw, &id); err != nil || !slices.Contains(want, id) {
			t.Fatalf("got %s (%v), want one of %v", raw, err, want)
		}
	}

	if _, err := GenerateKnownedType(db, nil, GenerationInfos{
		TypeName:     "enum",
		Restrictions: json.RawMessage(`{"tags":["0e8ef112-a906-4c50-8c68-26501212eaa0","dcc6c81e-d887-4df0-ab1f-76d0df97e63d"]}`),
	}); err == nil {
		t.Error("expected no enum found")
	}

	if _, err := GenerateKnownedType(db, nil, GenerationInfos{
		TypeName:     "enum",
		Restrictions: json.RawMessage(`{"tags":["not-a-uuid"]}`),
	}); err == nil {
		t.Error("expected db error for malformed uuid")
	}
}

func TestIntegrationGenerateRowTypes(t *testing.T) {
	dsn := os.Getenv("GDT_TEST_DSN")
	if dsn == "" {
		t.Skip("GDT_TEST_DSN not set")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	const typeLevel = "7786f8ca-40c0-4b22-8319-8c38da55f51c"
	rows, err := db.Query(`
		SELECT to_generate.id, tg.script, tg.type_name, tg.restrictions
		FROM unnest($1::uuid[]) AS to_generate(id)
		JOIN type_generation(to_generate.id) tg ON 1 = 1
	`, pq.Array([]string{typeLevel}))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = rows.Close() }()
	generated, err := GenerateRowTypes(db, nil, rows)
	if err != nil {
		t.Fatal(err)
	}
	var v int
	if err := json.Unmarshal(generated[typeLevel], &v); err != nil || v < 1 || v > 20 {
		t.Errorf("got %s (%v), want a number in [1, 20]", generated[typeLevel], err)
	}

	_, err = GenerateKnownedType(db, nil, GenerationInfos{TypeName: "enum", Restrictions: json.RawMessage(`{"tags":["not-a-uuid"]}`)})
	if !errors.Is(err, gdt.ErrDatabase) {
		t.Errorf("want gdt.ErrDatabase, got %v", err)
	}
}
