package builder

import (
	"errors"
	"reflect"
	"testing"
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
	reg := newTestRegistry()
	_ = reg.Add(okConfig{})
	_ = reg.Add(anotherConfig{})

	if err := Build(reg); err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	values := reg.values()
	if len(values) != 2 {
		t.Fatalf("expected 2 resources, got %d: %v", len(values), values)
	}
}

func TestBuild_inheritsTags(t *testing.T) {
	reg := newTestRegistry()
	_ = reg.AddWithTags(okConfig{}, "replaceable")

	if err := Build(reg); err != nil {
		t.Fatal(err)
	}

	tags := reg.firstTags()
	if len(tags) != 1 || tags[0] != "replaceable" {
		t.Fatalf("expected inherited tag, got %v", tags)
	}
}

func TestBuild_skipsNonBuilder(t *testing.T) {
	reg := newTestRegistry()
	_ = reg.Add("keep-me")
	_ = reg.Add(okConfig{})

	if err := Build(reg); err != nil {
		t.Fatal(err)
	}

	values := reg.values()
	if len(values) != 2 {
		t.Fatalf("expected kept string + built resource, got %v", values)
	}
	if values[0] != "keep-me" {
		t.Fatalf("non-builder entry removed or reordered: %v", values)
	}
}

func TestBuild_buildError(t *testing.T) {
	reg := newTestRegistry()
	_ = reg.Add(failConfig{})

	err := Build(reg)
	if err == nil || err.Error() != "builder: builder.failConfig: build failed" {
		t.Fatalf("expected build error, got %v", err)
	}

	if len(reg.values()) != 1 {
		t.Fatalf("config should remain on build error, entries=%d", len(reg.values()))
	}
}

func TestBuild_nilRegistry(t *testing.T) {
	err := Build(nil)
	if err == nil || err.Error() != "builder: nil registry" {
		t.Fatalf("expected nil registry error, got %v", err)
	}
}

func TestBuild_ptrReceiverConfig(t *testing.T) {
	reg := newTestRegistry()
	cfg := &ptrConfig{}
	_ = reg.Add(cfg)

	if err := Build(reg); err != nil {
		t.Fatal(err)
	}

	got, ok := reg.getOneByType(reflect.TypeOf(cfg))
	if !ok {
		t.Fatal("expected built value")
	}
	if got != cfg {
		t.Fatalf("expected built value %p, got %v", cfg, got)
	}
}

func TestBuild_removesConfigEntry(t *testing.T) {
	reg := newTestRegistry()
	cfg := okConfig{}
	_ = reg.Add(cfg)

	if err := Build(reg); err != nil {
		t.Fatal(err)
	}

	if _, ok := reg.getOneByType(reflect.TypeOf(okConfig{})); ok {
		t.Fatal("config type should be removed from registry")
	}
}
