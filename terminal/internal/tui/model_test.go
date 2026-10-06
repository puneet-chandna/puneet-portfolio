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
	"github.com/charmbracelet/x/ansi"
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

func TestScrollingReadsCurrentSectionWithoutSkippingExperience(t *testing.T) {
	for _, tab := range []int{0, 1, 2} {
		m := NewModel()
		m.State = StateMain
		m.MenuIndex = tab
		m = updateModel(t, m, tea.WindowSizeMsg{Width: 80, Height: 24})
		m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyDown})
		if m.MenuIndex != tab || m.ProjectIndex != 0 || m.Viewport.YOffset != 1 {
			t.Fatalf("down changed section or project instead of scrolling tab %d", tab)
		}
		m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
		if m.MenuIndex != tab || m.ProjectIndex != 0 || m.Viewport.YOffset != 2 {
			t.Fatalf("j did not scroll tab %d", tab)
		}
		m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyUp})
		if m.Viewport.YOffset != 1 {
			t.Fatal("up did not scroll back")
		}
		m = updateModel(t, m, tea.MouseMsg{Action: tea.MouseActionPress, Button: tea.MouseButtonWheelDown})
		if m.MenuIndex != tab || m.ProjectIndex != 0 || m.Viewport.YOffset <= 1 {
			t.Fatal("mouse wheel did not scroll the current content")
		}
	}
	m := NewModel()
	m.State = StateMain
	m = updateModel(t, m, tea.WindowSizeMsg{Width: 80, Height: 24})
	m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyRight})
	if m.ActiveTab() != "experience" {
		t.Fatal("right skipped Experience")
	}
	m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyRight})
	m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	if m.ProjectIndex != 1 || m.ActiveTab() != "projects" {
		t.Fatal("n did not select the next project")
	}
	m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	if m.ProjectIndex != 0 {
		t.Fatal("p did not select the previous project")
	}
}

func TestOverflowIsVisibleAndEverySectionCanBeRead(t *testing.T) {
	for _, size := range [][2]int{{32, 15}, {40, 15}, {80, 24}, {160, 45}} {
		for _, tab := range []int{0, 1} {
			m := NewModel()
			m.State = StateMain
			m.MenuIndex = tab
			m = updateModel(t, m, tea.WindowSizeMsg{Width: size[0], Height: size[1]})
			if !strings.Contains(m.View(), "More below") || !strings.Contains(m.View(), "↑↓") {
				t.Fatalf("%dx%d tab %d did not explain hidden content", size[0], size[1], tab)
			}
			var read strings.Builder
			for !m.Viewport.AtBottom() {
				read.WriteString(m.Viewport.View())
				m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyPgDown})
			}
			read.WriteString(m.Viewport.View())
			if !strings.Contains(m.View(), "More above") || strings.Contains(m.View(), "More below") {
				t.Fatal("bottom of content has an incorrect overflow hint")
			}
			last := "elegant."
			if tab == 1 {
				last = "validation"
			}
			if !strings.Contains(read.String(), last) {
				t.Fatalf("could not read the end of tab %d at %dx%d", tab, size[0], size[1])
			}
			view := m.View()
			if lipgloss.Width(view) > size[0] || lipgloss.Height(view) > size[1] || !strings.Contains(view, "[q] Exit") {
				t.Fatalf("scrolling broke the %dx%d layout", size[0], size[1])
			}
		}
	}
}

func TestHelpDoesNotNavigateBehindOverlay(t *testing.T) {
	m := NewModel()
	m.State = StateMain
	m.MenuIndex = 2
	m = updateModel(t, m, tea.WindowSizeMsg{Width: 80, Height: 24})
	m.ShowHelp = true
	for _, key := range []tea.KeyType{tea.KeyDown, tea.KeyRight, tea.KeyPgDown, tea.KeyEnter} {
		m = updateModel(t, m, tea.KeyMsg{Type: key})
		if m.MenuIndex != 2 || m.ProjectIndex != 0 || m.Viewport.YOffset != 0 || m.State != StateMain {
			t.Fatal("help allowed navigation behind the overlay")
		}
	}
	m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyEsc})
	if m.ShowHelp {
		t.Fatal("Esc did not close help")
	}
}

func TestResizeKeepsReadingPositionAndClampsAtBottom(t *testing.T) {
	m := NewModel()
	m.State = StateMain
	m.Experience = strings.Repeat("readable work history\n", 50)
	m.MenuIndex = 1
	m = updateModel(t, m, tea.WindowSizeMsg{Width: 80, Height: 24})
	m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyPgDown})
	offset := m.Viewport.YOffset
	m = updateModel(t, m, tea.WindowSizeMsg{Width: 80, Height: 20})
	if m.Viewport.YOffset != offset {
		t.Fatal("resize jumped to the top of the work history")
	}
	m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyEnd})
	if !m.Viewport.AtBottom() {
		t.Fatal("End did not scroll to the bottom")
	}
	m = updateModel(t, m, tea.WindowSizeMsg{Width: 80, Height: 40})
	if !m.Viewport.AtBottom() || m.Viewport.PastBottom() {
		t.Fatal("growing the terminal did not clamp the reading position")
	}
	m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyHome})
	if !m.Viewport.AtTop() {
		t.Fatal("Home did not scroll to the top")
	}
}

func TestIdleMainDoesNotScheduleAnimation(t *testing.T) {
	m := NewModel()
	m.State = StateMain
	_, cmd := m.Update(tickMsg(time.Now()))
	if cmd != nil {
		t.Fatal("idle main screen scheduled another animation tick")
	}
}

func TestNarrowMainKeepsSectionAndExitVisible(t *testing.T) {
	for _, width := range []int{32, 40, 60} {
		m := NewModel()
		m.State = StateMain
		m.MenuIndex = 2
		m = updateModel(t, m, tea.WindowSizeMsg{Width: width, Height: 15})
		view := m.View()
		for _, text := range []string{"/projects", "[?] Help", "[q] Exit"} {
			if !strings.Contains(view, text) {
				t.Fatalf("%d-column view hides %q:\n%s", width, text, view)
			}
		}
		if lipgloss.Width(m.GetHeader()) > width || lipgloss.Width(m.GetFooter()) > width {
			t.Fatalf("header or footer wraps beyond %d columns", width)
		}
	}
}

func TestNarrowBootKeepsIdentityAndProgressVisible(t *testing.T) {
	m := NewModel()
	m = updateModel(t, m, tea.WindowSizeMsg{Width: 32, Height: 15})
	if !strings.Contains(m.View(), "PUNEET-OS") {
		t.Fatal("narrow boot screen clipped its identity")
	}
	m.State = StateLoading
	m.BootProgress = 100
	if !strings.Contains(m.View(), "100%") {
		t.Fatal("narrow loading screen clipped its progress")
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

func TestContactAvailabilityIsVisibleBeforeScrolling(t *testing.T) {
	for _, configured := range []bool{false, true} {
		value := ""
		if configured {
			value = "configured"
		}
		t.Setenv("RESEND_API_KEY", value)
		t.Setenv("RESEND_FROM", value)
		t.Setenv("RESEND_TO", value)
		m := NewModel()
		m.State = StateMain
		m.MenuIndex = 3
		m = updateModel(t, m, tea.WindowSizeMsg{Width: 40, Height: 15})
		want := "Message form"
		if configured {
			want = "Press [Enter]"
		}
		if !strings.Contains(m.View(), want) {
			t.Fatal("contact page hid form availability below the fold")
		}
	}
}

func TestSmallContactSuccessKeepsReturnInstructionVisible(t *testing.T) {
	m := NewModel()
	m.State = StateContactSent
	m = updateModel(t, m, tea.WindowSizeMsg{Width: 32, Height: 15})
	view := m.View()
	for _, text := range []string{"TRANSMISSION COMPLETE", "email provider.", "Press any key"} {
		if !strings.Contains(view, text) {
			t.Fatalf("small delivery confirmation hides %q", text)
		}
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
	if updated.(Model).State != StateExiting {
		t.Fatal("Ctrl+C did not start the disconnect animation while sending")
	}
	if _, ok := cmd().(tea.QuitMsg); ok {
		t.Fatal("Ctrl+C skipped the disconnect animation while sending")
	}
}

func TestDisconnectPreservesFormTypingAndResize(t *testing.T) {
	m := NewModel()
	m.State = StateContactForm
	m.focusContact(0)
	m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if m.State != StateContactForm || m.ContactInputs[0].Value() != "q" {
		t.Fatal("q stopped working as form text")
	}
	m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyCtrlC})
	firstFrame := m.View()
	for range 12 {
		m = updateModel(t, m, exitTickMsg{})
	}
	if m.View() == firstFrame {
		t.Fatal("disconnect screen did not animate")
	}
	m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	m = updateModel(t, m, tea.WindowSizeMsg{Width: 40, Height: 16})
	if m.State != StateExiting || m.ExitFrame != 12 || lipgloss.Width(m.View()) > 40 || lipgloss.Height(m.View()) > 16 {
		t.Fatal("repeated exit or resize reset the animation or exceeded the terminal")
	}
}

func TestExitActionsAnimateBeforeQuitting(t *testing.T) {
	for _, test := range []struct {
		name  string
		state AppState
		key   tea.KeyMsg
		tab   int
		help  bool
	}{
		{"boot", StateBoot, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}, 0, false},
		{"main", StateMain, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}, 0, false},
		{"help", StateMain, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}, 0, true},
		{"exit page", StateMain, tea.KeyMsg{Type: tea.KeyEnter}, 4, false},
		{"form", StateContactForm, tea.KeyMsg{Type: tea.KeyCtrlC}, 3, false},
		{"confirmation", StateContactSent, tea.KeyMsg{Type: tea.KeyCtrlC}, 3, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			m := NewModel()
			m.State = test.state
			m.MenuIndex = test.tab
			m.ShowHelp = test.help
			m = updateModel(t, m, tea.WindowSizeMsg{Width: 80, Height: 24})
			updated, cmd := m.Update(test.key)
			m = updated.(Model)
			if m.State != StateExiting || cmd == nil {
				t.Fatal("exit action did not animate")
			}
			if _, quit := cmd().(tea.QuitMsg); quit {
				t.Fatal("exit action quit before animating")
			}
			m = updateModel(t, m, bootDoneMsg{})
			m = updateModel(t, m, sendResultMsg{err: context.Canceled})
			if m.State != StateExiting {
				t.Fatal("late result interrupted the exit")
			}
			_, cmd = m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
			if cmd == nil {
				t.Fatal("second Ctrl+C did not quit")
			}
			if _, quit := cmd().(tea.QuitMsg); !quit {
				t.Fatal("second Ctrl+C did not skip animation")
			}
		})
	}
}

func TestDisconnectFramesFitAndFinish(t *testing.T) {
	for _, size := range []tea.WindowSizeMsg{{Width: 32, Height: 15}, {Width: 80, Height: 24}, {Width: 160, Height: 45}, {Width: 8, Height: 5}, {Width: 1, Height: 1}, {Width: 250, Height: 70}} {
		m := NewModel()
		m.State = StateMain
		m = updateModel(t, m, size)
		updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
		m = updated.(Model)
		if cmd == nil {
			t.Fatal("no animation tick")
		}
		tick := cmd()
		if _, quit := tick.(tea.QuitMsg); quit {
			t.Fatal("quit before animation")
		}
		sawSignal, quit := false, false
		signalFrames := 0
		for i := 0; i < 100; i++ {
			view := m.View()
			sawSignal = sawSignal || strings.Contains(view, "SIGNAL LOST")
			if strings.Contains(view, "SIGNAL LOST") {
				signalFrames++
			}
			if lipgloss.Height(view) > size.Height {
				t.Fatal("shutdown exceeds terminal height")
			}
			for _, line := range strings.Split(view, "\n") {
				if lipgloss.Width(line) > size.Width {
					t.Fatal("shutdown exceeds terminal width")
				}
			}
			updated, cmd = m.Update(tick)
			m = updated.(Model)
			if cmd == nil {
				t.Fatal("shutdown lost its animation tick")
			}
			// Only the final command quits; intermediate commands schedule the same tick type.
			if i >= exitSignalFrame+exitHoldFrames-1 {
				if _, quit = cmd().(tea.QuitMsg); quit {
					break
				}
			}
		}
		if !quit {
			t.Fatal("shutdown did not finish within its bounded frame count")
		}
		if size.Width >= 32 && !sawSignal {
			t.Fatal("shutdown did not show SIGNAL LOST")
		}
		if size.Width >= 32 && time.Duration(signalFrames)*exitFrameInterval < time.Second {
			t.Fatal("SIGNAL LOST was not held for a full second")
		}
	}
}

func TestDisconnectStartsWithEntireCurrentScreen(t *testing.T) {
	for _, size := range []tea.WindowSizeMsg{{Width: 80, Height: 24}, {Width: 160, Height: 45}} {
		m := NewModel()
		m.State = StateMain
		m = updateModel(t, m, size)
		before := m.View()
		m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
		if m.View() != before {
			t.Errorf("%dx%d: disconnect cropped or replaced the first frame", size.Width, size.Height)
		}
		m = updateModel(t, m, exitTickMsg{})
		lines := strings.Split(m.View(), "\n")
		if strings.TrimSpace(ansi.Strip(lines[0])) == "" || strings.TrimSpace(ansi.Strip(lines[len(lines)-1])) == "" {
			t.Error("effect does not reach the top and bottom of the screen")
		}
		m.ExitFrame = 8
		if size.Width == 160 && !strings.Contains(m.View(), "OPERATOR FILE") {
			t.Error("full-screen effect dropped the right-hand inspector")
		}
		for range 48 {
			m = updateModel(t, m, exitTickMsg{})
		}
		if m.State != StateExiting || strings.Contains(m.View(), "SIGNAL LOST") {
			t.Error("effect finished too soon to see the full collapse")
		}
	}
}

func BenchmarkDisconnectRendering(b *testing.B) {
	m := NewModel()
	m.State = StateMain
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 160, Height: 45})
	m = updated.(Model)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	m = updated.(Model)
	m.ExitFrame = 12
	b.ResetTimer()
	for b.Loop() {
		_ = m.View()
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

func TestProjectArrowsBrowseAtScrollBoundaries(t *testing.T) {
	for _, size := range []tea.WindowSizeMsg{{Width: 160, Height: 50}, {Width: 80, Height: 24}} {
		m := NewModel()
		m.State = StateMain
		m.MenuIndex = 2
		m = updateModel(t, m, size)
		m.Viewport.GotoBottom()
		m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyDown})
		if m.ProjectIndex != 1 || m.ActiveTab() != "projects" || !m.Viewport.AtTop() {
			t.Fatalf("down at bottom did not open next project at %dx%d", size.Width, size.Height)
		}
		if !strings.Contains(m.View(), "SELECTED: Lattora") {
			t.Fatal("new project not visible")
		}
		m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
		if m.ProjectIndex != 0 || !m.Viewport.AtBottom() {
			t.Fatal("k at top did not return to previous project")
		}
		m = updateModel(t, m, tea.MouseMsg{Action: tea.MouseActionPress, Button: tea.MouseButtonWheelDown})
		if m.ProjectIndex != 1 || !m.Viewport.AtTop() {
			t.Fatal("wheel at bottom did not open next project")
		}
		m = updateModel(t, m, tea.MouseMsg{Action: tea.MouseActionPress, Button: tea.MouseButtonWheelUp})
		if m.ProjectIndex != 0 || !m.Viewport.AtBottom() {
			t.Fatal("wheel at top did not return to previous project")
		}
	}
}

func TestPermanentProviderFailureOffersDirectEmail(t *testing.T) {
	t.Setenv("RESEND_API_KEY", "test-key")
	t.Setenv("RESEND_FROM", "Portfolio <from@example.com>")
	t.Setenv("RESEND_TO", "owner@example.com")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		io.WriteString(w, `{"name":"validation_error","message":"The example.com domain is not verified"}`)
	}))
	defer server.Close()
	old := resendEndpoint
	resendEndpoint = server.URL
	t.Cleanup(func() { resendEndpoint = old })
	resetContactLimit()
	err := SendContactForm(context.Background(), "Visitor", "visitor@example.com", "Hello")
	if err == nil || !strings.Contains(err.Error(), "owner@example.com") || strings.Contains(err.Error(), "try again") {
		t.Fatalf("permanent provider rejection = %v", err)
	}
}

func TestSSHCoalescedNavigationKeysAreNotDropped(t *testing.T) {
	m := NewModel()
	m.State = StateMain
	m = updateModel(t, m, tea.WindowSizeMsg{Width: 80, Height: 24})
	m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("llj")})
	if m.ActiveTab() != "projects" || m.Viewport.YOffset != 1 {
		t.Fatal("coalesced navigation keys were dropped")
	}
	m.State = StateContactForm
	m.focusContact(0)
	m = updateModel(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("hello")})
	if m.ContactInputs[0].Value() != "hello" {
		t.Fatal("contact text was treated as navigation")
	}
}
