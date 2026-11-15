package scripting

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestEvalRoundTrip(t *testing.T) {
	script := `
def echo(raw, m, l, n=0):
    return {"raw": raw, "m": m, "l": l, "n": n, "none": None}
`
	got, err := NewStarlarkEngine().Eval(script, "echo",
		[]any{json.RawMessage(`{"a":[1,"b",true]}`), map[string]json.RawMessage{"x": json.RawMessage(`1.5`)}, []string{"p", "q"}},
		map[string]any{"n": 3},
	)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{
		"raw":  map[string]any{"a": []any{int64(1), "b", true}},
		"m":    map[string]any{"x": 1.5},
		"l":    []any{"p", "q"},
		"n":    int64(3),
		"none": nil,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v\nwant %#v", got, want)
	}
}

func TestEvalNilRawMessageIsNone(t *testing.T) {
	got, err := NewStarlarkEngine().Eval("def f(x):\n    return x == None\n", "f", []any{json.RawMessage(nil)}, nil)
	if err != nil || got != true {
		t.Errorf("got %v, %v", got, err)
	}
}

func TestEvalWhileAllowedButBounded(t *testing.T) {
	e := NewStarlarkEngine()
	got, err := e.Eval("def f():\n    i = 0\n    while i < 10:\n        i += 1\n    return i\n", "f", nil, nil)
	if err != nil || got != int64(10) {
		t.Fatalf("finite while: got %v, %v", got, err)
	}

	e.MaxSteps = 100_000
	_, err = e.Eval("def f():\n    while True:\n        pass\n", "f", nil, nil)
	if err == nil || !strings.Contains(err.Error(), "too many steps") {
		t.Errorf("infinite loop should be stopped, got %v", err)
	}
}

func TestEvalErrors(t *testing.T) {
	e := NewStarlarkEngine()
	if _, err := e.Eval("x = 1\n", "f", nil, nil); err == nil {
		t.Error("missing function should fail")
	}
	if _, err := e.Eval("f = 1\n", "f", nil, nil); err == nil {
		t.Error("non callable should fail")
	}
	if _, err := e.Eval("def f():\n    return set([1])\n", "f", nil, nil); err == nil {
		t.Error("unsupported return type should fail")
	}
}

func TestAsOptionalString(t *testing.T) {
	if s, err := AsOptionalString(nil); s != nil || err != nil {
		t.Error("nil")
	}
	if s, err := AsOptionalString(""); s != nil || err != nil {
		t.Error("empty")
	}
	if s, err := AsOptionalString("x"); s == nil || *s != "x" || err != nil {
		t.Error("value")
	}
	if _, err := AsOptionalString(int64(1)); err == nil {
		t.Error("wrong type")
	}
}
