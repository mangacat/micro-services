package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dchest/passwordreset"
	events "github.com/mangacat/micro-services/event-structs/email"
	"github.com/mangacat/micro-services/utils/email"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		emailstruct := &events.UserCreate{}
		err = json.Unmarshal(body, &emailstruct)
		if err != nil {
			log.Println(err)
		}
		if emailstruct.Event.Data.New.Email == "" {
			log.Println(err)
		}

		token := passwordreset.NewToken(emailstruct.Event.Data.New.Email.(string), 11*time.Hour, []byte(emailstruct.Event.Data.New.SecretToken), []byte(os.Getenv("secret")))

		url := fmt.Sprintf("%s/reset?token=%s", os.Getenv("MANGA_DOMAIN_URL"), token)
		err = email.SendWelcomeEmail(emailstruct.Event.Data.New.Email.(string), emailstruct.Event.Data.New.DisplayName, url)
		if err != nil {
			log.Println(err)
		}
	}
}
func main() {
	log.Print("helloworld: starting server...")

	// http.HandleFunc("/verify", handler)
	// http.HandleFunc("/reset", handler)
	http.HandleFunc("/webhook", handler)
	// http.HandleFunc("/about/", about)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("helloworld: listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
