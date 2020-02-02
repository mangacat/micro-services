package email

import (
	"github.com/matcornic/hermes"
)

type welcome struct {
}

func (w *welcome) Name() string {
	return "welcome"
}

func (w *welcome) Email(username string, link string) hermes.Email {
	return hermes.Email{
		Body: hermes.Body{
			Name: username,
			Intros: []string{
				"Welcome to MangaCat! We're very excited to have you on board.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "To get started with MangaCat, please click here:",
					Button: hermes.Button{
						Text: "Confirm your account",
						Link: link,
					},
				},
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}
}
