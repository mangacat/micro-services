package commento

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/astaxie/beego/httplib"
)

type Request struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}
type Commenter struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Website  string `json:"website"`
	Password string `json:"password"`
}
type CommenterSession struct {
	CommenterToken string    `json:"commenterToken"`
	CommenterHex   string    `json:"commenterHex"`
	CreationDate   time.Time `json:"creationDate"`
}
type Page struct {
	Domain       string `json:"domain"`
	Path         string `json:"path"`
	IsLocked     bool   `json:"isLocked"`
	CommentCount int    `json:"commentCount"`
}

type Comments struct {
	Count   int    `json:"count"`
	Success string `json:"success"`
}

func GetCommentCount(body *Page) (*Comments, error) {
	req := httplib.Get(beego.AppConfig.String("commento") + "/api/comment/count")
	req, err := req.JSONBody(body)
	if err != nil {
		return nil, err
	}
	resp, err := req.DoRequest()
	if err != nil {
		return nil, err
	}
	comments := &Comments{}
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body) // nolint
		if err != nil {
			return nil, err
		}
		logs.Critical(string(bodyBytes))

		err = json.Unmarshal(bodyBytes, comments)
		if err != nil {
			return nil, err
		}
		if comments.Success == "false" {
			return nil, errors.New("error commento responsed with error")
		}
		return comments, nil
	}
	logs.Critical(resp)
	return nil, err

}
func CreateUser(body *Commenter) error {
	req := httplib.Post(beego.AppConfig.String("commento") + "/api/commenter/new")
	req, err := req.JSONBody(body)
	if err != nil {
		return err
	}
	resp, err := req.DoRequest()
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body) // nolint
		if err != nil {
			return err
		}
		logs.Critical(string(bodyBytes))

	}
	logs.Critical(resp)
	return err
}

func Login(body *Request) (*CommenterSession, error) {

	req := httplib.Post(beego.AppConfig.String("commento") + "/api/commenter/login")
	req, err := req.JSONBody(body)
	if err != nil {
		return nil, err
	}
	resp, err := req.DoRequest()
	if err != nil {
		return nil, err
	}
	session := &CommenterSession{}
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body) // nolint
		if err != nil {
			return nil, err
		}
		logs.Critical(string(bodyBytes))

		err = json.Unmarshal(bodyBytes, session)
		if err != nil {
			return nil, err
		}

	}
	logs.Critical(session)
	return session, err
}
