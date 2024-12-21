package internal_test

import (
	"testing"

	"github.com/svdx9/aoc2024/internal"
)

func TestFloat64AlmostEqual(t *testing.T) {
	f1 := float64(2)
	f2 := float64(2)
	if !internal.Float64AlmostEqual(f1, f2) {
		t.Errorf("expected %f and %f to be almost equal", f1, f2)
	}
}
