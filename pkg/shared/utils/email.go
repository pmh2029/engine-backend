package utils

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"os"

	"gopkg.in/gomail.v2"
)

type TemplateData struct {
	Path    string
	Name    string
	To      string
	Subject string
	Url     string
}

func SendTemplateEMail(templateData TemplateData) error {
	var body bytes.Buffer
	t, err := template.ParseFiles(templateData.Path)
	// t, err := template.ParseFiles("pkg/shared/template/1.html")

	if err != nil {
		return err
	}

	t.Execute(&body, TemplateData{
		Url: templateData.Url,
	})

	m := gomail.NewMessage()
	// Set E-Mail sender
	m.SetHeader("From", os.Getenv("ENGINE_EMAIL"))

	// Set E-Mail receivers
	m.SetHeader("To", templateData.To)

	// Set E-Mail subject
	m.SetHeader("Subject", templateData.Subject)

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/html", body.String())

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("ENGINE_EMAIL"), os.Getenv("ENGINE_APP_PASS"))

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	err = d.DialAndSend(m)

	return err
}
