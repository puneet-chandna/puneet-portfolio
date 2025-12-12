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
	spinner := BootSpinnerStyle.Render(frame)
	text := DimStyle.Render(" Establishing secure link to PUNEET-MAINFRAME...")

	// Use a style without MarginBottom for ASCII art so lines are adjacent
	asciiStyle := TitleStyle.MarginBottom(0)
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

	bar := BootProgressStyle.Render("[" +
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
		TitleStyle.Render("SYSTEM INITIALIZATION"),
		"",
		bar+percent,
		"",
		DimStyle.Render(status),
		"",
	)

	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) renderAccessGranted() string {
	content := lipgloss.JoinVertical(lipgloss.Center,
		"",
		"",
		BootAccessGrantedStyle.Render("▓▓▓ ACCESS GRANTED ▓▓▓"),
		"",
		SuccessStyle.Render("Welcome, Operator"),
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
		label := DimStyle.Render(labels[i])
		inputs = append(inputs, label+"\n"+style.Render(input.View()))
	}

	formContent := lipgloss.JoinVertical(lipgloss.Left, inputs...)

	content := lipgloss.JoinVertical(lipgloss.Center,
		"",
		TitleStyle.Render("╔═══ SECURE UPLINK ═══╗"),
		"",
		SubtitleStyle.Render("Establish direct communication channel"),
		"",
		formContent,
		"",
		DimStyle.Render("[Tab] Next  [Enter] Transmit  [Esc] Abort"),
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
			encryptedBar += SuccessStyle.Render("█")
		} else {
			char := encryptChars[(m.SendingFrame+i)%len(encryptChars)]
			encryptedBar += BootSpinnerStyle.Render(char)
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
		TitleStyle.Render("▓▓▓ TRANSMISSION IN PROGRESS ▓▓▓"),
		"",
		"["+encryptedBar+"]",
		"",
		DimStyle.Render(statusText),
		"",
		"",
	)

	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) renderContactSent() string {
	content := lipgloss.JoinVertical(lipgloss.Center,
		"",
		"",
		SuccessStyle.Render("████████████████████████████████"),
		"",
		BootAccessGrantedStyle.Render("  ✓ TRANSMISSION COMPLETE  "),
		"",
		SuccessStyle.Render("████████████████████████████████"),
		"",
		DimStyle.Render("Message encrypted and delivered"),
		DimStyle.Render("Operator will respond shortly"),
		"",
		SubtitleStyle.Render("Press any key to continue..."),
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
		TitleStyle.Render("╔═══ KEYBOARD SHORTCUTS ═══╗"),
		"",
		MenuActiveStyle.Render("Navigation"),
		DimStyle.Render("  ↑/k     Move up"),
		DimStyle.Render("  ↓/j     Move down"),
		DimStyle.Render("  ←/h     Previous project"),
		DimStyle.Render("  →/l     Next project"),
		DimStyle.Render("  Enter   Select/Confirm"),
		"",
		MenuActiveStyle.Render("General"),
		DimStyle.Render("  ?       Toggle this help"),
		DimStyle.Render("  q       Disconnect"),
		DimStyle.Render("  Ctrl+C  Force quit"),
		"",
		MenuActiveStyle.Render("Contact Form"),
		DimStyle.Render("  Tab     Next field"),
		DimStyle.Render("  Esc     Cancel"),
		DimStyle.Render("  Enter   Send message"),
		"",
		SubtitleStyle.Render("Press ? to close"),
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
		style := MenuInactiveStyle
		if i == m.MenuIndex {
			cursor = "▸ "
			style = MenuActiveStyle
		}
		menuItems = append(menuItems, cursor+style.Render(item))
	}

	menu := lipgloss.JoinVertical(lipgloss.Left, menuItems...)

	// Add system stats at bottom
	stats := DimStyle.Render("\n\n━━━━━━━━━━━━━━━━\nSYS STATS:\nCPU: 12%\nMEM: 64MB\nUPTIME: 24d")
	menu = lipgloss.JoinVertical(lipgloss.Left, menu, stats)

	return MenuStyle.Width(width).Height(height).Render(menu)
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
			DimStyle.Render("Thank you for visiting.")
	}

	// Wrap content to fit viewport width
	contentWidth := width - 6
	if contentWidth < 20 {
		contentWidth = 20
	}
	wrapped := lipgloss.NewStyle().Width(contentWidth).Render(content)

	// Add navigation hint
	navHint := DimStyle.Render("\n[ ↑↓ Navigate Menu ]")

	return ViewportStyle.Width(width).Height(height).Render(wrapped + navHint)
}

// renderExperience renders the experience/work history section
func (m Model) renderExperience() string {
	var lines []string

	lines = append(lines, TitleStyle.Render("MISSION LOG"))
	lines = append(lines, "")

	// Experience entries
	lines = append(lines, MenuActiveStyle.Render("Research Intern")+" @ CeAT, VIT")
	lines = append(lines, DimStyle.Render("  May-Jul 2025"))
	lines = append(lines, "  CloudSim Plus framework for")
	lines = append(lines, "  VM placement optimization")
	lines = append(lines, "")

	lines = append(lines, MenuActiveStyle.Render("Full Stack Developer")+" @ Daira")
	lines = append(lines, DimStyle.Render("  Dec 2024 - Feb 2025"))
	lines = append(lines, "  Backend APIs, "+SuccessStyle.Render("30%")+" faster")
	lines = append(lines, "  data retrieval")
	lines = append(lines, "")

	lines = append(lines, MenuActiveStyle.Render("Web Dev Intern")+" @ IIT Bombay")
	lines = append(lines, DimStyle.Render("  Sept-Oct 2024"))
	lines = append(lines, "  Reduced payload by "+SuccessStyle.Render("90%"))
	lines = append(lines, "  Optimized DB queries")

	return strings.Join(lines, "\n")
}

// renderColorizedBio returns the bio with Tron-style coloring applied
func (m Model) renderColorizedBio() string {
	var lines []string

	// Operator header
	lines = append(lines, TitleStyle.Render("╔══════════════════════════════╗"))
	lines = append(lines, TitleStyle.Render("║  OPERATOR: PUNEET CHANDNA    ║"))
	lines = append(lines, TitleStyle.Render("╚══════════════════════════════╝"))
	lines = append(lines, "")

	// Status line with colors
	lines = append(lines, "> STATUS: "+SuccessStyle.Render("ONLINE"))
	lines = append(lines, "> ROLE: "+MenuActiveStyle.Render("SOFTWARE ENGINEER")+" | "+MenuActiveStyle.Render("FULL STACK DEV"))
	lines = append(lines, "> LOCATION: "+DimStyle.Render("VIT CHENNAI, INDIA"))
	lines = append(lines, "> GRADUATION: "+AccentTextStyle.Render("2026"))
	lines = append(lines, "")

	// Tech Arsenal - responsive grid
	lines = append(lines, TitleStyle.Render("TECH ARSENAL:"))
	lines = append(lines, "  "+MenuActiveStyle.Render("Python")+"  "+MenuActiveStyle.Render("JavaScript")+"  "+MenuActiveStyle.Render("Go"))
	lines = append(lines, "  "+DimStyle.Render("C++")+"     "+DimStyle.Render("Node.js")+"     "+DimStyle.Render("React"))
	lines = append(lines, "  "+DimStyle.Render("Next.js")+" "+DimStyle.Render("MongoDB")+"     "+DimStyle.Render("PostgreSQL"))
	lines = append(lines, "  "+AccentTextStyle.Render("AWS/GCP")+" "+AccentTextStyle.Render("Docker")+"      "+AccentTextStyle.Render("Linux/Git"))
	lines = append(lines, "")

	// Core Competencies
	lines = append(lines, TitleStyle.Render("CORE COMPETENCIES:"))
	lines = append(lines, "  • "+SuccessStyle.Render("Full Stack Dev")+" & APIs")
	lines = append(lines, "  • Data Structures & Algo")
	lines = append(lines, "  • System Design")
	lines = append(lines, "  • "+AccentTextStyle.Render("Cryptography")+" (AES)")
	lines = append(lines, "  • Cloud Computing")
	lines = append(lines, "")

	lines = append(lines, SubtitleStyle.Render("> \"Building performant &"))
	lines = append(lines, SubtitleStyle.Render("   elegant software.\""))

	return strings.Join(lines, "\n")
}

func (m Model) renderProjectList() string {
	var lines []string
	lines = append(lines, TitleStyle.Render("╔═══ PROJECT DATABASE ═══╗"))
	lines = append(lines, "")

	for i, p := range m.Projects {
		status := p.Status
		statusStyle := DimStyle
		if status == "Live" {
			statusStyle = SuccessStyle
		} else if status == "WIP" {
			statusStyle = AccentTextStyle
		}

		// Highlight selected project
		numStyle := DimStyle
		nameStyle := MenuInactiveStyle
		if i == m.ProjectIndex {
			numStyle = MenuActiveStyle
			nameStyle = MenuActiveStyle
		}

		lines = append(lines, numStyle.Render(fmt.Sprintf("%d. ", i+1))+nameStyle.Render(p.Name))
		lines = append(lines, "   "+DimStyle.Render(p.Description))
		lines = append(lines, "   Status: "+statusStyle.Render(status))
		lines = append(lines, "")
	}

	lines = append(lines, DimStyle.Render("[↑↓] Select project  [←→] Navigate"))

	return strings.Join(lines, "\n")
}

func (m Model) renderContactInfo() string {
	return TitleStyle.Render("╔═══ UPLINK TERMINAL ═══╗") + `

Establish communication via secure channels:

` + MenuActiveStyle.Render("EMAIL") + `
  puneetchandna@zohomail.in

` + MenuActiveStyle.Render("GITHUB") + `
  github.com/puneet-chandna

` + MenuActiveStyle.Render("LINKEDIN") + `
  linkedin.com/in/puneet-chandna

` + MenuActiveStyle.Render("WEBSITE") + `
  puneetchandna.com

` + DimStyle.Render(`
━━━━━━━━━━━━━━━━━━━━━━━━━━━
> Encrypted channel ready`) + `

` + AccentTextStyle.Render("Press [Enter] to open secure uplink form")
}

func (m Model) renderInspector(width, height int) string {
	var content string

	switch m.ActiveTab() {
	case "about":
		content = TitleStyle.Render("OPERATOR FILE") + "\n\n" +
			"Name: Puneet\n" +
			"Role: Engineer\n" +
			"Status: " + SuccessStyle.Render("Active")

	case "projects":
		if len(m.Projects) > 0 && m.ProjectIndex < len(m.Projects) {
			p := m.Projects[m.ProjectIndex]
			content = TitleStyle.Render("DETAILS") + "\n\n" +
				"Stack:\n"
			for _, s := range p.Stack {
				content += "  • " + s + "\n"
			}
			content += "\nYear: " + fmt.Sprint(p.Year)
			if p.URL != "" {
				content += "\n\n" + DimStyle.Render("URL:\n"+p.URL)
			}
		}

	case "contact":
		content = TitleStyle.Render("UPLINK STATUS") + "\n\n" +
			SuccessStyle.Render("● Channel Open") + "\n\n" +
			DimStyle.Render("Latency: 42ms\nEncryption: AES-256")

	case "exit":
		content = AccentTextStyle.Render("DISCONNECT?") + "\n\n" +
			DimStyle.Render("Press Enter\nto confirm")
	}

	return InspectorStyle.Width(width).Height(height).Render(content)
}
