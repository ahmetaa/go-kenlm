package kenlm

import (
	"math"
	"os"
	"testing"
)

const charModelPath = "char-6gram-kenlm.binary"
const tinyARPA = "testdata/tiny.arpa"

func TestLoadTinyARPA(t *testing.T) {
	if _, err := os.Stat(tinyARPA); err != nil {
		t.Fatalf("missing fixture %q: %v", tinyARPA, err)
	}
	m, err := LoadModel(tinyARPA)
	if err != nil {
		t.Fatalf("LoadModel: %v", err)
	}
	defer m.Close()

	if got := m.Order(); got != 2 {
		t.Errorf("Order() = %d, want 2", got)
	}
	got := m.Score("a b")
	if math.IsNaN(got) || math.IsInf(got, 0) {
		t.Fatalf("Score returned non-finite: %v", got)
	}
	// <s> a (-0.5) + a b (-0.5) + b </s> (-0.5) = -1.5
	want := -1.5
	if math.Abs(got-want) > 1e-4 {
		t.Errorf("Score(\"a b\") = %v, want %v", got, want)
	}
}

func TestCharModel(t *testing.T) {
	if _, err := os.Stat(charModelPath); err != nil {
		t.Skipf("model %q not present (this is fine in CI): %v", charModelPath, err)
	}
	m, err := LoadModel(charModelPath)
	if err != nil {
		t.Fatalf("LoadModel: %v", err)
	}
	defer m.Close()

	if got := m.Order(); got != 6 {
		t.Errorf("Order() = %d, want 6", got)
	}
	got := m.ScoreChars("merhaba dunya")
	if math.IsNaN(got) || math.IsInf(got, 0) || got >= 0 {
		t.Fatalf("ScoreChars returned %v, expected finite negative", got)
	}
	t.Logf("log10 P(\"merhaba dunya\") = %v", got)
}
