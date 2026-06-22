package builder

import "reflect"

type testRegistry struct {
	entries []testEntry
}

type testEntry struct {
	value any
	tags  []any
}

func newTestRegistry() *testRegistry {
	return &testRegistry{}
}

func (r *testRegistry) Add(v any) error {
	r.entries = append(r.entries, testEntry{value: v})
	return nil
}

func (r *testRegistry) AddWithTags(v any, tags ...any) error {
	r.entries = append(r.entries, testEntry{value: v, tags: append([]any(nil), tags...)})
	return nil
}

func (r *testRegistry) WalkEntries(fn func(value any, tags []any) bool) {
	for _, e := range r.entries {
		if !fn(e.value, e.tags) {
			return
		}
	}
}

func (r *testRegistry) Remove(v any) error {
	for i, e := range r.entries {
		if e.value == v {
			r.entries = append(r.entries[:i], r.entries[i+1:]...)
			return nil
		}
	}
	return nil
}

func (r *testRegistry) values() []any {
	out := make([]any, len(r.entries))
	for i, e := range r.entries {
		out[i] = e.value
	}
	return out
}

func (r *testRegistry) firstTags() []any {
	if len(r.entries) == 0 {
		return nil
	}
	return r.entries[0].tags
}

func (r *testRegistry) getOneByType(t reflect.Type) (any, bool) {
	for _, e := range r.entries {
		if reflect.TypeOf(e.value) == t {
			return e.value, true
		}
	}
	return nil, false
}
