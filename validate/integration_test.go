package validate

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"

	gdt "github.com/D-Gine/Gestionnaire2Type"
	"github.com/lib/pq"
)

const (
	testRuleset = "0e3196a0-f8d8-444f-a6b6-d5f6f81f6333"

	tagClass = "64f73b11-5699-46a2-9f63-4db41f18e472"

	enumClassA     = "a87fa8f3-d56a-49ec-8bd2-d62b7f07c367"
	enumBerserker  = "f552b243-81fb-4901-9b09-a140de4660f4"
	typeClassEntry = "b7c1d4e2-3a5f-4e8b-9c21-6f0a7d5e1b34"

	nodeClasses   = "5add1a83-bd43-4385-b04a-c3797138d684"
	nodeLevel     = "3e4ea0d2-1270-4816-bd30-e9f2108db35d"
	nodeHp        = "74460e93-5e8e-4fef-bbb7-0ca529bdef6d"
	nodeBaseStats = "4fbb7e27-8e26-47dd-bea1-f687cc4a98db"

	entityTest     = "c1f8b452-91fd-4c0a-bf93-48a6d03d7846"
	entityFreebird = "41371b0d-33c9-4a97-b902-c0a1b0108348"
)

func testDB(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv("GDT_TEST_DSN")
	if dsn == "" {
		t.Skip("GDT_TEST_DSN not set")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func TestIntegrationVerifyEnum(t *testing.T) {
	db := testDB(t)
	r := gdt.EnumRestrictions{Tags: []string{tagClass}}
	if err := VerifyEnum(db, nil, enumClassA, r); err != nil {
		t.Errorf("tagged enum: %v", err)
	}
	if err := VerifyEnum(db, nil, enumBerserker, r); err == nil {
		t.Error("untagged enum should fail")
	}
	if err := VerifyEnum(db, nil, "not-a-uuid", r); err == nil {
		t.Error("invalid uuid should fail")
	}
}

func TestIntegrationVerifyList(t *testing.T) {
	db := testDB(t)
	r := gdt.ListRestrictions{Type: typeClassEntry, MinLength: ptr(1), MaxLength: ptr(2)}
	entry := func(enum string, level int) json.RawMessage {
		b, _ := json.Marshal(map[string]any{"class": enum, "level": level})
		return b
	}
	for _, tc := range []struct {
		name  string
		value []json.RawMessage
		ok    bool
	}{
		{"valid", []json.RawMessage{entry(enumClassA, 3)}, true},
		{"empty", nil, false},
		{"too long", []json.RawMessage{entry(enumClassA, 1), entry(enumClassA, 1), entry(enumClassA, 1)}, false},
		{"bad enum", []json.RawMessage{entry(enumBerserker, 3)}, false},
		{"bad level", []json.RawMessage{entry(enumClassA, 99)}, false},
		{"missing key", []json.RawMessage{json.RawMessage(`{"class":"` + enumClassA + `"}`)}, false},
		{"extra key", []json.RawMessage{json.RawMessage(`{"class":"` + enumClassA + `","level":1,"x":1}`)}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := VerifyList(db, &gdt.ComponentMap{}, tc.value, r)
			if (err == nil) != tc.ok {
				t.Errorf("err=%v, want ok=%v", err, tc.ok)
			}
		})
	}
}

func TestIntegrationVerifyDictionaryDoesNotMutateRestrictions(t *testing.T) {
	db := testDB(t)
	iface := map[string]string{"class": "d3c23a4e-780b-49aa-a85a-802ed877c57d", "level": "7786f8ca-40c0-4b22-8319-8c38da55f51c"}
	r := gdt.DictionaryRestrictions{Interface: &iface}
	value := map[string]json.RawMessage{"class": json.RawMessage(`"` + enumClassA + `"`), "level": json.RawMessage(`2`)}
	if err := VerifyDictionary(db, &gdt.ComponentMap{}, value, r); err != nil {
		t.Fatal(err)
	}
	if len(iface) != 2 {
		t.Errorf("restrictions.Interface was mutated: %v", iface)
	}

	bad := map[string]string{"x": "00000000-0000-0000-0000-000000000000"}
	err := VerifyDictionary(db, &gdt.ComponentMap{}, map[string]json.RawMessage{"x": json.RawMessage(`1`)}, gdt.DictionaryRestrictions{Interface: &bad})
	if err == nil || !strings.Contains(err.Error(), "unknown type") {
		t.Errorf("unknown interface type should fail, got %v", err)
	}
}

func TestIntegrationValidateNodesTypes(t *testing.T) {
	db := testDB(t)
	classes := json.RawMessage(`{"classes":[{"class":"` + enumClassA + `","level":4}]}`)
	stats := func(v int) json.RawMessage {
		b, _ := json.Marshal(map[string]int{"cha": v, "con": v, "dex": v, "int": v, "str": v, "wis": v})
		return b
	}
	for _, tc := range []struct {
		name      string
		values    map[string]json.RawMessage
		wantFail  []string
		wantError bool
	}{
		{"all valid", map[string]json.RawMessage{
			nodeClasses: classes, nodeLevel: json.RawMessage(`10`), nodeHp: json.RawMessage(`50`), nodeBaseStats: stats(10),
		}, nil, false},
		{"level out of range", map[string]json.RawMessage{nodeLevel: json.RawMessage(`21`)}, []string{nodeLevel}, false},
		{"hp above level*5", map[string]json.RawMessage{nodeLevel: json.RawMessage(`2`), nodeHp: json.RawMessage(`11`)}, []string{nodeHp}, false},
		{"stats above level*6", map[string]json.RawMessage{nodeLevel: json.RawMessage(`1`), nodeBaseStats: stats(2)}, []string{nodeBaseStats}, false},
		{"stat above 20", map[string]json.RawMessage{nodeLevel: json.RawMessage(`20`), nodeBaseStats: stats(21)}, []string{nodeBaseStats}, false},
		{"empty class list", map[string]json.RawMessage{nodeClasses: json.RawMessage(`{"classes":[]}`)}, []string{nodeClasses}, false},
		{"unknown node", map[string]json.RawMessage{"00000000-0000-0000-0000-000000000000": json.RawMessage(`1`)}, nil, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ids := make([]string, 0, len(tc.values))
			for id := range tc.values {
				ids = append(ids, id)
			}
			errs, err := ValidateNodesTypes(db, tc.values, ids, testRuleset)
			if tc.wantError {
				if err == nil {
					t.Errorf("expected error, got %v", errs)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(errs) != len(tc.wantFail) {
				t.Fatalf("errors=%v, want failures on %v", errs, tc.wantFail)
			}
			for _, id := range tc.wantFail {
				if _, ok := errs[id]; !ok {
					t.Errorf("expected failure on %s, got %v", id, errs)
				}
			}
		})
	}
}

func TestIntegrationValidateEntity(t *testing.T) {
	db := testDB(t)
	errs, err := ValidateEntity(db, entityFreebird)
	if err != nil {
		t.Fatal(err)
	}
	if len(errs) != 0 {
		t.Errorf("Freebird should be valid, got %v", errs)
	}
}

func TestIntegrationValidateEntityOnlyChecksItsComponents(t *testing.T) {
	db := testDB(t)
	errs, err := ValidateEntity(db, entityTest)
	if err != nil {
		t.Fatal(err)
	}
	if len(errs) != 0 {
		t.Errorf("entity with only a valid Level should be valid, got %v", errs)
	}
}

func TestIntegrationValidateRowTypes(t *testing.T) {
	db := testDB(t)
	const typeLevel, typeStat = "7786f8ca-40c0-4b22-8319-8c38da55f51c", "f79c6225-baf0-45af-81c6-8eb45d1a961f"
	values := map[string]json.RawMessage{typeLevel: json.RawMessage(`5`), typeStat: json.RawMessage(`25`)}
	rows, err := db.Query(`
		SELECT t.id, gv.base_name, gv.scripts, gv.restrictions
		FROM types t
		JOIN get_verifications($2::uuid[]) gv ON gv.target_type_id = t.id
		WHERE t.ruleset_id = $1
	`, testRuleset, pq.Array([]string{typeLevel, typeStat}))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = rows.Close() }()
	errs, count, err := ValidateRowTypes(db, &values, &gdt.ComponentMap{}, rows)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Errorf("count=%d, want 2", count)
	}
	if _, ok := errs[typeStat]; !ok || len(errs) != 1 {
		t.Errorf("want only stat (25 > 20) to fail, got %v", errs)
	}
}

func TestIntegrationDatabaseErrorIsWrapped(t *testing.T) {
	db := testDB(t)
	_, err := ValidateNodesTypes(db, map[string]json.RawMessage{"x": json.RawMessage(`1`)}, []string{"x"}, testRuleset)
	if !errors.Is(err, gdt.ErrDatabase) {
		t.Errorf("malformed uuid should give gdt.ErrDatabase, got %v", err)
	}
}
