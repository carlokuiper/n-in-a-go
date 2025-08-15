package n_in_a_go

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Game struct {
	nInARow  int
	board    [][]int
	mu       sync.Mutex
	history  []Move
	finished bool
}

func (g *Game) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
	_, err := fmt.Fprintln(w, g)
	if err != nil {
		return
	}
}

type Config struct {
	NInARow   int `json:"NInARow"`
	BoardSize int `json:"BoardSize"`
}

func (g *Game) Start(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
	var config Config
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer func() { _ = r.Body.Close() }()
	if config.NInARow < config.BoardSize || config.NInARow < 0 || config.BoardSize < 0 {
		http.Error(w, "invalid config", http.StatusUnprocessableEntity)
		return
	}
	g.mu.Lock()
	g.New(config)
	g.mu.Unlock()
}

func (g *Game) New(config Config) {
	board := make([][]int, config.BoardSize) // create rows
	for i := range board {
		board[i] = make([]int, config.BoardSize) // for each row add the columns
	}
	g.nInARow = config.NInARow
	g.board = board
	g.history = nil
	g.finished = false
}

type Move struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (g *Game) Move(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
	var move Move
	if err := json.NewDecoder(r.Body).Decode(&move); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer func() { _ = r.Body.Close() }()
	g.mu.Lock()
	nextValue := g.nextValue()
	if nextValue == 0 {
		http.Error(w, "cannot determine next value", http.StatusInternalServerError)
	}
	err := g.update(move, g.nextValue())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
	g.mu.Unlock()
}

func (g *Game) nextValue() int {
	if g.history == nil {
		return 1
	}
	previousMove := g.history[len(g.history)-1]
	switch g.board[previousMove.X][previousMove.Y] {
	case 1:
		return 2
	case 2:
		return 1
	default:
		return 0
	}
}

func (g *Game) update(move Move, nextValue int) error {
	if move.X >= len(g.board) || move.Y >= len(g.board) || move.X < 0 || move.Y < 0 {
		return fmt.Errorf("move invalid")
	}
	if g.history == nil {
		g.board[move.X][move.Y] = nextValue
		g.history = []Move{move}
		return nil
	}
	previousMove := g.history[len(g.history)-1]
	if previousMove == move {
		return nil
	}
	if g.finished {
		return fmt.Errorf("game already finished")
	}
	if g.board[move.X][move.Y] != 0 {
		return fmt.Errorf("move not free")
	}
	g.board[move.X][move.Y] = nextValue
	g.history = append(g.history, move)
	g.finished = finished(g.board, g.nInARow)
	return nil
}

func finished(board [][]int, nInARow int) bool {
	for _, row := range board {
		if containsNInARow(row, nInARow) {
			return true
		}
	}
	for i := range board {
		column := make([]int, len(board))
		for j, row := range board {
			column[j] = row[i]
		}
		if containsNInARow(column, nInARow) {
			return true
		}
	}
	// (off) diagonal
	for i := range board {
		x := 0
		y := i
		offDiagonal := make([]int, 0, len(board))
		for j := range board {
			if y+j >= len(board) {
				continue
			}
			offDiagonal = append(offDiagonal, board[x+j][y+j])
		}
		if containsNInARow(offDiagonal, nInARow) {
			return true
		}
		x = i
		y = 0
		offDiagonal = make([]int, 0, len(board))
		for j := range board {
			if x+j >= len(board) {
				continue
			}
			offDiagonal = append(offDiagonal, board[x+j][y+j])
		}
		if containsNInARow(offDiagonal, nInARow) {
			return true
		}
	}
	// (off) anti diagonal
	for i := range board {
		x := len(board) - 1 - i
		y := 0
		offAntiDiagonal := make([]int, 0, len(board))
		for j := range board {
			if x-j < 0 {
				continue
			}
			offAntiDiagonal = append(offAntiDiagonal, board[x-j][y+j])
		}
		if containsNInARow(offAntiDiagonal, nInARow) {
			return true
		}
		x = len(board) - 1
		y = i
		offAntiDiagonal = make([]int, 0, len(board))
		for j := range board {
			if y+j >= len(board) {
				continue
			}
			offAntiDiagonal = append(offAntiDiagonal, board[x-j][y+j])
		}
		if containsNInARow(offAntiDiagonal, nInARow) {
			return true
		}
	}
	return false
}

func containsNInARow(row []int, n int) bool {
	var count, lastEl int
	for _, el := range row {
		switch el {
		case 0:
			count = 0
			lastEl = 0
		case lastEl:
			count++
			if count == n {
				return true
			}
		default:
			count = 1
			lastEl = el
		}
	}
	return false
}
