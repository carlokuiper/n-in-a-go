// https://gist.github.com/enricofoltran/10b4a980cd07cb02836f70a4ab3e72d7 looks cool
package main

import (
	"fmt"
	"net/http"

	n_in_a_go "github.com/carlokuiper/n-in-a-go"
)

func main() {
	game := &n_in_a_go.Game{}
	game.New(n_in_a_go.Config{NInARow: 3, BoardSize: 3})
	http.HandleFunc("/", game.Get)
	http.HandleFunc("/start", game.Start)
	http.HandleFunc("/turn", game.Move)
	fmt.Println("Starting server at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
