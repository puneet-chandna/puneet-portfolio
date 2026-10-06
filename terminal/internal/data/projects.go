package data

type Project struct {
	Name        string
	Description string
	Stack       []string
	Year        int
	Status      string
	URL         string
}

func GetProjects() []Project {
	return []Project{
		{
			Name:        "Real-Time Dealer Gamma Exposure (GEX) Analysis",
			Description: "Real-time quantitative analysis platform for SPX 0DTE options. Computes dealer gamma exposure, zero-gamma levels and Charm/Vanna hedging-flow estimates, with Kalman signal smoothing and Hawkes-style flow modeling. Combines live market feeds, historical replay and strategy backtesting to study intraday market structure.",
			Stack:       []string{"FastAPI", "Next.js", "PostgreSQL", "WebSocket"},
			Year:        2026,
			Status:      "Live",
			URL:         "https://github.com/puneet-chandna/0DTE-dealer-gamma",
		},
		{
			Name:        "Lattora",
			Description: "A terminal workbench for studying how virtual machines are placed across cloud servers. Compare four placement algorithms, follow experiments as they run, and revisit the inputs, results and independent validation behind each comparison, with offline packages for Linux and Apple Silicon macOS.",
			Stack:       []string{"Python", "Textual", "Java 21", "CloudSim Plus"},
			Year:        2026,
			Status:      "Live",
			URL:         "https://github.com/puneet-chandna/Lattora",
		},
		{
			Name:        "Requests Native",
			Description: "A Rust rewrite of Python’s Requests library that keeps its familiar API. Built around preserving the behaviors existing applications depend on—sessions, streaming, errors and custom adapters—with one shared HTTP engine for Python and Rust.",
			Stack:       []string{"Rust", "Python", "PyO3", "Tokio"},
			Year:        2026,
			Status:      "Beta",
			URL:         "https://github.com/puneet-chandna/requests-native",
		},
		{
			Name:        "ASCII Video Insanity",
			Description: "Play videos directly in your terminal as moving, full-color character art. Converts each frame into text, with different character sets and quality controls to balance visual detail with playback performance.",
			Stack:       []string{"Python", "OpenCV", "NumPy", "ANSI"},
			Year:        2025,
			Status:      "Live",
			URL:         "https://github.com/puneet-chandna/ascii-video-insanity",
		},
		{
			Name:        "Emotion-Aware Movie Recommender",
			Description: "Explores how movies can be grouped by emotional tone. Uses plot descriptions to identify basic and mixed moods, then visualizes those groupings in a research notebook.",
			Stack:       []string{"Python", "Sentence-BERT", "spaCy", "Jupyter"},
			Year:        2025,
			Status:      "Live",
			URL:         "https://github.com/puneet-chandna/Emotion-Aware-Movie-Recommendation-System-Using-Hybrid-Emotional-States",
		},
		{
			Name:        "Crop Stress Detection Model",
			Description: "Studies crop health through patterns in 30 days of sensor readings. Combines an attention-based LSTM model that follows changes over time with a second XGBoost classifier to distinguish healthy plants from stressed ones.",
			Stack:       []string{"Python", "LSTM", "XGBoost", "Scikit-learn", "Attention"},
			Year:        2025,
			Status:      "Live",
			URL:         "",
		},
		{
			Name:        "Water Brakes",
			Description: "Turns terrain survey data into interactive maps for exploring where water could collect or drain. Groups areas by slope and elevation to suggest swale and trench locations, with downloadable data and reports for planning.",
			Stack:       []string{"Streamlit", "Python", "Plotly"},
			Year:        2024,
			Status:      "Live",
			URL:         "https://water-brakes.streamlit.app/",
		},
		{
			Name:        "Terminal Portfolio",
			Description: "Explore my work, projects and experience through an interactive SSH portfolio. Built with Go and Charm libraries, with keyboard navigation and a Tron-inspired terminal interface.",
			Stack:       []string{"Go", "Bubble Tea", "Lipgloss", "Wish"},
			Year:        2025,
			Status:      "Live",
			URL:         "ssh puneet.space",
		},
	}
}
