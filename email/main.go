// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/micro/go-micro/web"
// //  	 k8s "github.com/mangacat/micro-services/utils/k8s"
// //   config "github.com/mangacat/micro-services/utils/config"
// )

// func main() {

// 	service := k8s.NewService(
// 		web.Name("email"),
//   )
//  	config := config.NewConfig()

// 	service.Init()

// 	service.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

// 		// if r.Method == "POST" {
// 		// 	r.ParseForm()

// 		// 	name := r.Form.Get("name")
// 		// 	if len(name) == 0 {
// 		// 		name = "World"
// 		// 	}

// 		// 	// cl := hello.NewSayService("go.micro.srv.greeter", client.DefaultClient)
// 		// 	// rsp, err := cl.Hello(context.Background(), &hello.Request{
// 		// 	// 	Name: name,
// 		// 	// })

// 		// 	if err != nil {
// 		// 		http.Error(w, err.Error(), 500)
// 		// 		return
// 		// 	}

// 		// 	w.Write([]byte(`<html><body><h1>` + rsp.Msg + `</h1></body></html>`))
// 		// 	return
// 		// }

// 		// fmt.Fprint(w, `<html><body><h1>Enter Name<h1><form method=post><input name=name type=text /></form></body></html>`)
// 	})

// 	// Reset ...
// 	// @Title Get
// 	// @Description create User
// 	// @Success 201 {object} models.Users
// 	// @Failure 403 body is empty
// 	// @router /reset [post]
// 	// func (c *AuthController) Reset() {
// 	// 	c.Ctx.Request.Header.Get("Authorization")
// 	// 	var v models.Auth
// 	// 	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
// 	// 		query := make(map[string]string)
// 	// 		query["email__exact"] = v.Email
// 	// 		ml, err := models.GetAllUsers(query, nil, nil, nil, 0, 1)
// 	// 		if len(ml) == 0 {
// 	// 			c.Data["json"] = err.Error()
// 	// 			c.Ctx.Output.SetStatus(401)
// 	// 			c.ServeJSON()
// 	// 			return
// 	// 		}
// 	// 		if err != nil {
// 	// 			c.Data["json"] = err.Error()
// 	// 			c.Ctx.Output.SetStatus(401)
// 	// 			c.ServeJSON()
// 	// 			return
// 	// 		}
// 	// 		user := ml[0].(models.Users)

// 	// 		token := passwordreset.NewToken(v.Email, 12*time.Hour, []byte(user.PasswordHash), []byte(beego.AppConfig.String("secret")))

// 	// 		url := fmt.Sprintf("https://manga.cat/reset?token=%s", token)
// 	// 		err = email.SendResetEmail(v.Email, v.Username, url)
// 	// 		if err != nil {
// 	// 			panic(err)
// 	// 		}
// 	// 	} else {

// 	// 		c.Data["json"] = err.Error()
// 	// 		c.Ctx.Output.SetStatus(401)
// 	// 		c.ServeJSON()
// 	// 		return
// 	// 	}

// 	// }

// 	if err := service.Init(); err != nil {
// 		log.Fatal(err)
// 	}

// 	if err := service.Run(); err != nil {
// 		log.Fatal(err)
// 	}
// }

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	events "github.com/mangacat/micro-services/event-structs/email"
	"github.com/mangacat/micro-services/utils/config"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		emailstruct := &events.UserCreate{}
		err = json.Unmarshal(body, &emailstruct)
		if err != nil {
			panic(err)
		}
	}
	fmt.Fprintf(w, "Welcome, %!", r.URL.Path[1:])
}
func main() {
	log.Print("helloworld: starting server...")

	config := config.NewConfig()
	// if err != nil {
	// panic(err)
	// }
	fmt.Println(config)
	http.HandleFunc("/", handler)
	// http.HandleFunc("/about/", about)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("helloworld: listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
