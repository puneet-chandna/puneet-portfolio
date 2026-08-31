package tui

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func updateModel(t *testing.T, m Model, msg tea.Msg) Model {
	t.Helper()
	updated, _ := m.Update(msg)
	return updated.(Model)
}

func TestProjectsFitAndPageDownScrolls(t *testing.T) {
	m := NewModel()
	m.State = StateMain
	m.MenuIndex = 2
	m = updateModel(t, m, tea.WindowSizeMsg{Width: 80, Height: 24})

	for _, line := range strings.Split(m.View(), "\n") {
		if got := lipgloss.Width(line); got > 80 {
			t.Fatalf("rendered line is %d columns wide", got)
		}
	}
	if got := len(strings.Split(m.View(), "\n")); got > 24 {
		t.Fatalf("rendered %d rows in a 24-row terminal", got)
	}

	m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyPgDown})
	if m.Viewport.YOffset == 0 {
		t.Fatal("PageDown did not move the viewport")
	}
}

func TestSelectedProjectDetailsAreVisibleAt80Columns(t *testing.T) {
	m := NewModel()
	m.State = StateMain
	m.MenuIndex = 2
	m = updateModel(t, m, tea.WindowSizeMsg{Width: 80, Height: 24})

	project := m.Projects[m.ProjectIndex]
	if !strings.Contains(m.View(), project.URL) {
		t.Fatalf("selected project URL %q is not visible", project.URL)
	}
	if !strings.Contains(m.View(), project.Stack[0]) {
		t.Fatalf("selected project stack is not visible")
	}
}

func TestContactReopensWithFirstInputFocused(t *testing.T) {
	t.Setenv("RESEND_API_KEY", "test-key")
	t.Setenv("RESEND_FROM", "Portfolio <from@example.com>")
	t.Setenv("RESEND_TO", "to@example.com")
	m := NewModel()
	m.State = StateMain
	m.MenuIndex = 3
	m.ContactFocus = 1
	m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyEnter})
	if m.State != StateContactForm || m.ContactFocus != 0 || !m.ContactInputs[0].Focused() {
		t.Fatal("contact form did not open on the first input")
	}
	m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyEsc})
	m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyEnter})
	if m.ContactFocus != 0 || !m.ContactInputs[0].Focused() {
		t.Fatal("contact form reopened on the wrong input")
	}
}

func TestContactDoesNotClaimSuccessWithoutConfiguration(t *testing.T) {
	t.Setenv("RESEND_API_KEY", "")
	t.Setenv("RESEND_FROM", "")
	t.Setenv("RESEND_TO", "")
	m := NewModel()
	m.State = StateContactForm
	m.ContactInputs[0].SetValue("Puneet")
	m.ContactInputs[1].SetValue("puneet@example.com")
	m.ContactInputs[2].SetValue("Hello")
	m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyEnter})
	if m.State == StateContactSending || m.State == StateContactSent {
		t.Fatal("unconfigured contact form claimed it could send")
	}
}

func TestWindowSizeIsClamped(t *testing.T) {
	m := NewModel()
	m = updateModel(t, m, tea.WindowSizeMsg{Width: 1_000_000, Height: 1_000_000})
	if m.Width > 500 || m.Height > 200 {
		t.Fatalf("window size was not clamped: %dx%d", m.Width, m.Height)
	}
}

func TestContactSendEscapesHTMLAndAcceptsProviderSuccess(t *testing.T) {
	t.Setenv("RESEND_API_KEY", "test-key")
	t.Setenv("RESEND_FROM", "Portfolio <from@example.com>")
	t.Setenv("RESEND_TO", "to@example.com")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("authorization = %q", got)
		}
		var payload ResendPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if strings.Contains(payload.Html, "<Puneet>") || !strings.Contains(payload.Html, "&lt;Puneet&gt;") {
			t.Fatalf("HTML was not escaped: %q", payload.Html)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()
	oldEndpoint := resendEndpoint
	resendEndpoint = server.URL
	t.Cleanup(func() { resendEndpoint = oldEndpoint })
	resetContactLimit()

	if err := SendContactForm(context.Background(), " <Puneet> ", "person@example.com", "<hello>"); err != nil {
		t.Fatal(err)
	}
}

func TestContactSendRespectsCanceledContext(t *testing.T) {
	t.Setenv("RESEND_API_KEY", "test-key")
	t.Setenv("RESEND_FROM", "Portfolio <from@example.com>")
	t.Setenv("RESEND_TO", "to@example.com")
	resetContactLimit()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := SendContactForm(ctx, "Puneet", "person@example.com", "Hello"); err == nil {
		t.Fatal("canceled context sent a request")
	}
}

func TestFailedContactSendDoesNotConsumeRateLimit(t *testing.T) {
	t.Setenv("RESEND_API_KEY", "test-key")
	t.Setenv("RESEND_FROM", "Portfolio <from@example.com>")
	t.Setenv("RESEND_TO", "to@example.com")
	oldEndpoint := resendEndpoint
	t.Cleanup(func() { resendEndpoint = oldEndpoint })

	t.Run("transport error", func(t *testing.T) {
		failedServer := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
		resendEndpoint = failedServer.URL
		failedServer.Close()
		resetContactLimit()
		if err := SendContactForm(context.Background(), "Puneet", "person@example.com", "Hello"); err == nil || !strings.Contains(err.Error(), "unable to send right now") {
			t.Fatalf("transport error = %v", err)
		}

		successServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusAccepted)
		}))
		defer successServer.Close()
		resendEndpoint = successServer.URL
		if err := SendContactForm(context.Background(), "Puneet", "person@example.com", "Retry"); err != nil {
			t.Fatalf("transport failure blocked immediate retry: %v", err)
		}
	})

	t.Run("provider rejection", func(t *testing.T) {
		attempts := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			attempts++
			if attempts == 1 {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			w.WriteHeader(http.StatusAccepted)
		}))
		defer server.Close()
		resendEndpoint = server.URL
		resetContactLimit()

		if err := SendContactForm(context.Background(), "Puneet", "person@example.com", "Hello"); err == nil || !strings.Contains(err.Error(), "provider could not accept") {
			t.Fatalf("provider rejection = %v", err)
		}
		if err := SendContactForm(context.Background(), "Puneet", "person@example.com", "Retry"); err != nil {
			t.Fatalf("provider rejection blocked immediate retry: %v", err)
		}
		if err := SendContactForm(context.Background(), "Puneet", "person@example.com", "Again"); err == nil || !strings.Contains(err.Error(), "please wait") {
			t.Fatalf("accepted message rate limit = %v", err)
		}
		if attempts != 2 {
			t.Fatalf("rate-limited request reached provider; attempts = %d", attempts)
		}
	})
}

func TestFailedSendReturnsToFormWithValues(t *testing.T) {
	m := NewModel()
	m.State = StateContactSending
	m.ContactInputs[0].SetValue("Puneet")
	m.ContactInputs[1].SetValue("person@example.com")
	m.ContactInputs[2].SetValue("Hello")
	m = updateModel(t, m, sendResultMsg{err: context.Canceled})
	if m.State != StateContactForm || m.ContactInputs[2].Value() != "Hello" || m.ContactError == "" {
		t.Fatal("failed send did not restore the form for retry")
	}
}

func TestSuccessfulRetryClearsContactError(t *testing.T) {
	t.Setenv("RESEND_API_KEY", "test-key")
	t.Setenv("RESEND_FROM", "Portfolio <from@example.com>")
	t.Setenv("RESEND_TO", "to@example.com")
	m := NewModel()
	m.State = StateContactSending
	m.MenuIndex = 3
	previousError := "email provider could not accept the message; please try again"
	m.ContactError = previousError
	m = updateModel(t, m, sendResultMsg{})
	contactInfo := m.renderContactInfo()
	if m.State != StateContactSent || strings.Contains(contactInfo, previousError) || !strings.Contains(contactInfo, "Press [Enter] to open the message form") {
		t.Fatal("contact page retained the previous delivery error after a successful retry")
	}
}

func TestContactSendingStillAllowsCtrlC(t *testing.T) {
	m := NewModel()
	m.State = StateContactSending
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if updated.(Model).State != StateContactSending {
		t.Fatal("Ctrl+C changed the sending state before quitting")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("Ctrl+C did not quit while sending")
	}
}

func TestContactRejectsControlCharactersAndRateLimits(t *testing.T) {
	if _, _, _, err := validateContact("Puneet\r\nBcc: other@example.com", "person@example.com", "Hello"); err == nil {
		t.Fatal("header injection was accepted")
	}
	resetContactLimit()
	if _, err := reserveContactAttempt(); err != nil {
		t.Fatal(err)
	}
	if _, err := reserveContactAttempt(); err == nil {
		t.Fatal("back-to-back contact attempts bypassed the rate limit")
	}
}

func TestDataFieldsDriveRenderedPages(t *testing.T) {
	m := NewModel()
	m.State = StateMain
	m.Bio = "bio from data"
	m = updateModel(t, m, tea.WindowSizeMsg{Width: 80, Height: 24})
	if !strings.Contains(m.View(), "bio from data") {
		t.Fatal("about page ignored Model.Bio")
	}
	m.MenuIndex = 1
	m.Experience = "experience from data"
	m.refreshViewport(true)
	if !strings.Contains(m.View(), "experience from data") {
		t.Fatal("experience page ignored Model.Experience")
	}
}

func TestMainLayoutKeepsFooterBordersAndInspectorURLVisible(t *testing.T) {
	for _, size := range [][2]int{{80, 24}, {120, 40}} {
		t.Run("terminal", func(t *testing.T) {
			m := NewModel()
			m.State = StateMain
			m.MenuIndex = 2
			m = updateModel(t, m, tea.WindowSizeMsg{Width: size[0], Height: size[1]})
			view := m.View()
			if got := lipgloss.Height(view); got != size[1] {
				t.Fatalf("rendered %d rows in a %d-row terminal", got, size[1])
			}
			if !strings.Contains(view, "[?] Help") {
				t.Fatalf("footer controls are not visible:\n%s", view)
			}
			lines := strings.Split(view, "\n")
			footer := -1
			for i, line := range lines {
				if strings.Contains(line, "[?] Help") {
					footer = i
					break
				}
			}
			if footer < 1 || !strings.Contains(lines[footer-1], "╰") {
				t.Fatal("pane bottom border is not visible above the footer")
			}
			if size[0] == 120 && !strings.Contains(strings.ReplaceAll(view, "\n", ""), "github.com/puneet-chandna/0DTE-dealer-gamma") {
				t.Fatal("inspector did not render the complete selected-project URL")
			}
		})
	}
}

func TestSmallContactAndHelpKeepActiveContentVisible(t *testing.T) {
	t.Setenv("RESEND_API_KEY", "test-key")
	t.Setenv("RESEND_FROM", "Portfolio <from@example.com>")
	t.Setenv("RESEND_TO", "to@example.com")
	m := NewModel()
	m.State = StateMain
	m.MenuIndex = 3
	m = updateModel(t, m, tea.WindowSizeMsg{Width: 60, Height: 20})
	m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyEnter})
	m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyTab})
	m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyTab})
	m.ContactInputs[m.ContactFocus].SetValue("visible-message")
	m.ContactError = "enter a valid email address"
	m = updateModel(t, m, tea.WindowSizeMsg{Width: 40, Height: 15})
	view := m.View()
	for _, text := range []string{"MESSAGE", "visible-message", "enter a valid email address", "[Tab]", "[Esc]"} {
		if !strings.Contains(view, text) {
			t.Fatalf("small contact form hides %q", text)
		}
	}

	m.State = StateMain
	m.ShowHelp = true
	view = m.View()
	for _, text := range []string{"PgUp/PgDn", "? Close", "q Exit"} {
		if !strings.Contains(view, text) {
			t.Fatalf("small help hides %q", text)
		}
	}
}

func TestShortWideProjectSelectionsKeepFooterWithoutInspector(t *testing.T) {
	for _, size := range [][2]int{{120, 24}, {120, 20}} {
		m := NewModel()
		m.State = StateMain
		m.MenuIndex = 2
		m = updateModel(t, m, tea.WindowSizeMsg{Width: size[0], Height: size[1]})
		if _, inspector, _, _ := m.paneLayout(); inspector != 0 {
			t.Fatalf("short %dx%d layout kept an inspector that cannot fit", size[0], size[1])
		}
		for i := range m.Projects {
			m.ProjectIndex = i
			m.refreshViewport(true)
			if !strings.Contains(m.View(), "[?] Help") {
				t.Fatalf("footer hidden at %dx%d selecting %q", size[0], size[1], m.Projects[i].Name)
			}
		}
	}
}

func TestSmallContactLongInputAndProviderErrorKeepControlsVisible(t *testing.T) {
	m := NewModel()
	m.State = StateContactForm
	m.focusContact(2)
	m = updateModel(t, m, tea.WindowSizeMsg{Width: 40, Height: 15})
	for _, r := range strings.Repeat("a", 100) + "FINAL" {
		m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	m.ContactError = "email provider could not accept the message; please try again"
	view := m.View()
	for _, text := range []string{"FINAL", "email provider could not accept", "[Esc]"} {
		if !strings.Contains(view, text) {
			t.Fatalf("small contact form hides %q", text)
		}
	}
}

func resetContactLimit() {
	contactMu.Lock()
	defer contactMu.Unlock()
	lastContactSend = time.Time{}
}

func TestOriginalThemeSurvivesDataAndDeliveryFixes(t *testing.T) {
	renderer := lipgloss.NewRenderer(io.Discard, termenv.WithProfile(termenv.TrueColor))
	renderer.SetColorProfile(termenv.TrueColor)
	renderer.SetHasDarkBackground(true)
	m := NewModelWithRenderer(renderer)
	for _, text := range []string{"OPERATOR: PUNEET CHANDNA", "Cryptography", "\x1b["} {
		if !strings.Contains(m.renderColorizedBio(), text) {
			t.Fatalf("bio lost %q", text)
		}
	}
	if !strings.Contains(m.renderExperience(), "\x1b[") || !strings.Contains(m.renderExperience(), "Apoliums") {
		t.Fatal("experience lost styling or updated content")
	}
	m.State = StateContactSending
	before := m.View()
	for range 100 {
		updated, cmd := m.Update(sendingTickMsg(time.Now()))
		m = updated.(Model)
		if m.State != StateContactSending || cmd == nil {
			t.Fatal("animation incorrectly completed delivery or stopped")
		}
	}
	if before == m.View() || !strings.Contains(m.View(), "TRANSMISSION IN PROGRESS") {
		t.Fatal("transmission animation was lost")
	}
	m = updateModel(t, m, sendResultMsg{})
	if m.State != StateContactSent {
		t.Fatal("provider result did not complete transmission")
	}
	_, cmd := m.Update(sendingTickMsg(time.Now()))
	if cmd != nil {
		t.Fatal("animation continued after provider result")
	}
}
