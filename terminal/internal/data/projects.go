package data

type Project struct {
	Name        string
	Description string
	Stack       []string
	Year        int
	Status      string // "Live", "WIP", "Archived"
	URL         string
}

func GetProjects() []Project {
	return []Project{
		{
			Name:        "Terminal Portfolio",
			Description: "A custom SSH server portfolio built with Charm libraries. Features multi-pane TUI, boot sequence animation, and Tron aesthetics.",
			Stack:       []string{"Go", "Bubble Tea", "Lipgloss", "Wish"},
			Year:        2025,
			Status:      "Live",
			URL:         "ssh puneet.sh",
		},
		{
			Name:        "Portfolio Website",
			Description: "Personal portfolio with pastel aesthetics, interactive cat companions, and smooth animations.",
			Stack:       []string{"React", "Vite", "CSS3"},
			Year:        2024,
			Status:      "Live",
			URL:         "puneetchandna.com",
		},
		{
			Name:        "Gemini Chat App",
			Description: "AI chat application with continuous conversation history and modern UI.",
			Stack:       []string{"React", "Gemini API"},
			Year:        2025,
			Status:      "WIP",
			URL:         "",
		},
	}
}
