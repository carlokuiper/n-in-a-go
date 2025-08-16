package kinago

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Game struct {
	kInARow int
	board   [][]int
	mu       sync.Mutex
	history  []Move
	finished bool
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
	g.finished = finished(g.board, g.kInARow)
	return nil
}

func finished(board [][]int, k int) bool {
	for _, row := range board {
		if kInARow(row, k) {
			return true
		}
	}
	if len(board) == 0 {
		return false
	}
	m := len(board[0])
	for i := range m {
		column := make([]int, len(board))
		for j, row := range board {
			column[j] = row[i]
		}
		if kInARow(column, k) {
			return true
		}
	}
	// (off) diagonal
	for i := range m {
		x := 0
		y := i
		offDiagonal := make([]int, 0, len(board))
		for j := range board {
			if y+j >= m || x+j >= len(board) {
				continue
			}
			offDiagonal = append(offDiagonal, board[x+j][y+j])
		}
		if kInARow(offDiagonal, k) {
			return true
		}
	}
	for i := range board {
		x := i
		y := 0
		offDiagonal := make([]int, 0, len(board))
		for j := range m {
			if y+j >= m || x+j >= len(board) {
				continue
			}
			offDiagonal = append(offDiagonal, board[x+j][y+j])
		}
		if kInARow(offDiagonal, k) {
			return true
		}
	}
	// (off) anti diagonal
	for i := range board {
		x := len(board) - 1 - i
		y := 0
		offAntiDiagonal := make([]int, 0, len(board))
		for j := range m {
			if x-j < 0 || y+j > m {
				continue
			}
			offAntiDiagonal = append(offAntiDiagonal, board[x-j][y+j])
		}
		if kInARow(offAntiDiagonal, k) {
			return true
		}
	}
	for i := range m {
		x := len(board) - 1
		y := i
		offAntiDiagonal := make([]int, 0, len(board))
		for j := range board {
			if x-j < 0 || y+j >= m {
				continue
			}
			offAntiDiagonal = append(offAntiDiagonal, board[x-j][y+j])
		}
		if kInARow(offAntiDiagonal, k) {
			return true
		}
	}
	return false
}

func kInARow(row []int, k int) bool {
	var count, lastEl int
	for _, el := range row {
		switch el {
		case 0:
			count = 0
			lastEl = 0
		case lastEl:
			count++
			if count == k {
				return true
			}
		default:
			count = 1
			lastEl = el
		}
	}
	return false
}

type Config struct {
	M int `json:"m"` // board dimensions are m x n
	N int `json:"n"` // board dimensions are m x n
	K int `json:"k"` // k-in-a-row to win
}

func (c *Config) valid() bool {
	if c.M <= 0 || c.N <= 0 || c.K <= 0 {
		return false
	}
	if c.K > c.M && c.K > c.N {
		return false
	}
	return true
}

type Move struct {
	X int `json:"x"`
	Y int `json:"y"`
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
	if !config.valid() {
		http.Error(w, "invalid config", http.StatusUnprocessableEntity)
		return
	}
	g.mu.Lock()
	g.New(config)
	g.mu.Unlock()
}

func (g *Game) New(config Config) {
	board := make([][]int, config.N) // create rows
	for i := range board {
		board[i] = make([]int, config.M) // for each row add the columns
	}
	g.kInARow = config.K
	g.board = board
	g.history = nil
	g.finished = false
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
