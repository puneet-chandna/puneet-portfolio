package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/joho/godotenv"
	"github.com/puneet/terminal-portfolio/internal/tui"
)

const host = "0.0.0.0"

// getPort returns the SSH port from environment variable or default
func getPort() string {
	if port := os.Getenv("SSH_PORT"); port != "" {
		return port
	}
	return "22" // Default to port 22 for production
}

func main() {
	// Load .env file (ignore error if not found - will use system env vars)
	_ = godotenv.Load()

	port := getPort()
	srv, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
		os.Exit(1)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Info("Starting SSH server", "host", host, "port", port)
	fmt.Printf("\n  %s\n\n", lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00d4ff")).
		Render("PUNEET-OS SSH Server"))
	fmt.Printf("  Connect with: %s\n\n",
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff6600")).
			Render("ssh localhost -p "+port))

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()
	m := tui.NewModel()
	m.Width = pty.Window.Width
	m.Height = pty.Window.Height

	return m, []tea.ProgramOption{tea.WithAltScreen()}
}
