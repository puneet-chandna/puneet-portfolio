package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var encryptChars = []string{"█", "▓", "▒", "░", "▪", "▫", "◆", "◇", "○", "●", "◐", "◑", "◒", "◓"}

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func (m Model) View() string {
	var content string
	switch m.State {
	case StateBoot:
		content = m.renderBoot()
	case StateLoading:
		content = m.renderLoading()
	case StateAccessGranted:
		content = m.renderAccessGranted()
	case StateContactForm:
		content = m.renderContactForm()
	case StateContactSending:
		content = m.renderContactSending()
	case StateContactSent:
		content = m.renderContactSent()
	default:
		content = m.renderMain()
	}
	return lipgloss.NewStyle().MaxWidth(m.Width).MaxHeight(m.Height).Render(content)
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
	compact := m.Width < 70 || m.Height < 24
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorDim).
		Padding(0, 1)
	formWidth := m.contactFormWidth()
	inputStyle = inputStyle.Width(max(0, formWidth-inputStyle.GetHorizontalBorderSize()))

	focusedInputStyle := inputStyle.BorderForeground(ColorPrimary)

	labels := []string{"CALLSIGN", "COMM FREQ (Email)", "MESSAGE"}

	var inputs []string
	for i, input := range m.ContactInputs {
		if compact && i != m.ContactFocus {
			continue
		}
		input.Width = m.contactInputWidth(input)
		style := inputStyle
		if i == m.ContactFocus {
			style = focusedInputStyle
		}
		label := m.Styles.Dim.Render(labels[i])
		inputs = append(inputs, label+"\n"+style.Render(input.View()))
	}

	formContent := lipgloss.JoinVertical(lipgloss.Left, inputs...)
	if compact {
		formContent = m.Styles.Dim.Render(fmt.Sprintf("Field %d/%d", m.ContactFocus+1, len(m.ContactInputs))) + "\n" + formContent
	}

	parts := []string{m.Styles.Title.MarginBottom(0).Render("╔═══ SECURE UPLINK ═══╗")}
	if !compact {
		parts = append(parts, m.Styles.Subtitle.Render("Establish direct communication channel"))
	}
	parts = append(parts, formContent)
	if m.ContactError != "" {
		parts = append(parts, m.Styles.AccentText.Width(formWidth).Render(m.ContactError))
	}
	parts = append(parts, m.Styles.Dim.Width(max(1, m.Width-2)).Render("[Tab] Next  [Enter] Transmit  [Esc] Abort"))
	content := lipgloss.JoinVertical(lipgloss.Center, parts...)

	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) renderContactSending() string {
	var bar strings.Builder
	// Keep the original animated glyphs, without inventing delivery progress.
	for i := 0; i < min(30, max(1, m.Width-4)); i++ {
		bar.WriteString(m.Styles.BootSpinner.Render(encryptChars[(m.SendingFrame+i)%len(encryptChars)]))
	}
	content := lipgloss.JoinVertical(lipgloss.Center,
		m.Styles.Title.MarginBottom(0).Width(max(1, m.Width-2)).Render("▓▓▓ TRANSMISSION IN PROGRESS ▓▓▓"),
		"", "["+bar.String()+"]", "",
		m.Styles.Dim.Width(max(1, m.Width-2)).Render("Waiting for the email provider..."),
		m.Styles.Dim.Render("[Ctrl+C] Disconnect"),
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
		m.Styles.Dim.Render("Accepted for delivery by the email provider."),
		"",
		m.Styles.Subtitle.Render("Press any key to continue..."),
		"",
	)

	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
}

func (m Model) renderMain() string {
	header := m.GetHeader()
	footer := m.GetFooter()

	menuWidth, inspectorWidth, viewportWidth, contentHeight := m.paneLayout()
	viewport := m.renderViewport(viewportWidth, contentHeight)

	panes := viewport
	if menuWidth > 0 {
		menu := m.renderMenu(menuWidth, contentHeight)
		panes = lipgloss.JoinHorizontal(lipgloss.Top, menu, viewport)
	}
	if inspectorWidth > 0 {
		inspector := m.renderInspector(inspectorWidth, contentHeight)
		panes = lipgloss.JoinHorizontal(lipgloss.Top, panes, inspector)
	}

	// Show help overlay if enabled
	if m.ShowHelp {
		return m.renderHelpOverlay()
	}

	return lipgloss.NewStyle().MaxWidth(m.Width).MaxHeight(m.Height).Render(lipgloss.JoinVertical(lipgloss.Left, header, panes, footer))
}

func (m Model) paneLayout() (menuWidth, inspectorWidth, viewportWidth, contentHeight int) {
	menuWidth = 22
	if m.Width < 60 {
		menuWidth = 16
	}
	if m.Width < 38 {
		menuWidth = 0
	}
	contentHeight = m.Height - lipgloss.Height(m.GetHeader()) - lipgloss.Height(m.GetFooter())
	if contentHeight < 1 {
		contentHeight = 1
	}
	if m.Width > 100 && contentHeight >= 28 {
		inspectorWidth = 26
	}
	viewportWidth = m.Width - menuWidth - inspectorWidth
	if viewportWidth < 1 {
		viewportWidth = 1
	}
	return
}

func (m *Model) refreshViewport(reset bool) {
	_, _, width, height := m.paneLayout()
	frameWidth := m.Styles.Viewport.GetHorizontalFrameSize()
	frameHeight := m.Styles.Viewport.GetVerticalFrameSize()
	m.Viewport.Width = max(1, width-frameWidth)
	m.Viewport.Height = max(1, height-frameHeight)
	m.Viewport.SetContent(lipgloss.NewStyle().Width(m.Viewport.Width).Render(m.viewportContent()))
	m.viewportTab = m.ActiveTab()
	if reset {
		m.Viewport.GotoTop()
	}
}

func (m Model) viewportContent() string {
	switch m.ActiveTab() {
	case "about":
		return m.renderColorizedBio()
	case "experience":
		return m.renderExperience()
	case "projects":
		return m.renderProjectList()
	case "contact":
		return m.renderContactInfo()
	case "exit":
		return "Press Enter to disconnect from PUNEET-OS\n\nThank you for visiting."
	}
	return ""
}

func (m Model) renderHelpOverlay() string {
	if m.Width < 60 || m.Height < 28 {
		help := lipgloss.JoinVertical(lipgloss.Left,
			m.Styles.Title.MarginBottom(0).Render("HELP"),
			m.Styles.Dim.Render("↑↓/jk navigate or select"),
			m.Styles.Dim.Render("h/l change section"),
			m.Styles.Dim.Render("PgUp/PgDn scroll"),
			m.Styles.Dim.Render("Enter select/send"),
			m.Styles.Dim.Render("? Close  q Exit"),
		)
		box := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(ColorPrimary).Padding(0, 1).Background(ColorBg)
		return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, box.Render(help))
	}
	helpContent := lipgloss.JoinVertical(lipgloss.Left,
		m.Styles.Title.Render("╔═══ KEYBOARD SHORTCUTS ═══╗"),
		"",
		m.Styles.MenuActive.Render("Navigation"),
		m.Styles.Dim.Render("  ↑/k     Menu or project"),
		m.Styles.Dim.Render("  ↓/j     Menu or project"),
		m.Styles.Dim.Render("  ←/h     Previous section"),
		m.Styles.Dim.Render("  →/l     Next section"),
		m.Styles.Dim.Render("  PgUp/PgDn Scroll page"),
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

	return m.Styles.Menu.Width(max(0, width-m.Styles.Menu.GetHorizontalBorderSize())).Height(max(0, height-m.Styles.Menu.GetVerticalBorderSize())).Render(menu)
}

func (m Model) renderViewport(width, height int) string {
	viewport := m.Viewport
	if m.viewportTab != m.ActiveTab() {
		viewport.SetContent(lipgloss.NewStyle().Width(viewport.Width).Render(m.viewportContent()))
		viewport.GotoTop()
	}
	return m.Styles.Viewport.Width(max(0, width-m.Styles.Viewport.GetHorizontalBorderSize())).Height(max(0, height-m.Styles.Viewport.GetVerticalBorderSize())).Render(viewport.View())
}

// renderExperience renders the experience/work history section
func (m Model) renderExperience() string {
	lines := strings.Split(m.Experience, "\n")
	for i, line := range lines {
		switch {
		case strings.TrimSpace(line) == "MISSION LOG:":
			lines[i] = m.Styles.Title.Render("MISSION LOG")
		case strings.HasPrefix(line, "• "):
			parts := strings.SplitN(line, " @ ", 2)
			lines[i] = m.Styles.MenuActive.Render(parts[0])
			if len(parts) == 2 {
				lines[i] += " @ " + parts[1]
			}
		case strings.Contains(line, "202"):
			lines[i] = m.Styles.Dim.Render(line)
		default:
			lines[i] = strings.NewReplacer("30%", m.Styles.Success.Render("30%"), "90%", m.Styles.Success.Render("90%")).Replace(line)
		}
	}
	return strings.Join(lines, "\n")
}

// renderColorizedBio returns the bio with Tron-style coloring applied
func (m Model) renderColorizedBio() string {
	title := m.Styles.Title.MarginBottom(0)
	header := title.Render("╔══════════════════════════════╗\n║  OPERATOR: PUNEET CHANDNA    ║\n╚══════════════════════════════╝")
	if m.Viewport.Width < 32 {
		header = title.Render("OPERATOR: PUNEET CHANDNA")
	}
	lines := strings.Split(m.Bio, "\n")
	for i, line := range lines {
		switch {
		case strings.HasPrefix(line, "> STATUS:"):
			lines[i] = strings.ReplaceAll(line, "ONLINE", m.Styles.Success.Render("ONLINE"))
		case strings.HasPrefix(line, "> ROLE:") || strings.HasPrefix(line, "  Python"):
			lines[i] = m.Styles.MenuActive.Render(line)
		case strings.HasPrefix(line, "> LOCATION:"):
			lines[i] = m.Styles.Dim.Render(line)
		case strings.Contains(line, "GRADUATION:") || strings.Contains(line, "AWS/GCP") || strings.Contains(line, "Cryptography"):
			lines[i] = m.Styles.AccentText.Render(line)
		case strings.HasSuffix(strings.TrimSpace(line), ":"):
			lines[i] = m.Styles.Title.Render(line)
		case strings.Contains(line, "\""):
			lines[i] = m.Styles.Subtitle.Render(line)
		}
	}
	return header + "\n\n" + strings.Join(lines, "\n")
}

func (m Model) renderProjectList() string {
	var lines []string
	lines = append(lines, m.Styles.Title.Render("╔═══ PROJECT DATABASE ═══╗"))
	if len(m.Projects) > 0 && m.ProjectIndex >= 0 && m.ProjectIndex < len(m.Projects) {
		p := m.Projects[m.ProjectIndex]
		lines = append(lines, m.Styles.MenuActive.Render("SELECTED: "+p.Name), p.Description)
		lines = append(lines, "Stack: "+strings.Join(p.Stack, ", "))
		if p.URL != "" {
			lines = append(lines, "URL:", p.URL)
		}
	}
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

		lines = append(lines, numStyle.Render(fmt.Sprintf("%d. ", i+1))+nameStyle.Render(p.Name)+" "+statusStyle.Render(status))
	}

	lines = append(lines, "", m.Styles.Dim.Render("[↑↓/jk] Select project  [PgUp/PgDn] Scroll"))

	return strings.Join(lines, "\n")
}

func (m Model) renderContactInfo() string {
	status := "Press [Enter] to open the message form"
	if !ContactFormConfigured() {
		status = "Message form is disabled until email delivery is configured. Use these links instead."
	}
	if m.ContactError != "" {
		status = m.ContactError
	}
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

` + m.Styles.Dim.Render("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━\n> "+status) + `

` + m.Styles.AccentText.Render("Resume: https://puneetchandna.com/Puneet-Chandna-Resume.pdf")
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
		if len(m.Projects) > 0 && m.ProjectIndex >= 0 && m.ProjectIndex < len(m.Projects) {
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
			m.Styles.Dim.Render("Use Enter to open the message form when configured.")

	case "exit":
		content = m.Styles.AccentText.Render("DISCONNECT?") + "\n\n" +
			m.Styles.Dim.Render("Press Enter\nto confirm")
	}

	return m.Styles.Inspector.Width(max(0, width-m.Styles.Inspector.GetHorizontalBorderSize())).Height(max(0, height-m.Styles.Inspector.GetVerticalBorderSize())).Render(content)
}
