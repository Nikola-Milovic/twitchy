package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"nikolamilovic/twitchy/auth/api"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	os.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/twitchy?sslmode=disable")

	srv := api.NewServer()

	fmt.Println("Go listening on 4003")

	http.ListenAndServe(":4003", srv)
}
