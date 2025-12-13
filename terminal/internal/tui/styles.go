package tui

import "github.com/charmbracelet/lipgloss"

// Tron Theme: Deep space, cyan glow, electric blue accents

// Core palette colors
var (
	ColorBg        = lipgloss.Color("#0a0a0f")
	ColorPrimary   = lipgloss.Color("#00d4ff") // Cyan glow
	ColorSecondary = lipgloss.Color("#0066ff") // Electric blue
	ColorAccent    = lipgloss.Color("#ff6600") // Orange warning
	ColorDim       = lipgloss.Color("#7a7a8a") // Muted grey (lighter for readability)
	ColorText      = lipgloss.Color("#e0e0e0") // Soft white
	ColorSuccess   = lipgloss.Color("#00ff88") // Green
)

// Styles holds all the lipgloss styles for a session
// Using a struct allows per-session renderer support for SSH
type Styles struct {
	Base              lipgloss.Style
	Header            lipgloss.Style
	Footer            lipgloss.Style
	Menu              lipgloss.Style
	MenuActive        lipgloss.Style
	MenuInactive      lipgloss.Style
	Viewport          lipgloss.Style
	Inspector         lipgloss.Style
	Title             lipgloss.Style
	Subtitle          lipgloss.Style
	AccentText        lipgloss.Style
	Success           lipgloss.Style
	Dim               lipgloss.Style
	BootSpinner       lipgloss.Style
	BootProgress      lipgloss.Style
	BootAccessGranted lipgloss.Style
}

// NewStyles creates styles using the given renderer
// For SSH sessions, pass bubbletea.MakeRenderer(s) to get correct colors
func NewStyles(r *lipgloss.Renderer) Styles {
	return Styles{
		Base: r.NewStyle().
			Background(ColorBg).
			Foreground(ColorText),

		Header: r.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Background(lipgloss.Color("#0d0d14")).
			Padding(0, 1),

		Footer: r.NewStyle().
			Foreground(ColorDim).
			Background(lipgloss.Color("#0d0d14")).
			Padding(0, 1),

		Menu: r.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorDim).
			Padding(1, 2),

		MenuActive: r.NewStyle().
			Bold(true).
			Foreground(ColorPrimary),

		MenuInactive: r.NewStyle().
			Foreground(ColorDim),

		Viewport: r.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSecondary).
			Padding(1, 2),

		Inspector: r.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorDim).
			Padding(1, 2),

		Title: r.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			MarginBottom(1),

		Subtitle: r.NewStyle().
			Foreground(ColorSecondary).
			Italic(true),

		AccentText: r.NewStyle().
			Foreground(ColorAccent),

		Success: r.NewStyle().
			Bold(true).
			Foreground(ColorSuccess),

		Dim: r.NewStyle().
			Foreground(ColorDim),

		BootSpinner: r.NewStyle().
			Foreground(ColorPrimary),

		BootProgress: r.NewStyle().
			Foreground(ColorSecondary),

		BootAccessGranted: r.NewStyle().
			Bold(true).
			Foreground(ColorSuccess).
			Background(lipgloss.Color("#0d1a0d")).
			Padding(1, 4),
	}
}

// DefaultStyles returns styles using the default renderer
// Only use for local testing, not for SSH sessions
func DefaultStyles() Styles {
	return NewStyles(lipgloss.DefaultRenderer())
}

// Legacy package-level styles for backward compatibility during refactor
// These will be removed once all code uses Styles struct
var (
	BaseStyle              = lipgloss.NewStyle().Background(ColorBg).Foreground(ColorText)
	HeaderStyle            = lipgloss.NewStyle().Bold(true).Foreground(ColorPrimary).Background(lipgloss.Color("#0d0d14")).Padding(0, 1)
	FooterStyle            = lipgloss.NewStyle().Foreground(ColorDim).Background(lipgloss.Color("#0d0d14")).Padding(0, 1)
	MenuStyle              = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(ColorDim).Padding(1, 2)
	MenuActiveStyle        = lipgloss.NewStyle().Bold(true).Foreground(ColorPrimary)
	MenuInactiveStyle      = lipgloss.NewStyle().Foreground(ColorDim)
	ViewportStyle          = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(ColorSecondary).Padding(1, 2)
	InspectorStyle         = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(ColorDim).Padding(1, 2)
	TitleStyle             = lipgloss.NewStyle().Bold(true).Foreground(ColorPrimary).MarginBottom(1)
	SubtitleStyle          = lipgloss.NewStyle().Foreground(ColorSecondary).Italic(true)
	AccentTextStyle        = lipgloss.NewStyle().Foreground(ColorAccent)
	SuccessStyle           = lipgloss.NewStyle().Bold(true).Foreground(ColorSuccess)
	DimStyle               = lipgloss.NewStyle().Foreground(ColorDim)
	BootSpinnerStyle       = lipgloss.NewStyle().Foreground(ColorPrimary)
	BootProgressStyle      = lipgloss.NewStyle().Foreground(ColorSecondary)
	BootAccessGrantedStyle = lipgloss.NewStyle().Bold(true).Foreground(ColorSuccess).Background(lipgloss.Color("#0d1a0d")).Padding(1, 4)
)
