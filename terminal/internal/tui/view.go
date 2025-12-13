package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
var encryptChars = []string{"█", "▓", "▒", "░", "▪", "▫", "◆", "◇", "○", "●", "◐", "◑", "◒", "◓"}

func (m Model) View() string {
	switch m.State {
	case StateBoot:
		return m.renderBoot()
	case StateLoading:
		return m.renderLoading()
	case StateAccessGranted:
		return m.renderAccessGranted()
	case StateContactForm:
		return m.renderContactForm()
	case StateContactSending:
		return m.renderContactSending()
	case StateContactSent:
		return m.renderContactSent()
	default:
		return m.renderMain()
	}
}

func (m Model) renderBoot() string {
	frame := spinnerFrames[m.BootFrame%len(spinnerFrames)]
	spinner := m.Styles.BootSpinner.Render(frame)
	text := m.Styles.Dim.Render(" Establishing secure link to PUNEET-MAINFRAME...")

	// Use a style without MarginBottom for ASCII art so lines are adjacent
	asciiStyle := m.Styles.Title.MarginBottom(0)
	asciiTitle := asciiStyle.Render("█▀█ █░█ █▄░█ █▀▀ █▀▀ ▀█▀ ▄▄ █▀█ █▀\n█▀▀ █▄█ █░▀█ ██▄ ██▄ ░█░ ░░ █▄█ ▄█")

	content := lipgloss.JoinVertical(lipgloss.Center,
		"",
		"",
		asciiTitle,
		"",
		spinner+text,
		"",
	)

	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) renderLoading() string {
	barWidth := 40
	filled := (m.BootProgress * barWidth) / 100
	empty := barWidth - filled

	bar := m.Styles.BootProgress.Render("[" +
		strings.Repeat("█", filled) +
		strings.Repeat("░", empty) + "]")

	percent := fmt.Sprintf(" %d%%", m.BootProgress)

	var status string
	switch {
	case m.BootProgress < 30:
		status = "Loading bio data..."
	case m.BootProgress < 60:
		status = "Compiling projects..."
	case m.BootProgress < 90:
		status = "Initializing UI subsystems..."
	default:
		status = "Finalizing connection..."
	}

	content := lipgloss.JoinVertical(lipgloss.Center,
		"",
		m.Styles.Title.Render("SYSTEM INITIALIZATION"),
		"",
		bar+percent,
		"",
		m.Styles.Dim.Render(status),
		"",
	)

	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) renderAccessGranted() string {
	content := lipgloss.JoinVertical(lipgloss.Center,
		"",
		"",
		m.Styles.BootAccessGranted.Render("▓▓▓ ACCESS GRANTED ▓▓▓"),
		"",
		m.Styles.Success.Render("Welcome, Operator"),
		"",
		"",
	)

	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) renderContactForm() string {
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorDim).
		Padding(0, 1).
		Width(50)

	focusedInputStyle := inputStyle.BorderForeground(ColorPrimary)

	labels := []string{"CALLSIGN", "COMM FREQ (Email)", "MESSAGE"}

	var inputs []string
	for i, input := range m.ContactInputs {
		style := inputStyle
		if i == m.ContactFocus {
			style = focusedInputStyle
		}
		label := m.Styles.Dim.Render(labels[i])
		inputs = append(inputs, label+"\n"+style.Render(input.View()))
	}

	formContent := lipgloss.JoinVertical(lipgloss.Left, inputs...)

	content := lipgloss.JoinVertical(lipgloss.Center,
		"",
		m.Styles.Title.Render("╔═══ SECURE UPLINK ═══╗"),
		"",
		m.Styles.Subtitle.Render("Establish direct communication channel"),
		"",
		formContent,
		"",
		m.Styles.Dim.Render("[Tab] Next  [Enter] Transmit  [Esc] Abort"),
		"",
	)

	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) renderContactSending() string {
	// Generate "encryption" animation
	barWidth := 30
	progress := (m.SendingFrame * barWidth) / 30

	var encryptedBar string
	for i := 0; i < barWidth; i++ {
		if i < progress {
			encryptedBar += m.Styles.Success.Render("█")
		} else {
			char := encryptChars[(m.SendingFrame+i)%len(encryptChars)]
			encryptedBar += m.Styles.BootSpinner.Render(char)
		}
	}

	var statusText string
	switch {
	case m.SendingFrame < 8:
		statusText = "Encrypting message..."
	case m.SendingFrame < 16:
		statusText = "Establishing secure tunnel..."
	case m.SendingFrame < 24:
		statusText = "Transmitting via Neural Link..."
	default:
		statusText = "Finalizing uplink..."
	}

	content := lipgloss.JoinVertical(lipgloss.Center,
		"",
		"",
		m.Styles.Title.Render("▓▓▓ TRANSMISSION IN PROGRESS ▓▓▓"),
		"",
		"["+encryptedBar+"]",
		"",
		m.Styles.Dim.Render(statusText),
		"",
		"",
	)

	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) renderContactSent() string {
	content := lipgloss.JoinVertical(lipgloss.Center,
		"",
		"",
		m.Styles.Success.Render("████████████████████████████████"),
		"",
		m.Styles.BootAccessGranted.Render("  ✓ TRANSMISSION COMPLETE  "),
		"",
		m.Styles.Success.Render("████████████████████████████████"),
		"",
		m.Styles.Dim.Render("Message encrypted and delivered"),
		m.Styles.Dim.Render("Operator will respond shortly"),
		"",
		m.Styles.Subtitle.Render("Press any key to continue..."),
		"",
	)

	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) renderMain() string {
	header := m.GetHeader()
	footer := m.GetFooter()

	contentHeight := m.Height - 4

	// Calculate pane widths - responsive design
	menuWidth := 22
	if m.Width < 60 {
		menuWidth = 16
	}

	inspectorWidth := 0
	if m.Width > 100 {
		inspectorWidth = 26
	}
	viewportWidth := m.Width - menuWidth - inspectorWidth - 6

	if viewportWidth < 20 {
		viewportWidth = 20
	}

	menu := m.renderMenu(menuWidth, contentHeight)
	viewport := m.renderViewport(viewportWidth, contentHeight)

	var panes string
	if inspectorWidth > 0 {
		inspector := m.renderInspector(inspectorWidth, contentHeight)
		panes = lipgloss.JoinHorizontal(lipgloss.Top, menu, viewport, inspector)
	} else {
		panes = lipgloss.JoinHorizontal(lipgloss.Top, menu, viewport)
	}

	// Show help overlay if enabled
	if m.ShowHelp {
		return m.renderHelpOverlay()
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, panes, footer)
}

func (m Model) renderHelpOverlay() string {
	helpContent := lipgloss.JoinVertical(lipgloss.Left,
		m.Styles.Title.Render("╔═══ KEYBOARD SHORTCUTS ═══╗"),
		"",
		m.Styles.MenuActive.Render("Navigation"),
		m.Styles.Dim.Render("  ↑/k     Move up"),
		m.Styles.Dim.Render("  ↓/j     Move down"),
		m.Styles.Dim.Render("  ←/h     Previous project"),
		m.Styles.Dim.Render("  →/l     Next project"),
		m.Styles.Dim.Render("  Enter   Select/Confirm"),
		"",
		m.Styles.MenuActive.Render("General"),
		m.Styles.Dim.Render("  ?       Toggle this help"),
		m.Styles.Dim.Render("  q       Disconnect"),
		m.Styles.Dim.Render("  Ctrl+C  Force quit"),
		"",
		m.Styles.MenuActive.Render("Contact Form"),
		m.Styles.Dim.Render("  Tab     Next field"),
		m.Styles.Dim.Render("  Esc     Cancel"),
		m.Styles.Dim.Render("  Enter   Send message"),
		"",
		m.Styles.Subtitle.Render("Press ? to close"),
	)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(ColorPrimary).
		Padding(1, 3).
		Background(ColorBg)

	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, boxStyle.Render(helpContent))
}

func (m Model) renderMenu(width, height int) string {
	items := []string{"About Me", "Experience", "Projects", "Contact", "Exit"}

	var menuItems []string
	for i, item := range items {
		cursor := "  "
		style := m.Styles.MenuInactive
		if i == m.MenuIndex {
			cursor = "▸ "
			style = m.Styles.MenuActive
		}
		menuItems = append(menuItems, cursor+style.Render(item))
	}

	menu := lipgloss.JoinVertical(lipgloss.Left, menuItems...)

	// Add system stats at bottom
	stats := m.Styles.Dim.Render("\n\n━━━━━━━━━━━━━━━━\nSYS STATS:\nCPU: 12%\nMEM: 64MB\nUPTIME: 24d")
	menu = lipgloss.JoinVertical(lipgloss.Left, menu, stats)

	return m.Styles.Menu.Width(width).Height(height).Render(menu)
}

func (m Model) renderViewport(width, height int) string {
	var content string

	switch m.ActiveTab() {
	case "about":
		content = m.renderColorizedBio()
	case "experience":
		content = m.renderExperience()
	case "projects":
		content = m.renderProjectList()
	case "contact":
		content = m.renderContactInfo()
	case "exit":
		content = "\n\nPress Enter to disconnect from PUNEET-OS\n\n" +
			m.Styles.Dim.Render("Thank you for visiting.")
	}

	// Wrap content to fit viewport width
	contentWidth := width - 6
	if contentWidth < 20 {
		contentWidth = 20
	}
	wrapped := lipgloss.NewStyle().Width(contentWidth).Render(content)

	// Add navigation hint
	navHint := m.Styles.Dim.Render("\n[ ↑↓ Navigate Menu ]")

	return m.Styles.Viewport.Width(width).Height(height).Render(wrapped + navHint)
}

// renderExperience renders the experience/work history section
func (m Model) renderExperience() string {
	var lines []string

	lines = append(lines, m.Styles.Title.Render("MISSION LOG"))
	lines = append(lines, "")

	// Experience entries
	lines = append(lines, m.Styles.MenuActive.Render("Research Intern")+" @ CeAT, VIT")
	lines = append(lines, m.Styles.Dim.Render("  May-Jul 2025"))
	lines = append(lines, "  CloudSim Plus framework for")
	lines = append(lines, "  VM placement optimization")
	lines = append(lines, "")

	lines = append(lines, m.Styles.MenuActive.Render("Full Stack Developer")+" @ Daira")
	lines = append(lines, m.Styles.Dim.Render("  Dec 2024 - Feb 2025"))
	lines = append(lines, "  Backend APIs, "+m.Styles.Success.Render("30%")+" faster")
	lines = append(lines, "  data retrieval")
	lines = append(lines, "")

	lines = append(lines, m.Styles.MenuActive.Render("Web Dev Intern")+" @ IIT Bombay")
	lines = append(lines, m.Styles.Dim.Render("  Sept-Oct 2024"))
	lines = append(lines, "  Reduced payload by "+m.Styles.Success.Render("90%"))
	lines = append(lines, "  Optimized DB queries")

	return strings.Join(lines, "\n")
}

// renderColorizedBio returns the bio with Tron-style coloring applied
func (m Model) renderColorizedBio() string {
	var lines []string

	// Operator header
	lines = append(lines, m.Styles.Title.Render("╔══════════════════════════════╗"))
	lines = append(lines, m.Styles.Title.Render("║  OPERATOR: PUNEET CHANDNA    ║"))
	lines = append(lines, m.Styles.Title.Render("╚══════════════════════════════╝"))
	lines = append(lines, "")

	// Status line with colors
	lines = append(lines, "> STATUS: "+m.Styles.Success.Render("ONLINE"))
	lines = append(lines, "> ROLE: "+m.Styles.MenuActive.Render("SOFTWARE ENGINEER")+" | "+m.Styles.MenuActive.Render("FULL STACK DEV"))
	lines = append(lines, "> LOCATION: "+m.Styles.Dim.Render("VIT CHENNAI, INDIA"))
	lines = append(lines, "> GRADUATION: "+m.Styles.AccentText.Render("2026"))
	lines = append(lines, "")

	// Tech Arsenal - responsive grid
	lines = append(lines, m.Styles.Title.Render("TECH ARSENAL:"))
	lines = append(lines, "  "+m.Styles.MenuActive.Render("Python")+"  "+m.Styles.MenuActive.Render("JavaScript")+"  "+m.Styles.MenuActive.Render("Go"))
	lines = append(lines, "  "+m.Styles.Dim.Render("C++")+"     "+m.Styles.Dim.Render("Node.js")+"     "+m.Styles.Dim.Render("React"))
	lines = append(lines, "  "+m.Styles.Dim.Render("Next.js")+" "+m.Styles.Dim.Render("MongoDB")+"     "+m.Styles.Dim.Render("PostgreSQL"))
	lines = append(lines, "  "+m.Styles.AccentText.Render("AWS/GCP")+" "+m.Styles.AccentText.Render("Docker")+"      "+m.Styles.AccentText.Render("Linux/Git"))
	lines = append(lines, "")

	// Core Competencies
	lines = append(lines, m.Styles.Title.Render("CORE COMPETENCIES:"))
	lines = append(lines, "  • "+m.Styles.Success.Render("Full Stack Dev")+" & APIs")
	lines = append(lines, "  • Data Structures & Algo")
	lines = append(lines, "  • System Design")
	lines = append(lines, "  • "+m.Styles.AccentText.Render("Cryptography")+" (AES)")
	lines = append(lines, "  • Cloud Computing")
	lines = append(lines, "")

	lines = append(lines, m.Styles.Subtitle.Render("> \"Building performant &"))
	lines = append(lines, m.Styles.Subtitle.Render("   elegant software.\""))

	return strings.Join(lines, "\n")
}

func (m Model) renderProjectList() string {
	var lines []string
	lines = append(lines, m.Styles.Title.Render("╔═══ PROJECT DATABASE ═══╗"))
	lines = append(lines, "")

	for i, p := range m.Projects {
		status := p.Status
		statusStyle := m.Styles.Dim
		if status == "Live" {
			statusStyle = m.Styles.Success
		} else if status == "WIP" {
			statusStyle = m.Styles.AccentText
		}

		// Highlight selected project
		numStyle := m.Styles.Dim
		nameStyle := m.Styles.MenuInactive
		if i == m.ProjectIndex {
			numStyle = m.Styles.MenuActive
			nameStyle = m.Styles.MenuActive
		}

		lines = append(lines, numStyle.Render(fmt.Sprintf("%d. ", i+1))+nameStyle.Render(p.Name))
		lines = append(lines, "   "+m.Styles.Dim.Render(p.Description))
		lines = append(lines, "   Status: "+statusStyle.Render(status))
		lines = append(lines, "")
	}

	lines = append(lines, m.Styles.Dim.Render("[↑↓] Select project  [←→] Navigate"))

	return strings.Join(lines, "\n")
}

func (m Model) renderContactInfo() string {
	return m.Styles.Title.Render("╔═══ UPLINK TERMINAL ═══╗") + `

Establish communication via secure channels:

` + m.Styles.MenuActive.Render("EMAIL") + `
  puneetchandna@zohomail.in

` + m.Styles.MenuActive.Render("GITHUB") + `
  github.com/puneet-chandna

` + m.Styles.MenuActive.Render("LINKEDIN") + `
  linkedin.com/in/puneet-chandna

` + m.Styles.MenuActive.Render("WEBSITE") + `
  puneetchandna.com

` + m.Styles.Dim.Render(`
━━━━━━━━━━━━━━━━━━━━━━━━━━━
> Encrypted channel ready`) + `

` + m.Styles.AccentText.Render("Press [Enter] to open secure uplink form")
}

func (m Model) renderInspector(width, height int) string {
	var content string

	switch m.ActiveTab() {
	case "about":
		content = m.Styles.Title.Render("OPERATOR FILE") + "\n\n" +
			"Name: Puneet\n" +
			"Role: Engineer\n" +
			"Status: " + m.Styles.Success.Render("Active")

	case "projects":
		if len(m.Projects) > 0 && m.ProjectIndex < len(m.Projects) {
			p := m.Projects[m.ProjectIndex]
			content = m.Styles.Title.Render("DETAILS") + "\n\n" +
				"Stack:\n"
			for _, s := range p.Stack {
				content += "  • " + s + "\n"
			}
			content += "\nYear: " + fmt.Sprint(p.Year)
			if p.URL != "" {
				content += "\n\n" + m.Styles.Dim.Render("URL:\n"+p.URL)
			}
		}

	case "contact":
		content = m.Styles.Title.Render("UPLINK STATUS") + "\n\n" +
			m.Styles.Success.Render("● Channel Open") + "\n\n" +
			m.Styles.Dim.Render("Latency: 42ms\nEncryption: AES-256")

	case "exit":
		content = m.Styles.AccentText.Render("DISCONNECT?") + "\n\n" +
			m.Styles.Dim.Render("Press Enter\nto confirm")
	}

	return m.Styles.Inspector.Width(width).Height(height).Render(content)
}
