package tui

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

const web3FormsURL = "https://api.web3forms.com/submit"

type ContactPayload struct {
	AccessKey string `json:"access_key"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Message   string `json:"message"`
	Subject   string `json:"subject"`
}

func SendContactForm(name, email, message string) error {
	payload := ContactPayload{
		AccessKey: os.Getenv("WEB3FORMS_KEY"),
		Name:      name,
		Email:     email,
		Message:   message,
		Subject:   "New message from Terminal Portfolio",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(web3FormsURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
