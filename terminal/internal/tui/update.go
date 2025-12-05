package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tickMsg:
		return m.handleTick()

	case bootDoneMsg:
		m.State = StateMain
		return m, nil

	case sendDoneMsg:
		m.State = StateContactSent
		return m, nil
	}

	// Update text inputs when in contact form
	if m.State == StateContactForm {
		return m.updateContactInputs(msg)
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

	// Contact sending animation - ignore keys
	if m.State == StateContactSending {
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
		} else if m.MenuIndex > 0 {
			m.MenuIndex--
		}

	case "down", "j":
		if m.ActiveTab() == "projects" && m.ProjectIndex < len(m.Projects)-1 {
			m.ProjectIndex++
		} else if m.MenuIndex < 3 {
			m.MenuIndex++
		}

	case "enter":
		if m.ActiveTab() == "exit" {
			return m, tea.Quit
		}
		if m.ActiveTab() == "contact" {
			m.State = StateContactForm
			m.ContactInputs[0].Focus()
			return m, nil
		}

	case "left", "h":
		if m.ProjectIndex > 0 {
			m.ProjectIndex--
		}

	case "right", "l":
		if m.ProjectIndex < len(m.Projects)-1 {
			m.ProjectIndex++
		}
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
		// Submit if we have content
		if m.ContactInputs[0].Value() != "" && m.ContactInputs[1].Value() != "" {
			m.State = StateContactSending
			m.SendingFrame = 0
			// Store values for sending
			name := m.ContactInputs[0].Value()
			email := m.ContactInputs[1].Value()
			message := m.ContactInputs[2].Value()
			// Send form in background
			go SendContactForm(name, email, message)
			return m, tickCmd()
		}
	}

	// Pass other keys to the focused input
	var cmd tea.Cmd
	m.ContactInputs[m.ContactFocus], cmd = m.ContactInputs[m.ContactFocus].Update(msg)
	return m, cmd
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

	case StateContactSending:
		m.SendingFrame++
		if m.SendingFrame > 30 {
			return m, func() tea.Msg {
				time.Sleep(500 * time.Millisecond)
				return sendDoneMsg{}
			}
		}
		return m, tickCmd()
	}

	return m, nil
}
