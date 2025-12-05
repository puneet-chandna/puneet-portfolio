package tui

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

const web3FormsURL = "https://api.web3forms.com/submit"
const web3FormsKey = "4e4fb280-8c90-4e14-b8b8-2f11286cfdc8"

type ContactPayload struct {
	AccessKey string `json:"access_key"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Message   string `json:"message"`
	Subject   string `json:"subject"`
}

func SendContactForm(name, email, message string) error {
	payload := ContactPayload{
		AccessKey: web3FormsKey,
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
