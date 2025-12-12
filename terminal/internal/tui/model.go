package tui

import (
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

	// Contact form
	ContactInputs []textinput.Model
	ContactFocus  int
	SendingFrame  int
	ShowHelp      bool
	ProjectIndex  int // For scrolling through projects

	// Scrollable viewport
	Viewport      viewport.Model
	ViewportReady bool

	// Animation frame for nav hints
	NavHintFrame int
}

type tickMsg time.Time
type bootDoneMsg struct{}
type sendDoneMsg struct{}

func NewModel() Model {
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
		ContactInputs: []textinput.Model{nameInput, emailInput, msgInput},
		ContactFocus:  0,
		ShowHelp:      false,
		ProjectIndex:  0,
		Viewport:      vp,
		ViewportReady: false,
		NavHintFrame:  0,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tickCmd(), tea.EnterAltScreen)
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
	status := SuccessStyle.Render("● ONLINE")
	tab := "/" + m.ActiveTab()
	title := TitleStyle.Render("PUNEET-OS v1.0")

	left := title
	center := DimStyle.Render(tab)
	right := status

	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(right)
	centerWidth := m.Width - leftWidth - rightWidth - 4

	if centerWidth < 0 {
		centerWidth = 0
	}

	header := lipgloss.JoinHorizontal(
		lipgloss.Center,
		left,
		lipgloss.NewStyle().Width(centerWidth).Align(lipgloss.Center).Render(center),
		right,
	)

	return HeaderStyle.Width(m.Width).Render(header)
}

func (m Model) GetFooter() string {
	var controls string
	switch m.State {
	case StateContactForm:
		controls = DimStyle.Render("[Tab] Next field  [Enter] Send  [Esc] Back")
	case StateMain:
		if m.ShowHelp {
			controls = DimStyle.Render("[?] Hide help  [↑↓/jk] Navigate  [Enter] Select  [q] Exit")
		} else {
			controls = DimStyle.Render("[?] Help  [↑↓/jk] Navigate  [Enter] Select  [q] Exit")
		}
	default:
		controls = DimStyle.Render("[↑↓/jk] Navigate  [Enter] Select  [q] Exit")
	}

	pos := DimStyle.Render("Ln 1, Col 1")

	leftWidth := lipgloss.Width(controls)
	rightWidth := lipgloss.Width(pos)
	gap := m.Width - leftWidth - rightWidth - 4

	if gap < 0 {
		gap = 0
	}

	footer := lipgloss.JoinHorizontal(
		lipgloss.Center,
		controls,
		lipgloss.NewStyle().Width(gap).Render(""),
		pos,
	)

	return FooterStyle.Width(m.Width).Render(footer)
}
