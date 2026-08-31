package tui

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)

	case tea.WindowSizeMsg:
		m.Width = clamp(msg.Width, 1, 500)
		m.Height = clamp(msg.Height, 1, 200)
		m.resizeContactInputs()
		m.refreshViewport(true)
		return m, nil

	case sendingTickMsg:
		if m.State == StateContactSending {
			m.SendingFrame = (m.SendingFrame + 1) % len(encryptChars)
			return m, sendingTickCmd()
		}
		return m, nil

	case tickMsg:
		return m.handleTick()

	case bootDoneMsg:
		m.State = StateMain
		return m, nil

	case sendResultMsg:
		if msg.err != nil {
			m.State = StateContactForm
			m.ContactError = msg.err.Error()
			m.focusContact(0)
		} else {
			m.ContactError = ""
			m.State = StateContactSent
		}
		return m, nil
	}

	// Update text inputs when in contact form
	if m.State == StateContactForm {
		return m.updateContactInputs(msg)
	}

	// Also update viewport messages for scrolling in main state
	if m.State == StateMain {
		var cmd tea.Cmd
		m.Viewport, cmd = m.Viewport.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) updateContactInputs(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.ContactInputs[m.ContactFocus], cmd = m.ContactInputs[m.ContactFocus].Update(msg)
	return m, cmd
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// During boot sequence, any key skips to main
	if m.State == StateBoot || m.State == StateLoading || m.State == StateAccessGranted {
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		m.State = StateMain
		return m, nil
	}

	// Allow an immediate exit while delivery observes the session context.
	if m.State == StateContactSending {
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		return m, nil
	}

	// Contact sent - any key returns to main
	if m.State == StateContactSent {
		m.State = StateMain
		// Reset form
		for i := range m.ContactInputs {
			m.ContactInputs[i].SetValue("")
		}
		m.ContactFocus = 0
		return m, nil
	}

	// Contact form mode
	if m.State == StateContactForm {
		return m.handleContactFormKey(msg)
	}

	// Main navigation
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "?":
		m.ShowHelp = !m.ShowHelp

	case "up", "k":
		if m.ActiveTab() == "projects" && m.ProjectIndex > 0 {
			m.ProjectIndex--
			m.refreshViewport(true)
		} else if m.MenuIndex > 0 {
			m.MenuIndex--
			m.refreshViewport(true)
		}

	case "down", "j":
		if m.ActiveTab() == "projects" && m.ProjectIndex < len(m.Projects)-1 {
			m.ProjectIndex++
			m.refreshViewport(true)
		} else if m.MenuIndex < 4 {
			m.MenuIndex++
			m.refreshViewport(true)
		}

	case "enter":
		if m.ActiveTab() == "exit" {
			return m, tea.Quit
		}
		if m.ActiveTab() == "contact" {
			if !ContactFormConfigured() {
				m.ContactError = "Message delivery is not configured. Use the contact links above."
				m.refreshViewport(false)
				return m, nil
			}
			m.State = StateContactForm
			m.ContactError = ""
			m.focusContact(0)
			return m, nil
		}

	case "left", "h":
		if m.MenuIndex > 0 {
			m.MenuIndex--
			m.refreshViewport(true)
		}

	case "right", "l":
		if m.MenuIndex < 4 {
			m.MenuIndex++
			m.refreshViewport(true)
		}

	case "pgup", "pgdown", "home", "end", "ctrl+u", "ctrl+d":
		var cmd tea.Cmd
		m.Viewport, cmd = m.Viewport.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) handleContactFormKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit

	case "esc":
		m.State = StateMain
		m.ContactInputs[m.ContactFocus].Blur()
		return m, nil

	case "tab", "down":
		m.ContactInputs[m.ContactFocus].Blur()
		m.ContactFocus = (m.ContactFocus + 1) % len(m.ContactInputs)
		m.ContactInputs[m.ContactFocus].Focus()
		return m, nil

	case "shift+tab", "up":
		m.ContactInputs[m.ContactFocus].Blur()
		m.ContactFocus--
		if m.ContactFocus < 0 {
			m.ContactFocus = len(m.ContactInputs) - 1
		}
		m.ContactInputs[m.ContactFocus].Focus()
		return m, nil

	case "enter":
		if _, _, _, err := validateContact(m.ContactInputs[0].Value(), m.ContactInputs[1].Value(), m.ContactInputs[2].Value()); err != nil {
			m.ContactError = err.Error()
			return m, nil
		}
		if !ContactFormConfigured() {
			m.ContactError = "Message delivery is not configured. Use the contact links above."
			return m, nil
		}
		m.State = StateContactSending
		m.SendingFrame = 0
		return m, tea.Batch(m.sendContactCmd(), sendingTickCmd())
	}

	// Pass other keys to the focused input
	var cmd tea.Cmd
	m.ContactInputs[m.ContactFocus], cmd = m.ContactInputs[m.ContactFocus].Update(msg)
	m.ContactError = ""
	return m, cmd
}

func (m *Model) focusContact(index int) {
	for i := range m.ContactInputs {
		m.ContactInputs[i].Blur()
	}
	m.ContactFocus = index
	m.ContactInputs[index].Focus()
}

func (m Model) sendContactCmd() tea.Cmd {
	ctx := m.Context
	if ctx == nil {
		ctx = context.Background()
	}
	name, email, message := m.ContactInputs[0].Value(), m.ContactInputs[1].Value(), m.ContactInputs[2].Value()
	return func() tea.Msg { return sendResultMsg{err: SendContactForm(ctx, name, email, message)} }
}

func (m Model) handleTick() (tea.Model, tea.Cmd) {
	switch m.State {
	case StateBoot:
		m.BootFrame++
		if m.BootFrame > 15 {
			m.State = StateLoading
			m.BootFrame = 0
		}
		return m, tickCmd()

	case StateLoading:
		m.BootProgress += 5
		if m.BootProgress >= 100 {
			m.State = StateAccessGranted
			m.BootFrame = 0
		}
		return m, tickCmd()

	case StateAccessGranted:
		m.BootFrame++
		if m.BootFrame > 12 {
			return m, func() tea.Msg {
				time.Sleep(200 * time.Millisecond)
				return bootDoneMsg{}
			}
		}
		return m, tickCmd()

	case StateMain:
		return m, tickCmd()
	}

	return m, nil
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func sendingTickCmd() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg { return sendingTickMsg(t) })
}
