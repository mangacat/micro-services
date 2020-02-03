package email

import (
	"errors"
	"net/mail"
	"os"
	"strconv"

	"github.com/go-gomail/gomail"
	"github.com/matcornic/hermes"
)

type email interface {
	Email() hermes.Email
	Name() string
}

type smtpAuthentication struct {
	Server         string
	Port           int
	SenderEmail    string
	SenderIdentity string
	SMTPUser       string
	SMTPPassword   string
}

// sendOptions are options for sending an email
type sendOptions struct {
	To      string
	Subject string
}

var herme hermes.Hermes
var smtpConfig smtpAuthentication

func init() {
	herme = hermes.Hermes{
		Product: hermes.Product{
			Name: os.Getenv("APP_NAME"),
			Link: os.Getenv("APP_LINK"),
			Logo: os.Getenv("APP_IMAGE"),
		},
	}
	herme.Theme = new(hermes.Flat)

	portStr := os.Getenv("HERMES_SMTP_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(err)
	}
	smtpConfig = smtpAuthentication{
		Server:         os.Getenv("HERMES_SMTP_SERVER"),
		Port:           port,
		SenderEmail:    os.Getenv("HERMES_SENDER_EMAIL"),
		SenderIdentity: os.Getenv("HERMES_SENDER_IDENTITY"),
		SMTPPassword:   os.Getenv("HERMES_SMTP_PASSWORD"),
		SMTPUser:       os.Getenv("HERMES_SMTP_USER"),
	}
}
func SendWelcomeEmail(email string, username string, link string) error {

	tmpl := new(welcome)
	em := tmpl.Email(username, link)
	htmlBytes, err := herme.GenerateHTML(em)
	if err != nil {
		return err
	}
	txtBytes, err := herme.GeneratePlainText(em)
	if err != nil {
		return err
	}
	options := sendOptions{
		To:      email,
		Subject: "Welcome to MangaCat",
	}
	err = send(smtpConfig, options, string(htmlBytes), string(txtBytes))
	if err != nil {
		return err
	}
	return nil
}
func SendResetEmail(email string, username string, link string) error {

	tmpl := new(reset)
	em := tmpl.Email(username, link)
	htmlBytes, err := herme.GenerateHTML(em)
	if err != nil {
		return err
	}
	txtBytes, err := herme.GeneratePlainText(em)
	if err != nil {
		return err
	}
	options := sendOptions{
		To:      email,
		Subject: "Reset Password",
	}
	err = send(smtpConfig, options, string(htmlBytes), string(txtBytes))
	if err != nil {
		return err
	}
	return nil
}

// send sends the email
func send(smtpConfig smtpAuthentication, options sendOptions, htmlBody string, txtBody string) error {

	if smtpConfig.Server == "" {
		return errors.New("SMTP server config is empty")
	}
	if smtpConfig.Port == 0 {
		return errors.New("SMTP port config is empty")
	}

	if smtpConfig.SMTPUser == "" {
		return errors.New("SMTP user is empty")
	}

	if smtpConfig.SenderIdentity == "" {
		return errors.New("SMTP sender identity is empty")
	}

	if smtpConfig.SenderEmail == "" {
		return errors.New("SMTP sender email is empty")
	}

	if options.To == "" {
		return errors.New("no receiver emails configured")
	}

	from := mail.Address{
		Name:    smtpConfig.SenderIdentity,
		Address: smtpConfig.SenderEmail,
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from.String())
	m.SetHeader("To", options.To)
	m.SetHeader("Subject", options.Subject)

	m.SetBody("text/plain", txtBody)
	m.AddAlternative("text/html", htmlBody)

	d := gomail.NewPlainDialer(smtpConfig.Server, smtpConfig.Port, smtpConfig.SMTPUser, smtpConfig.SMTPPassword)

	return d.DialAndSend(m)
}
