package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/carlokuiper/k-in-a-go"
)

func main() {
	config := kinago.Config{M: 3, N: 3, K: 3}
	if m, err := strconv.Atoi(os.Getenv("M")); err == nil && m != 0 {
		config.M = m
	}
	if n, err := strconv.Atoi(os.Getenv("N")); err == nil && n != 0 {
		config.N = n
	}
	if k, err := strconv.Atoi(os.Getenv("K")); err == nil && k != 0 {
		config.K = k
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
