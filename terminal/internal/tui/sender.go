package tui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

const resendAPIURL = "https://api.resend.com/emails"

type ResendPayload struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Html    string   `json:"html"`
	ReplyTo string   `json:"reply_to"`
}

type ResendResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func SendContactForm(name, email, message string) error {
	apiKey := os.Getenv("RESEND_API_KEY")

	// Log for debugging
	if apiKey == "" {
		log.Error("RESEND_API_KEY environment variable is not set")
		return fmt.Errorf("RESEND_API_KEY not configured")
	}
	log.Info("Sending contact form via Resend", "name", name, "email", email)

	// Create HTML email body
	htmlBody := fmt.Sprintf(`
		<h2>New Contact from Terminal Portfolio</h2>
		<p><strong>Name:</strong> %s</p>
		<p><strong>Email:</strong> %s</p>
		<p><strong>Message:</strong></p>
		<p>%s</p>
		<hr>
		<p><em>Sent from SSH Terminal Portfolio</em></p>
	`, name, email, message)

	// Resend requires a verified domain or use onboarding@resend.dev for testing
	payload := ResendPayload{
		From:    "Terminal Portfolio <onboarding@resend.dev>",
		To:      []string{"puneetchandna7@gmail.com"},
		Subject: fmt.Sprintf("New message from %s via Terminal Portfolio", name),
		Html:    htmlBody,
		ReplyTo: email,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Error("Failed to marshal JSON", "error", err)
		return err
	}

	req, err := http.NewRequest("POST", resendAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error("Failed to create request", "error", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("HTTP request failed", "error", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed to read response body", "error", err)
		return err
	}

	var apiResp ResendResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		log.Error("Failed to parse API response", "error", err, "body", string(body))
		return err
	}

	if resp.StatusCode != 200 {
		log.Error("Resend API returned error", "status", resp.StatusCode, "message", apiResp.Message, "body", string(body))
		return fmt.Errorf("API error: %s", apiResp.Message)
	}

	log.Info("Contact form sent successfully", "id", apiResp.ID)
	return nil
}
