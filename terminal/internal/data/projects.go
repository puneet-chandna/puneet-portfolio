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
			Name:        "ASCII Video Insanity",
			Description: "High-performance CLI media player rendering videos as colorized ASCII art at 30+ FPS with real-time image processing.",
			Stack:       []string{"Python", "OpenCV", "PIL", "ANSI"},
			Year:        2025,
			Status:      "Live",
			URL:         "github.com/puneet-chandna/ascii-video-insanity",
		},
		{
			Name:        "High-Performance AES Framework",
			Description: "Reduced encryption time by 51% through parallel processing. Achieved 362 KB/sec throughput with dynamic S-boxes via SHA-512.",
			Stack:       []string{"Python", "Multiprocessing", "SHA-512", "NumPy"},
			Year:        2025,
			Status:      "Live",
			URL:         "github.com/puneet-chandna/High-Performance-Parallel-AES-Encryption-Framework",
		},
		{
			Name:        "Emotion-Aware Movie Recommender",
			Description: "NLP-powered recommendation engine using SBERT embeddings. Achieved 0.71 Macro F1 score with hybrid emotional state interpretation.",
			Stack:       []string{"Python", "SBERT", "spaCy", "ML", "Pandas"},
			Year:        2025,
			Status:      "",
			URL:         "github.com/puneet-chandna/Emotion-Aware-Movie-Recommendation-System-Using-Hybrid-Emotional-States",
		},
		{
			Name:        "Crop Stress Detection Model",
			Description: "Hybrid ML pipeline combining LSTM with attention and XGBoost for binary plant-health classification using 30-day sensor data.",
			Stack:       []string{"Python", "LSTM", "XGBoost", "Scikit-learn"},
			Year:        2025,
			Status:      "",
			URL:         "",
		},
		{
			Name:        "Water Brakes",
			Description: "Streamlit app helping farmers optimize water management by analyzing contour maps for swale and trench placement.",
			Stack:       []string{"Streamlit", "Python", "Plotly", "Matplotlib"},
			Year:        2024,
			Status:      "Live",
			URL:         "water-brakes.streamlit.app/",
		},
		{
			Name:        "Terminal Portfolio",
			Description: "This SSH-based portfolio you're viewing now. Built with Go and Charm libraries featuring Tron aesthetics.",
			Stack:       []string{"Go", "Bubble Tea", "Lipgloss", "Wish"},
			Year:        2025,
			Status:      "Live",
			URL:         "ssh puneet.space",
		},
		{
			Name:        "Road Safety Expert System",
			Description: "AI-Powered Expert System that analyzes user-reported road safety issues and provides detailed, guideline-based intervention recommendations. Built on Google Gemini.",
			Stack:       []string{"React", "Vite", "TailwindCSS", "Node.js", "Vercel", "Google Gemini"},
			Year:        2025,
			Status:      "Live",
			URL:         "road-safety-gpt-cyan.vercel.app",
		},
	}
}
