package email

import (
	"github.com/matcornic/hermes"
)

type reset struct {
}

func (r *reset) Name() string {
	return "reset"
}

func (r *reset) Email(username string, link string) hermes.Email {
	return hermes.Email{
		Body: hermes.Body{
			Name: username,
			Intros: []string{
				"You have received this email because a password reset request for your MangaCat account was received.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Click the button below to reset your password:",
					Button: hermes.Button{
						Color: "#DC4D2F",
						Text:  "Reset your password",
						Link:  link,
					},
				},
			},
			Outros: []string{
				"If you did not request a password reset, no further action is required on your part.",
			},
			Signature: "Thanks",
		},
	}
}
