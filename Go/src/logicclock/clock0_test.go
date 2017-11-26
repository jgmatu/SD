package logicclock

import (
	"testing"
)

func TestClock(t *testing.T) {
	c1 := NewClock()
	c2 := NewClock()

	c1.Add("a")
	c2.Add("b")
	// [0, 1]
	// [1, 0]
	if _, err := c1.IsLess(c2); err == nil {
		// Not comparables.
		t.Error()
	}
	c2.Add("a")
	// [0, 1] c1
	// [1, 1] c2
	if less, err := c1.IsLess(c2); err != nil || !less {
		// Comparables and c1 is less c2
		t.Error()
	}
	// [0, 1] c1
	// [1 ,2] c2
	c2.Add("a")
	if _, err := c1.IsLess(c2); err != nil {
		// Comparables
		t.Error()
	}
	c1.Add("c")
	// [0, 1, 1] c1
	// [1, 2, 0] c2
	if _, err := c1.IsLess(c2); err == nil {
		// No son comparables.
		t.Error()
	}
	c2.Add("c")
	// [1, 2, 1] c2
	// [0, 1, 1] c1
	if less, err := c2.IsLess(c1); err != nil || less {
		// Son comparables y c1 es mayor que c2.
		t.Error()
	}
	// [0,0,0,1]
	c3 := NewClock()
	c3.Add("d")
	if _, err := c3.IsLess(c2); err == nil {
		// Not comparables.
		t.Error()
	}
	c2.Add("d")
	if less, err := c3.IsLess(c2); err != nil || !less {
		// Comparables and c3 is less c2.
		t.Error()
	}
}
