import { 
  SiPython, SiJavascript, SiGo, SiCplusplus, SiReact, 
  SiNodedotjs, SiNextdotjs, SiMongodb, SiPostgresql, 
  SiAmazonwebservices, SiDocker, SiLinux, SiGit, 
  SiThreedotjs, SiPytorch, SiGooglecloud
} from 'react-icons/si'

export const projects = [
  {
    id: 7,
    name: "Real-Time Dealer Gamma Exposure (GEX) Analysis",
    description: "Full-stack 0DTE SPX monitor computing strike-level dealer GEX, zero-gamma, Charm/Vanna flow signals, and Kalman-smoothed regimes from live option chains.",
    stack: ["FastAPI", "Next.js", "PostgreSQL", "WebSocket"],
    year: 2026,
    status: "Live",
    url: "https://github.com/puneet-chandna/0DTE-dealer-gamma",
    image: "/projects/gex.webp"
  },
  {
    id: 8,
    name: "CloudSim-HO Research",
    description: "Centre for e-Automation Technologies (CeAT) research simulator implementing Hippopotamus Optimization for VM placement with CloudSim Plus benchmarking and statistical analysis.",
    stack: ["Java 21", "Maven", "CloudSim Plus", "ANOVA"],
    year: 2025,
    status: "Research",
    url: "https://github.com/puneet-chandna/cloudsim-ho-research",
    image: "/projects/cloudsim-ho.webp"
  },
  {
    id: 1,
    name: "ASCII Video Insanity",
    description: "High-performance CLI media player rendering videos as colorized ASCII art at 30+ FPS with real-time image processing.",
    stack: ["Python", "OpenCV", "PIL", "ANSI"],
    year: 2025,
    status: "Live",
    url: "https://github.com/puneet-chandna/ascii-video-insanity",
    image: "/projects/ascii.webp"
  },
  {
    id: 2,
    name: "Encryption Algorithm Using Dynamic S-boxes",
    description: "Modified the AES algorithm to reduce encryption time by 51% through parallel processing. Achieved higher throughput  and robust security with dynamic S-boxes via SHA-512.",
    stack: ["Python", "Multiprocessing", "SHA-512", "NumPy"],
    year: 2025,
    status: "Live",
    url: "",
    image: "/projects/aes.webp"
  },
  {
    id: 3,
    name: "Emotion-Aware Movie Recommender",
    description: "NLP-powered recommendation engine using SBERT embeddings. Achieved 0.71 Macro F1 score with hybrid emotional state interpretation.",
    stack: ["Python", "SBERT", "spaCy", "ML"],
    year: 2025,
    status: "Live",
    url: "https://github.com/puneet-chandna/Emotion-Aware-Movie-Recommendation-System-Using-Hybrid-Emotional-States",
    image: "/projects/movie.webp"
  },
  {
    id: 4,
    name: "Crop Stress Detection Model",
    description: "Hybrid ML pipeline combining LSTM with attention and XGBoost for binary plant-health classification using 30-day sensor data.",
    stack: ["Python", "LSTM", "XGBoost", "Scikit-learn"],
    year: 2025,
    status: "Live",
    url: "",
    image: "/projects/crop.webp"
  },
  {
    id: 5,
    name: "Water Brakes",
    description: "Streamlit app helping farmers optimize water management by analyzing contour maps for swale and trench placement.",
    stack: ["Streamlit", "Python", "Plotly"],
    year: 2024,
    status: "Live",
    url: "https://water-brakes.streamlit.app/",
    image: "/projects/water.webp"
  },
  {
    id: 6,
    name: "BookMyEvent",
    description: "Event management app processing ~10,000 ticket authentications per event with real-time QR code scanning via ESP32.",
    stack: ["Node.js", "Express", "Flutter", "MongoDB"],
    year: 2024,
    status: "Live",
    url: "https://github.com/puneet-chandna/bookmyshow_server",
    image: "/projects/event.webp"
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
  { name: "AWS", icon: SiAmazonwebservices, color: "#FF9900" },
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
