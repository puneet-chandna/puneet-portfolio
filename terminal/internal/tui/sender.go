package tui

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"sync"
	"time"
	"unicode"
)

const (
	resendAPIURL       = "https://api.resend.com/emails"
	contactMaxName     = 50
	contactMaxEmail    = 100
	contactMaxMessage  = 500
	contactMinInterval = 10 * time.Second
)

var (
	resendEndpoint  = resendAPIURL
	contactMu       sync.Mutex
	lastContactSend time.Time
)

type ResendPayload struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Html    string   `json:"html"`
	ReplyTo string   `json:"reply_to"`
}

func ContactFormConfigured() bool {
	return os.Getenv("RESEND_API_KEY") != "" && os.Getenv("RESEND_FROM") != "" && os.Getenv("RESEND_TO") != ""
}

func validateContact(name, email, message string) (string, string, string, error) {
	name, email, message = strings.TrimSpace(name), strings.TrimSpace(email), strings.TrimSpace(message)
	if name == "" || message == "" {
		return "", "", "", fmt.Errorf("name and message are required")
	}
	if len([]rune(name)) > contactMaxName || len([]rune(email)) > contactMaxEmail || len([]rune(message)) > contactMaxMessage {
		return "", "", "", fmt.Errorf("message is too long")
	}
	if strings.IndexFunc(name+email+message, unicode.IsControl) >= 0 {
		return "", "", "", fmt.Errorf("control characters are not allowed")
	}
	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Address != email {
		return "", "", "", fmt.Errorf("enter a valid email address")
	}
	return name, email, message, nil
}

func reserveContactAttempt() (time.Time, error) {
	contactMu.Lock()
	defer contactMu.Unlock()
	if wait := contactMinInterval - time.Since(lastContactSend); !lastContactSend.IsZero() && wait > 0 {
		return time.Time{}, fmt.Errorf("please wait %d seconds before trying again", int(wait.Seconds())+1)
	}
	// ponytail: one process-wide limit; use shared storage if multiple replicas need a global quota.
	lastContactSend = time.Now()
	return lastContactSend, nil
}

func rollbackContactAttempt(attempt time.Time) {
	contactMu.Lock()
	defer contactMu.Unlock()
	if lastContactSend == attempt {
		lastContactSend = time.Time{}
	}
}

func SendContactForm(ctx context.Context, name, email, message string) error {
	name, email, message, err := validateContact(name, email, message)
	if err != nil {
		return err
	}
	if !ContactFormConfigured() {
		return fmt.Errorf("contact form is unavailable until email delivery is configured")
	}
	payload := ResendPayload{
		From:    os.Getenv("RESEND_FROM"),
		To:      []string{os.Getenv("RESEND_TO")},
		Subject: "New message via Terminal Portfolio",
		Html: fmt.Sprintf("<h2>New Contact from Terminal Portfolio</h2><p><strong>Name:</strong> %s</p><p><strong>Email:</strong> %s</p><p><strong>Message:</strong></p><p>%s</p>",
			html.EscapeString(name), html.EscapeString(email), html.EscapeString(message)),
		ReplyTo: email,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("could not prepare message")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, resendEndpoint, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("could not prepare message")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("RESEND_API_KEY"))

	attempt, err := reserveContactAttempt()
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		rollbackContactAttempt(attempt)
		return fmt.Errorf("unable to send right now; please try again")
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		rollbackContactAttempt(attempt)
		return fmt.Errorf("email provider could not accept the message; please try again")
	}
	if _, err := io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20)); err != nil {
		return fmt.Errorf("unable to read delivery response")
	}
	return nil
}
