package ndoapi

import "testing"

func TestTranslationAPIAndSchemaValues(t *testing.T) {
	translation := NewTranslationMap(map[string]string{
		"alpha": "alpha-api",
	})

	if got := translation.ToAPI("alpha"); got != "alpha-api" {
		t.Fatalf("unexpected API translation: got %q, want %q", got, "alpha-api")
	}
	if got := translation.ToSchema("alpha-api"); got != "alpha" {
		t.Fatalf("unexpected schema translation: got %q, want %q", got, "alpha")
	}
	if got := translation.ToAPI("unknown"); got != "unknown" {
		t.Fatalf("unexpected unmapped API value: got %q, want %q", got, "unknown")
	}
	if got := translation.ToSchema("unknown-api"); got != "unknown-api" {
		t.Fatalf("unexpected unmapped schema value: got %q, want %q", got, "unknown-api")
	}
	if got := translation.ToAPI("af11"); got != "af11" {
		t.Fatalf("unexpected passthrough API value: got %q, want %q", got, "af11")
	}
	if got := translation.ToSchema("af11"); got != "af11" {
		t.Fatalf("unexpected passthrough schema value: got %q, want %q", got, "af11")
	}
}

func TestNewTranslationMapRejectsAmbiguousAPIValue(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected ambiguous translation to panic")
		}
	}()

	NewTranslationMap(map[string]string{
		"first":  "same-api-value",
		"second": "same-api-value",
	})
}
