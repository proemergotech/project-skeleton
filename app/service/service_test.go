package service

import (
	"strconv"
	"testing"
)

// Example unit tests
func TestDummy(t *testing.T) {
	runDummyTestCase(t, 5, 5, 10)
}

func TestDummyWithTables(t *testing.T) {
	tables := []struct {
		x int
		y int
		n int
	}{
		{1, 1, 2},
		{1, 2, 3},
		{2, 2, 4},
		{5, 2, 7},
	}

	for index, table := range tables {
		t.Run("Test Case "+strconv.Itoa(index), func(t *testing.T) {
			runDummyTestCase(t, table.x, table.y, table.n)
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
