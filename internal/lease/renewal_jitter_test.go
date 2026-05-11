package lease

import (
	"math/rand"
	"testing"
	"time"
)

func deterministicRng(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}

func TestDefaultJitterPolicy_Values(t *testing.T) {
	p := DefaultJitterPolicy()
	if p.MaxFraction != 0.20 {
		t.Errorf("expected MaxFraction 0.20, got %v", p.MaxFraction)
	}
}

func TestJitterPolicy_Apply_ZeroBase(t *testing.T) {
	p := DefaultJitterPolicy()
	result := p.Apply(0, deterministicRng(42))
	if result != 0 {
		t.Errorf("expected 0 for zero base, got %v", result)
	}
}

func TestJitterPolicy_Apply_NegativeBase(t *testing.T) {
	p := DefaultJitterPolicy()
	result := p.Apply(-5*time.Second, deterministicRng(42))
	if result != -5*time.Second {
		t.Errorf("expected unchanged negative base, got %v", result)
	}
}

func TestJitterPolicy_Apply_ZeroFraction(t *testing.T) {
	p := JitterPolicy{MaxFraction: 0}
	base := 10 * time.Second
	result := p.Apply(base, deterministicRng(42))
	if result != base {
		t.Errorf("expected unchanged base with zero fraction, got %v", result)
	}
}

func TestJitterPolicy_Apply_ResultAtLeastBase(t *testing.T) {
	p := DefaultJitterPolicy()
	base := 30 * time.Second
	for i := int64(0); i < 20; i++ {
		result := p.Apply(base, deterministicRng(i))
		if result < base {
			t.Errorf("seed %d: jitter result %v is less than base %v", i, result, base)
		}
	}
}

func TestJitterPolicy_Apply_ResultWithinBounds(t *testing.T) {
	p := DefaultJitterPolicy()
	base := 60 * time.Second
	maxExpected := base + time.Duration(float64(base)*p.MaxFraction)
	for i := int64(0); i < 20; i++ {
		result := p.Apply(base, deterministicRng(i))
		if result > maxExpected {
			t.Errorf("seed %d: result %v exceeds max expected %v", i, result, maxExpected)
		}
	}
}

func TestJitterPolicy_ApplyGlobal_ReturnsAtLeastBase(t *testing.T) {
	p := DefaultJitterPolicy()
	base := 10 * time.Second
	result := p.ApplyGlobal(base)
	if result < base {
		t.Errorf("ApplyGlobal result %v is less than base %v", result, base)
	}
}
