package n_in_a_go

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
