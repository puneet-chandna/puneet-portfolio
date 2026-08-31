package main

import (
	"context"
	"io"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/ssh"
	gossh "golang.org/x/crypto/ssh"
)

func TestPublicSSHIsOnlyATUI(t *testing.T) {
	key := filepath.Join(t.TempDir(), "keys", "id_ed25519")
	srv, err := newServer("127.0.0.1:0", key)
	if err != nil {
		t.Fatal(err)
	}
	listener, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		t.Fatal(err)
	}
	defer srv.Close()
	go srv.Serve(listener)
	config := &gossh.ClientConfig{User: "visitor", HostKeyCallback: gossh.FixedHostKey(srv.HostSigners[0].PublicKey()), Timeout: 3 * time.Second}
	client, err := gossh.Dial("tcp", listener.Addr().String(), config)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	if c, err := client.Dial("tcp", "127.0.0.1:1"); err == nil {
		c.Close()
		t.Fatal("forwarding accepted")
	}
	session, err := client.NewSession()
	if err != nil {
		t.Fatal(err)
	}
	if err := session.Run("id"); err == nil {
		t.Fatal("exec request accepted")
	}
	session.Close()
	session, err = client.NewSession()
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	if err := session.RequestPty("xterm-256color", 24, 80, gossh.TerminalModes{}); err != nil {
		t.Fatal(err)
	}
	output, err := session.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	input, err := session.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := session.Shell(); err != nil {
		t.Fatal(err)
	}
	// No terminal-query replies: startup must not depend on OSC support.
	started := make(chan string, 1)
	go func() {
		var rendered strings.Builder
		b := make([]byte, 8192)
		for {
			n, err := output.Read(b)
			rendered.Write(b[:n])
			text := rendered.String()
			if strings.Contains(text, "38;5;") || strings.Contains(text, "\x1b]11;?") || err != nil {
				started <- text
				return
			}
		}
	}()
	select {
	case first := <-started:
		if strings.Contains(first, "\x1b]11;?") {
			t.Fatal("blocking background-color query")
		}
		if !strings.Contains(first, "38;5;") {
			t.Fatal("SSH terminal lost its256-color theme")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("TUI did not produce output")
	}
	go io.Copy(io.Discard, output)
	input.Write([]byte("q"))
	quit := make(chan error, 1)
	go func() { quit <- session.Wait() }()
	select {
	case err := <-quit:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("q did not quit")
	}
	// A second server with the same persisted key must keep its identity.
	restarted, err := newServer("127.0.0.1:0", key)
	if err != nil {
		t.Fatal(err)
	}
	if string(restarted.HostSigners[0].PublicKey().Marshal()) != string(srv.HostSigners[0].PublicKey().Marshal()) {
		t.Fatal("host key changed")
	}
}

func startSSHTestServer(t *testing.T) (*gossh.Client, <-chan struct{}, func() *gossh.Client) {
	t.Helper()
	srv, err := newServer("127.0.0.1:0", filepath.Join(t.TempDir(), "hostkey"))
	if err != nil {
		t.Fatal(err)
	}
	finished := make(chan struct{}, 128)
	handler := srv.ChannelHandlers["session"]
	srv.ChannelHandlers["session"] = func(s *ssh.Server, c *gossh.ServerConn, ch gossh.NewChannel, ctx ssh.Context) {
		handler(s, c, ch, ctx)
		finished <- struct{}{}
	}
	listener, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		t.Fatal(err)
	}
	go srv.Serve(listener)
	t.Cleanup(func() { _ = srv.Close() })
	dial := func() *gossh.Client {
		client, err := gossh.Dial("tcp", listener.Addr().String(), &gossh.ClientConfig{
			User: "visitor", Timeout: 3 * time.Second, HostKeyCallback: gossh.FixedHostKey(srv.HostSigners[0].PublicKey()),
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = client.Close() })
		return client
	}
	return dial(), finished, dial
}

func awaitSessionExit(t *testing.T, finished <-chan struct{}) {
	t.Helper()
	select {
	case <-finished:
	case <-time.After(3 * time.Second):
		t.Fatal("session handler did not release its resources")
	}
}

func TestSSHSessionRequestsAreBounded(t *testing.T) {
	client, finished, dial := startSSHTestServer(t)
	// The old handler leaked one admission slot per pre-shell resize, even
	// after the entire transport disconnected. Exercise more than its32 slots.
	for i := 0; i < 40; i++ {
		session, err := client.NewSession()
		if err != nil {
			t.Fatal(err)
		}
		if err := session.Setenv("IGNORED", strings.Repeat("x", 16<<10)); err == nil {
			t.Fatal("unbounded environment accepted")
		}
		if err := session.Setenv("COLORTERM", "truecolor"); err != nil {
			t.Fatal(err)
		}
		if err := session.Setenv("COLORTERM", "truecolor"); err == nil {
			t.Fatal("repeated environment accepted")
		}
		if err := session.RequestPty("xterm", 24, 80, gossh.TerminalModes{}); err != nil {
			t.Fatal(err)
		}
		if err := session.WindowChange(25, 81); err != nil {
			t.Fatal(err)
		}
		_ = session.Close()
		awaitSessionExit(t, finished)
	}
	_ = client.Close()
	client = dial()
	var held []*gossh.Session
	for i := 0; i < 32; i++ {
		session, err := client.NewSession()
		if err != nil {
			t.Fatal(err)
		}
		held = append(held, session)
	}
	if session, err := client.NewSession(); err == nil {
		session.Close()
		t.Fatal("session admission limit bypassed")
	}
	awaitSessionExit(t, finished) // The rejected channel's handler.
	for _, session := range held {
		_ = session.Close()
	}
	for range held {
		awaitSessionExit(t, finished)
	}
	// Keep the transport open while repeatedly closing actual running TUIs.
	for i := 0; i < 40; i++ {
		session, err := client.NewSession()
		if err != nil {
			t.Fatal(err)
		}
		session.Stdout, session.Stderr = io.Discard, io.Discard
		if err := session.RequestPty("xterm", 24, 80, gossh.TerminalModes{}); err != nil {
			t.Fatal(err)
		}
		if err := session.Shell(); err != nil {
			t.Fatal(err)
		}
		_ = session.Close()
		awaitSessionExit(t, finished)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestClosingSSHSessionCancelsEmail(t *testing.T) {
	t.Setenv("RESEND_API_KEY", "test-key")
	t.Setenv("RESEND_FROM", "portfolio@example.test")
	t.Setenv("RESEND_TO", "owner@example.test")
	started, stopped := make(chan context.Context, 1), make(chan struct{})
	originalTransport := http.DefaultTransport
	http.DefaultTransport = roundTripFunc(func(request *http.Request) (*http.Response, error) {
		started <- request.Context()
		<-request.Context().Done()
		close(stopped)
		return nil, request.Context().Err()
	})
	t.Cleanup(func() { http.DefaultTransport = originalTransport })
	client, finished, _ := startSSHTestServer(t)
	session, err := client.NewSession()
	if err != nil {
		t.Fatal(err)
	}
	session.Stdout, session.Stderr = io.Discard, io.Discard
	input, err := session.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := session.RequestPty("xterm", 24, 80, gossh.TerminalModes{}); err != nil {
		t.Fatal(err)
	}
	if err := session.Shell(); err != nil {
		t.Fatal(err)
	}
	// Skip boot; move three sections right; open contact and submit dummy data.
	_, err = io.WriteString(input, "\r\x1b[C\x1b[C\x1b[C\rTester\tvisitor@example.test\tHello\r")
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-started:
	case <-time.After(3 * time.Second):
		t.Fatal("contact request did not start")
	}
	_ = session.Close() // Deliberately retain the multiplexed SSH transport.
	select {
	case <-stopped:
	case <-time.After(time.Second):
		t.Fatal("email outlived its SSH session")
	}
	awaitSessionExit(t, finished)
	probe, err := client.NewSession()
	if err != nil {
		t.Fatalf("transport closed with session: %v", err)
	}
	_ = probe.Close()
}
