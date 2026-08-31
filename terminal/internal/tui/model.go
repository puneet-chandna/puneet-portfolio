package tui

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/puneet/terminal-portfolio/internal/data"
)

type AppState int

const (
	StateBoot AppState = iota
	StateLoading
	StateAccessGranted
	StateMain
	StateContactForm
	StateContactSending
	StateContactSent
)

type Model struct {
	Width        int
	Height       int
	State        AppState
	BootProgress int
	BootFrame    int
	MenuIndex    int
	Projects     []data.Project
	Bio          string
	Experience   string
	Context      context.Context

	// Contact form
	ContactInputs []textinput.Model
	ContactFocus  int
	ContactError  string
	SendingFrame  int
	ShowHelp      bool
	ProjectIndex  int // For scrolling through projects

	// Scrollable viewport
	Viewport    viewport.Model
	viewportTab string

	// Per-session styles (for proper SSH color rendering)
	Styles Styles
}

type tickMsg time.Time
type sendingTickMsg time.Time
type bootDoneMsg struct{}
type sendResultMsg struct{ err error }

const contactInputFrameWidth = 4

// NewModel creates a model with default styles (for local testing)
func NewModel() Model {
	return NewModelWithRenderer(lipgloss.DefaultRenderer())
}

// NewModelWithRenderer creates a model with styles from the given renderer
// For SSH sessions, use bubbletea.MakeRenderer(s) for proper color support
func NewModelWithRenderer(r *lipgloss.Renderer) Model {
	// Create text inputs for contact form
	nameInput := textinput.New()
	nameInput.Placeholder = "Your callsign..."
	nameInput.CharLimit = 50
	nameInput.Width = 40

	emailInput := textinput.New()
	emailInput.Placeholder = "your@email.com"
	emailInput.CharLimit = 100
	emailInput.Width = 40

	msgInput := textinput.New()
	msgInput.Placeholder = "Your message..."
	msgInput.CharLimit = 500
	msgInput.Width = 40

	// Initialize viewport with default dimensions
	vp := viewport.New(60, 15)
	vp.Style = lipgloss.NewStyle()

	return Model{
		Width:         80,
		Height:        24,
		State:         StateBoot,
		BootProgress:  0,
		BootFrame:     0,
		MenuIndex:     0,
		Projects:      data.GetProjects(),
		Bio:           data.GetBio(),
		Experience:    data.GetExperience(),
		Context:       context.Background(),
		ContactInputs: []textinput.Model{nameInput, emailInput, msgInput},
		ContactFocus:  0,
		ShowHelp:      false,
		ProjectIndex:  0,
		Viewport:      vp,
		Styles:        NewStyles(r),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tickCmd(), tea.EnterAltScreen)
}

func (m Model) contactFormWidth() int {
	return min(54, max(1, m.Width-2))
}

func (m Model) contactInputWidth(input textinput.Model) int {
	return max(1, m.contactFormWidth()-contactInputFrameWidth-lipgloss.Width(input.Prompt)-1)
}

func (m *Model) resizeContactInputs() {
	for i := range m.ContactInputs {
		m.ContactInputs[i].Width = m.contactInputWidth(m.ContactInputs[i])
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*80, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) ActiveTab() string {
	tabs := []string{"about", "experience", "projects", "contact", "exit"}
	if m.MenuIndex >= 0 && m.MenuIndex < len(tabs) {
		return tabs[m.MenuIndex]
	}
	return "about"
}

func (m Model) GetHeader() string {
	tab := "/" + m.ActiveTab()
	title := m.Styles.Title.MarginBottom(0).Render("PUNEET-OS v1.0")

	left := title
	center := m.Styles.Dim.Render(tab)
	right := m.Styles.Success.Render("● ONLINE")

	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(right)
	contentWidth := max(0, m.Width-m.Styles.Header.GetHorizontalPadding())
	centerWidth := contentWidth - leftWidth - rightWidth

	if centerWidth < 0 {
		centerWidth = 0
	}

	header := lipgloss.JoinHorizontal(
		lipgloss.Center,
		left,
		lipgloss.NewStyle().Width(centerWidth).Align(lipgloss.Center).Render(center),
		right,
	)

	return m.Styles.Header.Width(max(0, m.Width-m.Styles.Header.GetHorizontalBorderSize())).Render(header)
}

func (m Model) GetFooter() string {
	var controls string
	switch m.State {
	case StateContactForm:
		controls = m.Styles.Dim.Render("[Tab] Next field  [Enter] Send  [Esc] Back")
	case StateMain:
		if m.ShowHelp {
			controls = m.Styles.Dim.Render("[?] Help  [↑↓/jk] Select  [PgUp/PgDn] Scroll  [q] Exit")
		} else {
			controls = m.Styles.Dim.Render("[?] Help  [↑↓/jk] Select  [PgUp/PgDn] Scroll  [q] Exit")
		}
	default:
		controls = m.Styles.Dim.Render("[↑↓/jk] Navigate  [Enter] Select  [q] Exit")
	}

	leftWidth := lipgloss.Width(controls)
	contentWidth := max(0, m.Width-m.Styles.Footer.GetHorizontalPadding())
	pos := m.Styles.Dim.Render(fmt.Sprintf("Ln %d, Col 1", m.Viewport.YOffset+1))
	gap := contentWidth - leftWidth - lipgloss.Width(pos)
	if gap < 2 {
		pos = ""
		gap = contentWidth - leftWidth
	}

	if gap < 0 {
		gap = 0
	}

	footer := lipgloss.JoinHorizontal(
		lipgloss.Center,
		controls,
		lipgloss.NewStyle().Width(gap).Render(""),
		pos,
	)

	return m.Styles.Footer.Width(max(0, m.Width-m.Styles.Footer.GetHorizontalBorderSize())).Render(footer)
}
