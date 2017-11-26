package logicclock

import (
	"fmt"
	"testing"
)

func isGreather(c1, c2 *Clock) bool {
	less, err := c1.IsLess(c2)

	return err == nil && !less
}

func isTestEquals(c1, c2 *Clock) bool {
	if len(c1.clock) != len(c2.clock) {
		return false
	}
	for k, v1 := range c1.clock {
		v2 := c2.clock[k]
		if v1 != v2 {
			return false
		}
	}
	return true
}

func failMax(c1, c2, c3 *Clock) bool {
	return isGreather(c1, c2) || !isGreather(c1, c3) || !isTestEquals(c1, c2)
}

func TestMax(t *testing.T) {
	c1 := NewClock()
	c2 := NewClock()
	c3 := NewClock()

	c1.Add("a") // [1, 0]
	c1.Add("a") // [2, 0]
	c2.Add("b") // [0, 1]
	c1.Max(c2)  // c1 [2, 1] c2 [2, 1] c3 [0, 0]
	if failMax(c1, c2, c3) {
		t.Error()
	}

	c1.Add("b")
	c1.Add("b")
	c1.Max(c2) // c1 [1,3] c2 [1,3] c3 [0, 0]
	if failMax(c1, c2, c3) {
		t.Error()
	}

	c3.Add("d")
	c3.Add("e")
	c3.Add("f")
	c3.Add("g")
	c3.Add("h")
	c3.Max(c2) // c1[1,3] c2[2,3,1,1,1,1,1] c3[2,3,1,1,1,1,1]
	if failMax(c2 , c3 , c1) {
		t.Error()
	}
}
