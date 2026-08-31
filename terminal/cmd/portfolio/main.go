package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/joho/godotenv"
	"github.com/muesli/termenv"
	"github.com/puneet/terminal-portfolio/internal/tui"
	gossh "golang.org/x/crypto/ssh"
)

func main() {
	_ = godotenv.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx); err != nil {
		log.Error("SSH server stopped", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	port := os.Getenv("SSH_PORT")
	if port == "" {
		port = "2222"
	}
	keyPath := os.Getenv("SSH_HOST_KEY_PATH")
	if keyPath == "" {
		keyPath = ".ssh/id_ed25519"
	}
	srv, err := newServer(net.JoinHostPort("0.0.0.0", port), keyPath)
	if err != nil {
		return err
	}
	listener, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}
	defer srv.Close()
	stopped := make(chan error, 1)
	go func() { stopped <- srv.Serve(listener) }()
	log.Info("Portfolio SSH listening", "port", port)
	fmt.Printf("\n  %s\n\n", lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00d4ff")).
		Render("PUNEET-OS SSH Server"))
	fmt.Printf("  Connect with: %s\n\n",
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff6600")).
			Render("ssh localhost -p "+port))
	select {
	case err := <-stopped:
		if errors.Is(err, ssh.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// Close remaining sessions after the grace period, cancelling their requests.
		if err := srv.Shutdown(shutdown); err != nil {
			return srv.Close()
		}
		return nil
	}
}

func newServer(address, keyPath string) (*ssh.Server, error) {
	srv, err := wish.NewServer(
		wish.WithAddress(address),
		wish.WithHostKeyPath(keyPath),
		wish.WithIdleTimeout(2*time.Minute),
		wish.WithMaxTimeout(15*time.Minute),
	)
	if err != nil {
		return nil, err
	}
	// Limit both unauthenticated connections and channels opened on one connection.
	connections := make(chan struct{}, 64)
	srv.ConnCallback = func(ctx ssh.Context, conn net.Conn) net.Conn {
		select {
		case connections <- struct{}{}:
			go func() { <-ctx.Done(); <-connections }()
			return conn
		default:
			return nil
		}
	}
	sessions := make(chan struct{}, 32)
	srv.ChannelHandlers = map[string]ssh.ChannelHandler{
		"session": func(srv *ssh.Server, conn *gossh.ServerConn, channel gossh.NewChannel, ctx ssh.Context) {
			select {
			case sessions <- struct{}{}:
				defer func() { <-sessions }()
				servePortfolio(channel, ctx)
			default:
				_ = channel.Reject(gossh.ResourceShortage, "Portfolio is busy; try again shortly")
			}
		},
	}
	return srv, nil
}

// servePortfolio accepts only the session requests this read-only TUI needs.
// Owning the request loop also ties the program and email context to the channel,
// rather than the longer-lived SSH connection used by multiplexed clients.
func servePortfolio(incoming gossh.NewChannel, parent context.Context) {
	channel, requests, err := incoming.Accept()
	if err != nil {
		return
	}
	ctx, cancel := context.WithCancel(parent)
	var program *tea.Program
	var finished <-chan error
	defer func() {
		cancel()
		_ = channel.Close() // Unblock input/output before waiting for the program.
		if finished != nil {
			<-finished
		}
	}()

	type ptyRequest struct {
		Term                                   string
		Width, Height, PixelWidth, PixelHeight uint32
		Modes                                  string
	}
	var pty *ptyRequest
	colorTerm := ""
	colorSet := false
	windows := make(chan tea.WindowSizeMsg, 1)
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-finished:
			finished = nil
			status := uint32(0)
			if err != nil {
				status = 1
			}
			_, _ = channel.SendRequest("exit-status", false, gossh.Marshal(struct{ Status uint32 }{status}))
			return
		case request, ok := <-requests:
			if !ok {
				return
			}
			accepted := false
			// crypto/ssh bounds packets; use a smaller cap for our few tiny requests.
			if len(request.Payload) <= 4096 {
				switch request.Type {
				case "pty-req":
					if pty == nil && program == nil {
						candidate := new(ptyRequest)
						if gossh.Unmarshal(request.Payload, candidate) == nil && len(candidate.Term) <= 64 {
							pty, accepted = candidate, true
						}
					}
				case "env":
					var env struct{ Key, Value string }
					if program == nil && !colorSet && gossh.Unmarshal(request.Payload, &env) == nil &&
						env.Key == "COLORTERM" && (env.Value == "truecolor" || env.Value == "24bit") {
						colorTerm, colorSet, accepted = env.Value, true, true
					}
				case "shell":
					if program == nil && pty != nil && len(request.Payload) == 0 {
						program = teaHandler(ctx, channel, pty.Term, colorTerm,
							tea.WindowSizeMsg{Width: int(pty.Width), Height: int(pty.Height)})
						done := make(chan error, 1)
						finished = done
						go func() { _, err := program.Run(); done <- err }()
						go func() {
							for {
								select {
								case <-ctx.Done():
									return
								case window := <-windows:
									program.Send(window)
								}
							}
						}()
						accepted = true
					}
				case "window-change":
					var window struct{ Width, Height, PixelWidth, PixelHeight uint32 }
					if pty != nil && gossh.Unmarshal(request.Payload, &window) == nil {
						pty.Width, pty.Height = window.Width, window.Height
						if program != nil {
							// Keep the latest size without blocking the SSH request loop.
							select {
							case <-windows:
							default:
							}
							windows <- tea.WindowSizeMsg{Width: int(window.Width), Height: int(window.Height)}
						}
						accepted = true
					}
				}
			}
			_ = request.Reply(accepted, nil)
		}
	}
}

func teaHandler(ctx context.Context, channel io.ReadWriter, term, colorTerm string, window tea.WindowSizeMsg) *tea.Program {
	profile := termenv.ANSI
	if term == "" || term == "dumb" {
		profile = termenv.Ascii
	} else if strings.Contains(term, "256color") {
		profile = termenv.ANSI256
	}
	if colorTerm == "truecolor" || colorTerm == "24bit" {
		profile = termenv.TrueColor
	}
	// Avoid OSC queries: terminals that do not answer must still reach the TUI.
	output := ssh.NewPtyWriter(channel)
	renderer := lipgloss.NewRenderer(output, termenv.WithProfile(profile))
	renderer.SetColorProfile(profile)
	renderer.SetHasDarkBackground(true)
	m := tui.NewModelWithRenderer(renderer)
	m.Context = ctx
	sized, _ := m.Update(window)
	return tea.NewProgram(sized, tea.WithContext(ctx), tea.WithInput(channel), tea.WithOutput(output),
		tea.WithAltScreen(), tea.WithoutSignalHandler())
}
