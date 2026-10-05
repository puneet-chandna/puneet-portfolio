import { 
  SiPython, SiJavascript, SiGo, SiCplusplus, SiReact, 
  SiNodedotjs, SiNextdotjs, SiMongodb, SiPostgresql, 
  SiDocker, SiLinux, SiGit,
  SiThreedotjs, SiPytorch, SiGooglecloud
} from 'react-icons/si'
import { FaAws } from 'react-icons/fa'

export const projects = [
  {
    id: 7,
    name: "Real-Time Dealer Gamma Exposure (GEX) Analysis",
    description: "Real-time quantitative analysis platform for SPX 0DTE options. Computes dealer gamma exposure, zero-gamma levels and Charm/Vanna hedging-flow estimates, with Kalman signal smoothing and Hawkes-style flow modeling. Combines live market feeds, historical replay and strategy backtesting to study intraday market structure.",
    stack: ["FastAPI", "Next.js", "PostgreSQL", "WebSocket"],
    year: 2026,
    status: "Live",
    url: "https://github.com/puneet-chandna/0DTE-dealer-gamma",
    image: "/projects/gex.webp"
  },
  {
    id: 8,
    name: "Lattora",
    description: "Terminal workbench for reproducible VM placement experiments with CloudSim Plus. Compare four placement algorithms, inspect retained results and independently validate evidence, with offline packages for Linux and Apple Silicon macOS.",
    stack: ["Python", "Textual", "Java 21", "CloudSim Plus"],
    year: 2026,
    status: "Live",
    url: "https://github.com/puneet-chandna/Lattora",
    image: "/projects/lattora.png",
    imageFit: "contain"
  },
  {
    id: 9,
    name: "Requests Native",
    description: "HTTP library pairing the familiar Python Requests API with a Rust transport via PyO3, plus native async and blocking Rust clients. Built around preserving Requests behavior.",
    stack: ["Rust", "Python", "PyO3", "Tokio"],
    year: 2026,
    status: "Beta",
    url: "https://github.com/puneet-chandna/requests-native",
    image: "/projects/requests-native.png"
  },
  {
    id: 1,
    name: "ASCII Video Insanity",
    description: "Terminal video player that turns MP4 frames into color ASCII art. Uses OpenCV and NumPy with adjustable character sets, quality presets and terminal-aware sizing.",
    stack: ["Python", "OpenCV", "NumPy", "ANSI"],
    year: 2025,
    status: "Live",
    url: "https://github.com/puneet-chandna/ascii-video-insanity",
    image: "/projects/ascii.webp"
  },
  {
    id: 3,
    name: "Emotion-Aware Movie Recommender",
    description: "Exploratory NLP notebook that uses movie overviews, sentiment and Sentence-BERT embeddings to assign basic and hybrid mood labels, with visualizations of the resulting groups.",
    stack: ["Python", "Sentence-BERT", "spaCy", "Jupyter"],
    year: 2025,
    status: "Live",
    url: "https://github.com/puneet-chandna/Emotion-Aware-Movie-Recommendation-System-Using-Hybrid-Emotional-States",
    image: "/projects/movie.webp"
  },
  {
    id: 4,
    name: "Crop Stress Detection Model",
    description: "Plant-health classifier combining an attention-based LSTM with XGBoost to detect crop stress from 30-day sensor sequences.",
    stack: ["Python", "LSTM", "XGBoost", "Scikit-learn"],
    year: 2025,
    status: "Live",
    url: "",
    image: "/projects/crop.webp"
  },
  {
    id: 5,
    name: "Water Brakes",
    description: "Streamlit tool for exploring contour maps and planning swale and trench placement to support farm water management.",
    stack: ["Streamlit", "Python", "Plotly"],
    year: 2024,
    status: "Live",
    url: "https://water-brakes.streamlit.app/",
    image: "/projects/water.webp"
  }
];

export const experience = [
  {
    id: 4,
    title: "Product Developer",
    company: "Hyr.works (Zofa AI Solutions Pvt. Ltd.)",
    date: "Mar 2026 – Present",
    description: "Developed and stabilized Hyr’s V2 recruiter platform in Next.js, resolving recurring client-reported bugs across frontend and backend workflows. Built the Hyr Live Chrome extension and integrated Hyr Agent APIs. Hardened tenant isolation across Next.js and Java services using Supabase row-level security, JWT authentication, and role-based access controls, with integration testing for authentication and workflow regressions. Rebuilt the FFmpeg recording pipeline with dual recorders, 2-minute chunks, browser buffering, and verified uploads to preserve partial interviews during network failures. Implemented Grafana/Sentry monitoring and Telegram alerts across DigitalOcean production services. Independently researched, designed, and built hyr.works in Next.js, exploring 16+ design prototypes and implementing responsive pages, technical SEO, and generative engine optimization (GEO)."
  },
  {
    id: 0,
    title: "Backend Engineering Intern",
    company: "Apoliums Infotech India Pvt. Ltd.",
    date: "Dec 2025 – Jan 2026",
    description: "Migrated legacy Node.js services to Golang (Gin) and optimized MySQL schemas, reducing API latency and enhancing concurrency via clean architecture. Engineered high-performance RESTful APIs deployed on GCP (Compute Engine, Cloud SQL), ensuring robust system reliability and seamless production operations."
  },
  {
    id: 1,
    title: "Research Intern",
    company: "Centre for e-Automation Technologies (CeAT), VIT Chennai",
    date: "May 2025 – July 2025",
    description: "Developed CloudSim Plus framework for VM placement optimization using Hippopotamus Optimization algorithm. Authored research paper on simulation results."
  },
  {
    id: 2,
    title: "Full Stack Developer Intern",
    company: "Daira Edtech Pvt Limited",
    date: "Dec 2024 – Feb 2025",
    description: "Built RESTful APIs for Vidhira EdTech platform. Designed data models leading to 30% improvement in data retrieval efficiency."
  },
  {
    id: 3,
    title: "Web Development Intern",
    company: "Indian Institute of Technology Bombay (IIT Bombay)",
    date: "Sept 2024 - Oct 2024",
    description: "Reduced data payload by 90% for low-bandwidth environments. Optimized backend database queries and implemented Joi validation."
  }
];

export const skills = [
  { name: "Python", icon: SiPython, color: "#3776AB" },
  { name: "JavaScript", icon: SiJavascript, color: "#F7DF1E" },
  { name: "Go", icon: SiGo, color: "#00ADD8" },
  { name: "C++", icon: SiCplusplus, color: "#00599C" },
  { name: "React", icon: SiReact, color: "#61DAFB" },
  { name: "Node.js", icon: SiNodedotjs, color: "#339933" },
  { name: "Next.js", icon: SiNextdotjs, color: "#FFFFFF" },
  { name: "MongoDB", icon: SiMongodb, color: "#47A248" },
  { name: "PostgreSQL", icon: SiPostgresql, color: "#4169E1" },
  { name: "AWS", icon: FaAws, color: "#FF9900" },
  { name: "Google Cloud", icon: SiGooglecloud, color: "#4285F4" },
  { name: "Docker", icon: SiDocker, color: "#2496ED" },
  { name: "Linux", icon: SiLinux, color: "#FCC624" },
  { name: "Git", icon: SiGit, color: "#F05032" },
  { name: "Three.js", icon: SiThreedotjs, color: "#FFFFFF" },
  { name: "ML/PyTorch", icon: SiPytorch, color: "#EE4C2C" }
];

export const featuredCert = {
  name: "Professional Cloud Architect",
  issuer: "Google Cloud",
  type: "Certification",
  badge: "/professional-cloud-architect-certification.png",
  url: "https://www.credly.com/badges/5be24c81-d882-476b-bd82-0504a3ab46f3/public_url"
};

export const courseCerts = [
  { name: "Introduction to Cybersecurity", issuer: "Cisco", color: "#049fd9" },
  { name: "Introduction to Packet Tracer", issuer: "Cisco", color: "#049fd9" },
  { name: "Load Balancing on Compute Engine", issuer: "Google Cloud", color: "#4285F4" },
  { name: "C++ Intermediate", issuer: "Sololearn", color: "#41c473" },
  { name: "HTML and CSS in Depth", issuer: "Coursera", color: "#0056d2" },
  { name: "Programming with JavaScript", issuer: "Coursera", color: "#0056d2" },
  { name: "React Basics", issuer: "Coursera", color: "#0056d2" },
  { name: "SQL Course", issuer: "Coursera", color: "#0056d2" }
];

export const credlyProfile = "https://www.credly.com/users/puneet-chandna.ea5ad5b2";

export const socialLinks = {
  github: "https://github.com/puneet-chandna",
  linkedin: "https://www.linkedin.com/in/puneet-chandna",
  email: "puneetchandna@zohomail.in",
  website: "https://puneetchandna.com"
};
