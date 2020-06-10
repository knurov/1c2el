package helper

import "testing"

type rangeFixture struct {
	groupRange GroupRange
	start      uint8
	end        uint8
}

var testRange = []rangeFixture{
	{"1..2", 1, 2},
}

//TestGetRange test for GetRange
func TestGetRange(t testing.T) {
	for _, fixture := range testRange {
		start, end := fixture.groupRange.GetRange()
		if start != fixture.start {
			t.Error("for", fixture.groupRange, "expected start", fixture.start, "got", start)
		}
		if end != fixture.end {
			t.Error("for", fixture.groupRange, "expected end", fixture.end, "got", end)
		}
	}

}
