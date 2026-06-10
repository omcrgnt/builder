package builder

import (
	"errors"
	"fmt"
	"testing"
)

type mockRegistrar struct {
	added []any
	err   error
}

func (m *mockRegistrar) Add(v any) error {
	if m.err != nil {
		return m.err
	}
	m.added = append(m.added, v)
	return nil
}

type okModule struct{}

func (okModule) Build() (any, error) { return "resource-a", nil }

type anotherModule struct{}

func (anotherModule) Build() (any, error) { return 42, nil }

type failModule struct{}

func (failModule) Build() (any, error) { return nil, errors.New("build failed") }

type ptrReceiverModule struct{}

func (m *ptrReceiverModule) Build() (any, error) {
	return fmt.Sprintf("ptr-%p", m), nil
}

type valueReceiverModule struct{}

func (valueReceiverModule) Build() (any, error) { return "value-receiver", nil }

func TestBuild_success(t *testing.T) {
	cfg := struct {
		Skip int
		A    okModule
		B    anotherModule
	}{
		A: okModule{},
		B: anotherModule{},
	}

	reg := &mockRegistrar{}
	if err := Build(cfg, reg); err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if len(reg.added) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(reg.added))
	}
	if reg.added[0] != "resource-a" {
		t.Errorf("unexpected first resource: %v", reg.added[0])
	}
	if reg.added[1] != 42 {
		t.Errorf("unexpected second resource: %v", reg.added[1])
	}
}

func TestBuild_pointerCfg(t *testing.T) {
	cfg := &struct {
		A okModule
	}{A: okModule{}}

	reg := &mockRegistrar{}
	if err := Build(cfg, reg); err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	if len(reg.added) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(reg.added))
	}
}

func TestBuild_ptrReceiverModule(t *testing.T) {
	cfg := struct {
		M *ptrReceiverModule
	}{
		M: &ptrReceiverModule{},
	}

	reg := &mockRegistrar{}
	if err := Build(cfg, reg); err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	if len(reg.added) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(reg.added))
	}
}

func TestBuild_skipsNilBuilderField(t *testing.T) {
	cfg := struct {
		M *ptrReceiverModule
	}{}

	reg := &mockRegistrar{}
	if err := Build(cfg, reg); err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	if len(reg.added) != 0 {
		t.Fatalf("expected 0 resources, got %d", len(reg.added))
	}
}

func TestBuild_buildError(t *testing.T) {
	cfg := struct {
		A okModule
		B failModule
	}{}

	err := Build(cfg, &mockRegistrar{})
	if err == nil || err.Error() != "builder: B: build failed" {
		t.Fatalf("expected build error, got %v", err)
	}
}

func TestBuild_registrarError(t *testing.T) {
	cfg := struct {
		A okModule
	}{}

	reg := &mockRegistrar{err: errors.New("duplicate")}
	err := Build(cfg, reg)
	if err == nil || err.Error() != "builder: A: duplicate" {
		t.Fatalf("expected registrar error, got %v", err)
	}
}

func TestBuild_notStruct(t *testing.T) {
	err := Build(10, &mockRegistrar{})
	if err == nil || err.Error() != "builder: want struct, got int" {
		t.Fatalf("expected type error, got %v", err)
	}
}

func TestBuild_nilConfig(t *testing.T) {
	err := Build(nil, &mockRegistrar{})
	if err == nil || err.Error() != "builder: nil config" {
		t.Fatalf("expected nil config error, got %v", err)
	}
}

func TestBuild_nilRegistrar(t *testing.T) {
	err := Build(struct{}{}, nil)
	if err == nil || err.Error() != "builder: nil registrar" {
		t.Fatalf("expected nil registrar error, got %v", err)
	}
}

func TestBuild_nilPointerConfig(t *testing.T) {
	var cfg *struct {
		A okModule
	}

	err := Build(cfg, &mockRegistrar{})
	if err == nil || err.Error() != "builder: nil config" {
		t.Fatalf("expected nil config error, got %v", err)
	}
}

func TestBuild_valueReceiverOnPointerField(t *testing.T) {
	cfg := struct {
		M *valueReceiverModule
	}{
		M: &valueReceiverModule{},
	}

	reg := &mockRegistrar{}
	if err := Build(cfg, reg); err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	if len(reg.added) != 1 || reg.added[0] != "value-receiver" {
		t.Fatalf("unexpected resources: %v", reg.added)
	}
}

func TestBuild_skipsUnexportedFieldContinues(t *testing.T) {
	cfg := struct {
		hidden int
		A      okModule
		B      anotherModule
	}{
		hidden: 99,
		A:      okModule{},
		B:      anotherModule{},
	}

	reg := &mockRegistrar{}
	if err := Build(cfg, reg); err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	if len(reg.added) != 2 {
		t.Fatalf("expected 2 resources after unexported field, got %d", len(reg.added))
	}
}

func TestBuild_skipsNilFieldContinues(t *testing.T) {
	cfg := struct {
		M *ptrReceiverModule
		A okModule
	}{
		A: okModule{},
	}

	reg := &mockRegistrar{}
	if err := Build(cfg, reg); err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	if len(reg.added) != 1 || reg.added[0] != "resource-a" {
		t.Fatalf("expected resource after nil field, got %v", reg.added)
	}
}
