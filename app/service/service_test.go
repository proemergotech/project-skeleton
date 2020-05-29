package service

import (
	"strconv"
	"testing"
)

// todo: remove
//  Example unit tests
func TestDummy(t *testing.T) {
	runDummyTestCase(t, 5, 5, 10)
}

func TestDummyWithTables(t *testing.T) {
	for index, test := range []struct {
		x int
		y int
		n int
	}{
		{1, 1, 2},
		{1, 2, 3},
		{2, 2, 4},
		{5, 2, 7},
	} {
		testCase := test
		t.Run("Test Case "+strconv.Itoa(index), func(t *testing.T) {
			runDummyTestCase(t, testCase.x, testCase.y, testCase.n)
		})

	}
}

func runDummyTestCase(t *testing.T, x int, y int, expected int) {
	total := sum(x, y)
	if total != expected {
		t.Fatalf("Sum of (%d+%d) was incorrect, got: %d, want: %d.", x, y, total, expected)
	}
}

func sum(x int, y int) int {
	return x + y
}
