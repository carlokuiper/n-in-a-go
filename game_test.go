package kinago

import "testing"

func TestFinished(t *testing.T) {
	tests := []struct {
		name     string
		board    [][]int
		nInARow  int
		expected bool
	}{{
		name:     "unfinished column",
		board:    [][]int{{0, 0, 1}, {0, 0, 1}, {0, 0, 0}},
		nInARow:  3,
		expected: false,
	}, {
		name:     "finished row",
		board:    [][]int{{1, 1, 1}, {0, 0, 0}, {0, 0, 0}},
		nInARow:  3,
		expected: true,
	}, {
		name:     "off anti diagonal",
		board:    [][]int{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 1}, {0, 0, 1, 0}},
		nInARow:  2,
		expected: true,
	}, {
		name:     "m large finished",
		board:    [][]int{{0, 0, 0, 0, 0, 1}, {0, 0, 0, 0, 1, 0}, {0, 0, 0, 1, 0, 0}},
		nInARow:  3,
		expected: true,
	}, {
		name:     "m large unfinished",
		board:    [][]int{{0, 0, 0, 0, 0, 1}, {0, 0, 0, 0, 1, 0}, {0, 0, 0, 1, 0, 0}},
		nInARow:  4,
		expected: false,
	}, {
		name:     "n large finished",
		board:    [][]int{{0, 0}, {0, 0}, {1, 0}, {0, 1}},
		nInARow:  2,
		expected: true,
	}, {
		name:     "n large unfinished",
		board:    [][]int{{0, 0}, {0, 0}, {1, 0}, {0, 1}},
		nInARow:  3,
		expected: false,
	}}
	t.Parallel()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := finished(test.board, test.nInARow)
			if got != test.expected {
				t.Fail()
			}
		})
	}
}
