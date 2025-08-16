package kinago

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Move struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Game struct {
	K        int     `json:"k"`
	Board    [][]int `json:"board"`
	History  []Move  `json:"history"`
	Finished bool    `json:"finished"`
	mu       sync.Mutex
}

func (g *Game) New(config Config) {
	board := make([][]int, config.N) // create rows
	for i := range board {
		board[i] = make([]int, config.M) // for each row add the columns
	}
	g.K = config.K
	g.Board = board
	g.History = nil
	g.Finished = false
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
	g.writeResponse(w)
}
func (g *Game) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
	response, err := json.Marshal(g)
	if err != nil {
		http.Error(w, "cannot marshall response", http.StatusInternalServerError)
	}
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, "cannot write response", http.StatusInternalServerError)
	}
	g.writeResponse(w)
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
	if err := g.update(move, nextValue); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
	g.mu.Unlock()
	g.writeResponse(w)
}

func (g *Game) writeResponse(w http.ResponseWriter) {
	response, err := json.Marshal(g)
	if err != nil {
		http.Error(w, "cannot marshall response", http.StatusInternalServerError)
	}
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, "cannot write response", http.StatusInternalServerError)
	}
}

func (g *Game) nextValue() int {
	if g.History == nil {
		return 1
	}
	previousMove := g.History[len(g.History)-1]
	switch g.Board[previousMove.Y][previousMove.X] {
	case 1:
		return 2
	case 2:
		return 1
	default:
		return 0
	}
}

func (g *Game) update(move Move, nextValue int) error {
	m, n := mxn(g.Board)
	if move.X < 0 || move.Y < 0 || move.X >= n || move.Y >= m {
		return fmt.Errorf("move invalid")
	}
	if g.History == nil {
		g.Board[move.Y][move.X] = nextValue
		g.History = []Move{move}
		return nil
	}
	previousMove := g.History[len(g.History)-1]
	if previousMove == move {
		return nil
	}
	if g.Finished {
		return fmt.Errorf("game already finished")
	}
	if g.Board[move.Y][move.X] != 0 {
		return fmt.Errorf("move not free")
	}
	g.Board[move.Y][move.X] = nextValue
	g.History = append(g.History, move)
	g.Finished = finished(g.Board, g.K)
	return nil
}

func finished(board [][]int, k int) bool {
	m, n := mxn(board)
	for _, row := range board {
		if kInARow(row, k) {
			return true
		}
	}
	for i := range m {
		column := make([]int, n)
		for j, row := range board {
			column[j] = row[i]
		}
		if kInARow(column, k) {
			return true
		}
	}
	// (off) diagonal
	for i := -(n - 1); i < m; i++ {
		x := i
		y := 0
		offDiagonal := make([]int, 0, n)
		for range n {
			if x >= 0 && y >= 0 && x < m && y < n {
				offDiagonal = append(offDiagonal, board[y][x])
			}
			x++
			y++
		}
		if kInARow(offDiagonal, k) {
			return true
		}
	}
	// (off) anti diagonal
	for i := 0; i < m+n; i++ {
		x := 0
		y := i
		offAntiDiagonal := make([]int, 0, m)
		for range m {
			if x >= 0 && y >= 0 && x < m && y < n {
				offAntiDiagonal = append(offAntiDiagonal, board[y][x])
			}
			x++
			y--
		}
		if kInARow(offAntiDiagonal, k) {
			return true
		}
	}
	return false
}

func mxn(board [][]int) (int, int) {
	n := len(board)
	if n == 0 {
		return 0, 0
	}
	return len(board[0]), n
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
