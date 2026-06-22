package builder

import (
	"errors"
	"testing"
)

type seedNewPtr struct{}

func (p *seedNewPtr) NewResource() (any, error) { return "new-resource", nil }

type seedNew struct{}

func (seedNew) NewResource() (any, error) { return "new-resource", nil }

type seedBuild struct{}

func (seedBuild) BuildConfig() (Builder, error) { return okConfig{}, nil }

type seedDual struct{}

func (seedDual) NewResource() (any, error) { return nil, nil }
func (seedDual) BuildConfig() (Builder, error) { return okConfig{}, nil }

type seedFailNew struct{}

func (seedFailNew) NewResource() (any, error) { return nil, errors.New("new failed") }

type seedFailBuild struct{}

func (seedFailBuild) BuildConfig() (Builder, error) { return nil, errors.New("build config failed") }

type seedApp struct {
	New seedNew
	Bld seedBuild
}

type seedAppMissing struct {
	X int
}

func TestSeed_registersNewAndBuild(t *testing.T) {
	reg := newTestRegistry()
	app := seedApp{New: seedNew{}, Bld: seedBuild{}}

	if err := Seed(reg, &app); err != nil {
		t.Fatal(err)
	}

	m, ok := SeedMapFor(reg)
	if !ok {
		t.Fatal("expected seed map")
	}
	if _, ok := m["Bld"]; !ok {
		t.Fatal("expected Bld in seed map")
	}
	if _, ok := m["New"]; ok {
		t.Fatal("NewResourceer should not be in seed map")
	}

	values := reg.values()
	if len(values) != 2 {
		t.Fatalf("expected 2 specs, got %d: %v", len(values), values)
	}
	if _, ok := values[0].(newResourceSpec); !ok {
		if _, ok := values[1].(newResourceSpec); !ok {
			t.Fatalf("expected newResourceSpec in reg after Seed, got %T and %T", values[0], values[1])
		}
	}
	for _, v := range values {
		if v == "new-resource" {
			t.Fatal("NewResource must not run during Seed")
		}
	}

	if err := Build(reg); err != nil {
		t.Fatal(err)
	}
	var found bool
	for _, v := range reg.values() {
		if v == "new-resource" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected materialized resource after Build")
	}
}

func TestSeed_dualInterfaceError(t *testing.T) {
	reg := newTestRegistry()
	app := struct{ F seedDual }{F: seedDual{}}
	err := Seed(reg, &app)
	if err == nil || err.Error() != "builder: F: implements both NewResourceer and BuildConfiger" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSeed_missingInterfaceError(t *testing.T) {
	reg := newTestRegistry()
	app := seedAppMissing{}
	err := Seed(reg, &app)
	if err == nil || err.Error() != "builder: X: must implement NewResourceer or BuildConfiger" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuild_newResourceError(t *testing.T) {
	reg := newTestRegistry()
	app := struct{ N seedFailNew }{N: seedFailNew{}}
	if err := Seed(reg, &app); err != nil {
		t.Fatal(err)
	}
	if err := Build(reg); err == nil {
		t.Fatal("expected error")
	}
}

func TestSeed_buildConfigError(t *testing.T) {
	reg := newTestRegistry()
	app := struct{ B seedFailBuild }{B: seedFailBuild{}}
	err := Seed(reg, &app)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSeed_nilPointerReceiver(t *testing.T) {
	reg := newTestRegistry()
	app := struct {
		N *seedNewPtr
	}{}

	if err := Seed(reg, &app); err != nil {
		t.Fatal(err)
	}
	if app.N != nil {
		t.Fatal("Seed must not allocate pointer fields")
	}
}

func TestSeed_nilRegistry(t *testing.T) {
	err := Seed(nil, &seedApp{})
	if err == nil || err.Error() != "builder: nil registry" {
		t.Fatalf("unexpected error: %v", err)
	}
}
