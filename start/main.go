package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	n_in_a_go "github.com/carlokuiper/n-in-a-go"
)

func main() {
	config := n_in_a_go.Config{
		NInARow:   3,
		BoardSize: 3,
	}
	body, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	r, err := http.NewRequest("POST", "http://localhost:8080/start", bytes.NewBuffer(body))
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
