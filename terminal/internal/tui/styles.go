package tui

import "github.com/charmbracelet/lipgloss"

// Tron Theme: Deep space, cyan glow, electric blue accents
var (
	// Core palette
	ColorBg        = lipgloss.Color("#0a0a0f")
	ColorPrimary   = lipgloss.Color("#00d4ff") // Cyan glow
	ColorSecondary = lipgloss.Color("#0066ff") // Electric blue
	ColorAccent    = lipgloss.Color("#ff6600") // Orange warning
	ColorDim       = lipgloss.Color("#7a7a8a") // Muted grey (lighter for readability)
	ColorText      = lipgloss.Color("#e0e0e0") // Soft white
	ColorSuccess   = lipgloss.Color("#00ff88") // Green

	// Styles
	BaseStyle = lipgloss.NewStyle().
			Background(ColorBg).
			Foreground(ColorText)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Background(lipgloss.Color("#0d0d14")).
			Padding(0, 1)

	FooterStyle = lipgloss.NewStyle().
			Foreground(ColorDim).
			Background(lipgloss.Color("#0d0d14")).
			Padding(0, 1)

	MenuStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorDim).
			Padding(1, 2)

	MenuActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary)

	MenuInactiveStyle = lipgloss.NewStyle().
				Foreground(ColorDim)

	ViewportStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSecondary).
			Padding(1, 2)

	InspectorStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorDim).
			Padding(1, 2)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Italic(true)

	AccentTextStyle = lipgloss.NewStyle().
			Foreground(ColorAccent)

	SuccessStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSuccess)

	DimStyle = lipgloss.NewStyle().
			Foreground(ColorDim)

	// Boot sequence styles
	BootSpinnerStyle = lipgloss.NewStyle().
				Foreground(ColorPrimary)

	BootProgressStyle = lipgloss.NewStyle().
				Foreground(ColorSecondary)

	BootAccessGrantedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorSuccess).
				Background(lipgloss.Color("#0d1a0d")).
				Padding(1, 4)
)
