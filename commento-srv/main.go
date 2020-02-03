package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mangacat/micro-services/utils/commento"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print("helloworld: received a request")
	target := os.Getenv("TARGET")
	if target == "" {
		target = "World"
	}
	fmt.Fprintf(w, "Hello %s!\n", target)
	k := &commento.Commenter{
		// Email:    m.Email,
		// Name:     m.Username,
		// Password: v.Password,
		// Website:  fmt.Sprintf("https://manga.cat/user/%d", id),
	}
	err = commento.CreateUser(k)
	if err != nil {
		panic(err)
	}
}

func main() {
	log.Print("helloworld: starting server...")

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("helloworld: listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
