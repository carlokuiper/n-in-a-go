package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	n_in_a_go "github.com/carlokuiper/n-in-a-go"
)

func main() {
	move := n_in_a_go.Move{
		X: 2,
		Y: 0,
	}
	body, err := json.Marshal(move)
	if err != nil {
		panic(err)
	}
	r, err := http.NewRequest("POST", "http://localhost:8080/turn", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	response, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer func() { _ = response.Body.Close() }()
	fmt.Println(response.StatusCode)
	//err = json.NewDecoder(response.Body).Decode(obj)
}
