package ndoapi

import "testing"

func TestPathResolveObjectFromResponse(t *testing.T) {
	template := map[string]any{
		"l3outTemplate": map[string]any{
			"l3outs": []any{
				map[string]any{"uuid": "other", "name": "other", "externalEpgs": []any{map[string]any{"name": "edge"}}},
				map[string]any{
					"uuid": "l3out-1", "name": "router",
					"externalEpgs": []any{
						map[string]any{"uuid": "other-epg", "tenant": "other", "name": "edge"},
						map[string]any{
							"uuid": "epg-1", "tenant": "tenant-a", "name": "edge",
							"annotations": []any{
								map[string]any{"key": "owner", "scope": "other"},
								map[string]any{"key": "unrelated"},
								map[string]any{"key": "owner", "scope": "service", "value": "team-a"},
							},
						},
					},
				},
			},
		},
	}

	newPath := func(uuid, l3outName, epgTenant, epgName string) Path {
		return NewPath(
			PathStep{Field: "l3outTemplate"},
			PathStep{Field: "l3outs", Selector: &ObjectSelector{
				UUID: ObjectIdentifier{Field: "uuid", Value: uuid},
				Keys: []ObjectIdentifier{{Field: "name", Value: l3outName}},
			}},
			PathStep{Field: "externalEpgs", Selector: &ObjectSelector{
				Keys: []ObjectIdentifier{{Field: "tenant", Value: epgTenant}, {Field: "name", Value: epgName}},
			}},
			PathStep{Field: "annotations", Selector: &ObjectSelector{
				Keys: []ObjectIdentifier{{Field: "key", Value: "owner"}, {Field: "scope", Value: "service"}},
			}},
		)
	}
	path := newPath("l3out-1", "other", "tenant-a", "edge")

	resolved, found, err := path.ResolveObjectFromResponse(template)
	if err != nil || !found || resolved.Object["value"] != "team-a" {
		t.Fatalf("unexpected nested lookup: %#v, found=%t, err=%v", resolved, found, err)
	}
	if got, want := resolved.PatchPath(), "/l3outTemplate/l3outs/1/externalEpgs/1/annotations/2"; got != want {
		t.Fatalf("unexpected nested patch path: got %q, want %q", got, want)
	}
	operations := NewPatchOperations(resolved.Object, resolved.PatchPath())
	operations.Set("value", "team-b")
	if got, want := operations.Operations()[0].Path, "/l3outTemplate/l3outs/1/externalEpgs/1/annotations/2/value"; got != want {
		t.Fatalf("unexpected nested field patch path: got %q, want %q", got, want)
	}
	if _, found, err := newPath("missing", "router", "tenant-a", "edge").ResolveObjectFromResponse(template); err != nil || found {
		t.Fatal("UUID mismatch must not fall back to keys")
	}
	if _, found, err := newPath("", "router", "tenant-a", "edge").ResolveObjectFromResponse(template); err != nil || !found {
		t.Fatal("expected key fallback when UUID is absent")
	}
	if _, found, err := newPath("", "router", "tenant-a", "missing").ResolveObjectFromResponse(template); err != nil || found {
		t.Fatal("all key identifiers must match the same object")
	}
	if _, found, err := newPath("", "router", "", "edge").ResolveObjectFromResponse(template); err != nil || found {
		t.Fatal("incomplete keys must not select an object")
	}
	malformed := map[string]any{"l3outTemplate": map[string]any{"l3outs": []any{map[string]any{"uuid": "l3out-1", "name": "router", "externalEpgs": []any{map[string]any{"tenant": "tenant-a", "name": "edge", "annotations": []any{"invalid object"}}}}}}}
	if _, found, err := path.ResolveObjectFromResponse(malformed); err == nil || found {
		t.Fatalf("expected malformed annotation entry error, found=%t, err=%v", found, err)
	}
	malformed = map[string]any{"l3outTemplate": map[string]any{"l3outs": "invalid collection"}}
	if _, found, err := path.ResolveObjectFromResponse(malformed); err == nil || found {
		t.Fatalf("expected malformed L3Out collection error, found=%t, err=%v", found, err)
	}
	malformed = map[string]any{"l3outTemplate": []any{}}
	if _, found, err := path.ResolveObjectFromResponse(malformed); err == nil || found {
		t.Fatalf("expected malformed template object error, found=%t, err=%v", found, err)
	}
	malformed = map[string]any{"l3outTemplate": map[string]any{"l3outs": []any{map[string]any{"uuid": 42}}}}
	if _, found, err := path.ResolveObjectFromResponse(malformed); err == nil || found {
		t.Fatalf("expected malformed L3Out UUID error, found=%t, err=%v", found, err)
	}
}

func TestPathTopLevelAppend(t *testing.T) {
	path := NewPath(
		PathStep{Field: "l3outTemplate"},
		PathStep{Field: "l3outs", Selector: &ObjectSelector{}},
	)
	appendPath, found, err := path.AppendPath()
	if err != nil || !found || appendPath != "/l3outTemplate/l3outs/-" {
		t.Fatalf("unexpected append path without lookup: %q, found=%t, err=%v", appendPath, found, err)
	}
	if _, found, err := path.ResolveObjectFromResponse(nil); err != nil || found {
		t.Fatal("object lookup requires a template")
	}
}

func TestPathRejectsUnresolvedParent(t *testing.T) {
	path := NewPath(
		PathStep{Field: "container"},
		PathStep{Field: "parents", Selector: &ObjectSelector{Keys: []ObjectIdentifier{{Field: "name", Value: "parent"}}}},
		PathStep{Field: "children", Selector: &ObjectSelector{Keys: []ObjectIdentifier{{Field: "name", Value: "child"}}}},
	)
	if _, found, err := path.AppendPath(); err != nil || found {
		t.Fatal("nested append must resolve the parent first")
	}
	template := map[string]any{"container": map[string]any{"parents": []any{map[string]any{"name": "parent", "children": []any{map[string]any{"name": "child"}}}}}}
	if appendPath, found, err := path.AppendPath(template); err != nil || !found || appendPath != "/container/parents/0/children/-" {
		t.Fatalf("unexpected nested append path: %q, found=%t, err=%v", appendPath, found, err)
	}
	resolved, found, err := path.ResolveObjectFromResponse(template)
	if err != nil || !found || resolved.PatchPath() != "/container/parents/0/children/0" {
		t.Fatalf("unexpected path: %#v, found=%t, err=%v", resolved, found, err)
	}
}

func TestPathCopiesStepsAndSelectors(t *testing.T) {
	selector := &ObjectSelector{
		UUID: ObjectIdentifier{Field: "uuid", Value: "l3out-1"},
		Keys: []ObjectIdentifier{{Field: "name", Value: "router"}},
	}
	steps := []PathStep{
		{Field: "l3outTemplate"},
		{Field: "l3outs", Selector: selector},
	}
	path := NewPath(steps...)
	steps[0].Field = "changed"
	steps[1].Field = "changed"
	selector.UUID.Value = "changed"
	selector.Keys[0].Value = "changed"

	if appendPath, found, err := path.AppendPath(); err != nil || !found || appendPath != "/l3outTemplate/l3outs/-" {
		t.Fatalf("path changed after construction: %q, found=%t, err=%v", appendPath, found, err)
	}
	template := map[string]any{
		"l3outTemplate": map[string]any{
			"l3outs": []any{map[string]any{"uuid": "l3out-1", "name": "router"}},
		},
	}
	if resolved, found, err := path.ResolveObjectFromResponse(template); err != nil || !found || resolved.PatchPath() != "/l3outTemplate/l3outs/0" {
		t.Fatalf("selector changed after construction: %#v, found=%t, err=%v", resolved, found, err)
	}
	selector.UUID.Value = ""
	selector.Keys[0].Value = "router"
	steps[0].Field = "l3outTemplate"
	steps[1].Field = "l3outs"
	fallbackPath := NewPath(steps...)
	selector.Keys[0].Value = "changed"
	if resolved, found, err := fallbackPath.ResolveObjectFromResponse(template); err != nil || !found || resolved.PatchPath() != "/l3outTemplate/l3outs/0" {
		t.Fatalf("keys changed after construction: %#v, found=%t, err=%v", resolved, found, err)
	}

	var empty Path
	if _, found, err := empty.AppendPath(); err != nil || found {
		t.Fatal("zero-value path cannot append")
	}
	if _, found, err := empty.ResolveObjectFromResponse(template); err != nil || found {
		t.Fatal("zero-value path cannot resolve an object")
	}
}
