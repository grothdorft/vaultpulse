package lease

import (
	"testing"
	"time"
)

func TestDefaultWindowPolicy_Values(t *testing.T) {
	p := DefaultWindowPolicy()
	if p.MinWindow != 30*time.Second {
		t.Errorf("expected MinWindow=30s, got %v", p.MinWindow)
	}
	if p.MaxWindow != 24*time.Hour {
		t.Errorf("expected MaxWindow=24h, got %v", p.MaxWindow)
	}
	if p.DefaultFraction != 0.25 {
		t.Errorf("expected DefaultFraction=0.25, got %v", p.DefaultFraction)
	}
}

func TestRenewalWindow_Compute_Typical(t *testing.T) {
	rw := NewRenewalWindow(DefaultWindowPolicy())
	ttl := 4 * time.Hour
	w := rw.Compute("lease-1", ttl)
	expected := time.Duration(float64(ttl) * 0.25) // 1h
	if w != expected {
		t.Errorf("expected %v, got %v", expected, w)
	}
}

func TestRenewalWindow_Compute_ClampedToMin(t *testing.T) {
	rw := NewRenewalWindow(DefaultWindowPolicy())
	// 10s * 0.25 = 2.5s < MinWindow(30s)
	w := rw.Compute("lease-2", 10*time.Second)
	if w != 30*time.Second {
		t.Errorf("expected 30s (min), got %v", w)
	}
}

func TestRenewalWindow_Compute_ClampedToMax(t *testing.T) {
	rw := NewRenewalWindow(DefaultWindowPolicy())
	// 1000h * 0.25 = 250h > MaxWindow(24h)
	w := rw.Compute("lease-3", 1000*time.Hour)
	if w != 24*time.Hour {
		t.Errorf("expected 24h (max), got %v", w)
	}
}

func TestRenewalWindow_GetAfterCompute(t *testing.T) {
	rw := NewRenewalWindow(DefaultWindowPolicy())
	rw.Compute("lease-4", 4*time.Hour)
	w, ok := rw.Get("lease-4")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if w != time.Hour {
		t.Errorf("expected 1h, got %v", w)
	}
}

func TestRenewalWindow_GetUnknown(t *testing.T) {
	rw := NewRenewalWindow(DefaultWindowPolicy())
	_, ok := rw.Get("unknown")
	if ok {
		t.Error("expected not found for unknown lease")
	}
}

func TestRenewalWindow_Remove(t *testing.T) {
	rw := NewRenewalWindow(DefaultWindowPolicy())
	rw.Compute("lease-5", 2*time.Hour)
	rw.Remove("lease-5")
	_, ok := rw.Get("lease-5")
	if ok {
		t.Error("expected entry to be removed")
	}
}

func TestRenewalWindow_Snapshot(t *testing.T) {
	rw := NewRenewalWindow(DefaultWindowPolicy())
	rw.Compute("a", 4*time.Hour)
	rw.Compute("b", 8*time.Hour)
	snap := rw.Snapshot()
	if len(snap) != 2 {
		t.Errorf("expected 2 entries, got %d", len(snap))
	}
	if snap["a"] != time.Hour {
		t.Errorf("expected a=1h, got %v", snap["a"])
	}
	if snap["b"] != 2*time.Hour {
		t.Errorf("expected b=2h, got %v", snap["b"])
	}
}
