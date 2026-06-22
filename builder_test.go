package builder

import (
	"errors"
	"reflect"
	"testing"

	"github.com/omcrgnt/res"
)

type okConfig struct{}

func (okConfig) Build() (any, error) { return "resource-a", nil }

type anotherConfig struct{}

func (anotherConfig) Build() (any, error) { return 42, nil }

type failConfig struct{}

func (failConfig) Build() (any, error) { return nil, errors.New("build failed") }

type ptrConfig struct{}

func (c *ptrConfig) Build() (any, error) { return c, nil }

func TestBuild_success(t *testing.T) {
	reg := res.New()
	_ = reg.Add(okConfig{})
	_ = reg.Add(anotherConfig{})

	if err := Build(reg); err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	var values []any
	reg.WalkEntries(func(e res.Entry) bool {
		values = append(values, e.Value)
		return true
	})
	if len(values) != 2 {
		t.Fatalf("expected 2 resources, got %d: %v", len(values), values)
	}
}

func TestBuild_inheritsReplaceableTag(t *testing.T) {
	reg := res.New()
	_ = reg.AddWithTags(okConfig{}, res.TagReplaceable)

	if err := Build(reg); err != nil {
		t.Fatal(err)
	}

	reg.WalkEntries(func(e res.Entry) bool {
		if !e.Replaceable() {
			t.Fatal("expected Replaceable on built resource")
		}
		return false
	})
}

func TestBuild_inheritsFixedTag(t *testing.T) {
	reg := res.New()
	_ = reg.AddWithTags(okConfig{}, res.TagFixed)

	if err := Build(reg); err != nil {
		t.Fatal(err)
	}

	reg.WalkEntries(func(e res.Entry) bool {
		if !e.Fixed() {
			t.Fatal("expected Fixed on built resource")
		}
		return false
	})
}

func TestBuild_skipsNonBuilder(t *testing.T) {
	reg := res.New()
	_ = reg.Add("keep-me")
	_ = reg.Add(okConfig{})

	if err := Build(reg); err != nil {
		t.Fatal(err)
	}

	var values []any
	reg.WalkEntries(func(e res.Entry) bool {
		values = append(values, e.Value)
		return true
	})
	if len(values) != 2 {
		t.Fatalf("expected kept string + built resource, got %v", values)
	}
	if values[0] != "keep-me" {
		t.Fatalf("non-builder entry removed or reordered: %v", values)
	}
}

func TestBuild_buildError(t *testing.T) {
	reg := res.New()
	_ = reg.Add(failConfig{})

	err := Build(reg)
	if err == nil || err.Error() != "builder: builder.failConfig: build failed" {
		t.Fatalf("expected build error, got %v", err)
	}

	n := 0
	reg.WalkEntries(func(_ res.Entry) bool {
		n++
		return true
	})
	if n != 1 {
		t.Fatalf("config should remain on build error, entries=%d", n)
	}
}

func TestBuild_nilRegistry(t *testing.T) {
	err := Build(nil)
	if err == nil || err.Error() != "builder: nil registry" {
		t.Fatalf("expected nil registry error, got %v", err)
	}
}

func TestBuild_ptrReceiverConfig(t *testing.T) {
	reg := res.New()
	cfg := &ptrConfig{}
	_ = reg.Add(cfg)

	if err := Build(reg); err != nil {
		t.Fatal(err)
	}

	got, err := reg.GetOneByType(reflect.TypeOf(cfg))
	if err != nil {
		t.Fatal(err)
	}
	if got != cfg {
		t.Fatalf("expected built value %p, got %v", cfg, got)
	}
}

func TestBuild_removesConfigEntry(t *testing.T) {
	reg := res.New()
	cfg := okConfig{}
	_ = reg.Add(cfg)

	if err := Build(reg); err != nil {
		t.Fatal(err)
	}

	if _, err := reg.GetOneByType(reflect.TypeOf(okConfig{})); err == nil {
		t.Fatal("config type should be removed from registry")
	}
}
