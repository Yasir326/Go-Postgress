package main

import (
	"go-postgres/go-postgres/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	r := router.Router()

	fmt.Println("Server started on port 8080:")

	log.Fatal(http.ListenAndServe(":8080", r))
}