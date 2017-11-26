package logicclock

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type Clock struct {
	clock map[string]int
}

func NewClock() *Clock {
	c := &Clock{}

	c.clock = make(map[string]int)
	return c
}

func isLast(size int, pos int) bool {
	return size == 1 || pos == size-1
}

func (c *Clock) String() string {
	var s string
	var pos int

	s += fmt.Sprint("[")
	for id, v := range c.clock {
		if isLast(len(c.clock), pos) {
			s += fmt.Sprintf("%q:%d] ", id, v)
		} else {
			s += fmt.Sprintf("%q:%d, ", id, v)
		}
		pos++
	}
	if len(c.clock) == 0 {
		s += fmt.Sprint("] ")
	}
	return s
}

func (c *Clock) Add(id string) {
	if _, ok := c.clock[id]; ok {
		c.clock[id]++
	} else {
		c.clock[id] = 1
	}
}

// Take two cloks and get the max
// between them...
func maxclock(source *Clock, dest *Clock) {
	for k, s := range source.clock {
		d := dest.clock[k]
		if d > s {
			source.clock[k] = d
		}
	}
	for k, d := range dest.clock {
		s := source.clock[k]
		if d > s {
			source.clock[k] = d
		}
	}
}

func (source *Clock) Max(dest *Clock) {
	maxclock(source, dest)
	maxclock(dest, source)
}

func order(source *Clock, dest *Clock) bool {
	for k, s := range source.clock {
		if d := dest.clock[k]; s > d {
			return false
		}
	}
	return true
}

func iscmp(c1 *Clock, c2 *Clock) bool {
	return order(c1, c2) || order(c2, c1)
}

func (source *Clock) IsLess(dest *Clock) (bool, error) {
	err := errors.New("The clocks are not comparables")

	if !iscmp(source, dest) {
		return false, err
	}
	for k, s := range source.clock {
		if d := dest.clock[k]; s > d {
			return false, nil
		}
	}
	return true, nil
}

func (c *Clock) Json() string {
	data, err := json.Marshal(c.clock)

	if err != nil {
		log.Fatalf("JSON Marshalling failed : %v\n", err)
	}
	return string(data)
}

func (c *Clock) Data(text string) {
	data := []byte(text)
	err := json.Unmarshal(data, &c.clock)

	if err != nil {
		log.Fatalf("Error data clock json : %v\n", err)
	}
}
