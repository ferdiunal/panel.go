package resource

import "testing"

func TestBaseGridEnabled_DefaultTrue(t *testing.T) {
	r := &Base{}
	if !r.IsGridEnabled() {
		t.Fatal("expected base resource grid to be enabled by default")
	}
}

func TestBaseGridEnabled_Setter(t *testing.T) {
	r := &Base{}

	if returned := r.SetGridEnabled(false); returned != r {
		t.Fatal("expected SetGridEnabled to support method chaining")
	}
	if r.IsGridEnabled() {
		t.Fatal("expected base resource grid to be disabled")
	}

	r.SetGridEnabled(true)
	if !r.IsGridEnabled() {
		t.Fatal("expected base resource grid to be enabled")
	}
}

func TestOptimizedBaseGridEnabled_DefaultTrue(t *testing.T) {
	r := &OptimizedBase{}
	if !r.IsGridEnabled() {
		t.Fatal("expected optimized resource grid to be enabled by default")
	}
}

func TestOptimizedBaseGridEnabled_Setter(t *testing.T) {
	r := &OptimizedBase{}

	if returned := r.SetGridEnabled(false); returned != r {
		t.Fatal("expected SetGridEnabled to support method chaining")
	}
	if r.IsGridEnabled() {
		t.Fatal("expected optimized resource grid to be disabled")
	}

	r.SetGridEnabled(true)
	if !r.IsGridEnabled() {
		t.Fatal("expected optimized resource grid to be enabled")
	}
}
