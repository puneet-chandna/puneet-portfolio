import { 
  SiPython, SiJavascript, SiGo, SiCplusplus, SiReact, 
  SiNodedotjs, SiNextdotjs, SiMongodb, SiPostgresql, 
  SiAmazonwebservices, SiDocker, SiLinux, SiGit, 
  SiThreedotjs, SiPytorch, SiGooglecloud
} from 'react-icons/si'

export const projects = [
  {
    id: 1,
    name: "ASCII Video Insanity",
    description: "High-performance CLI media player rendering videos as colorized ASCII art at 30+ FPS with real-time image processing.",
    stack: ["Python", "OpenCV", "PIL", "ANSI"],
    year: 2025,
    status: "Live",
    url: "https://github.com/puneet-chandna",
    image: "/projects/ascii.png"
  },
  {
    id: 2,
    name: "High-Performance AES Framework",
    description: "Reduced encryption time by 51% through parallel processing. Achieved 362 KB/sec throughput with dynamic S-boxes via SHA-512.",
    stack: ["Python", "Multiprocessing", "SHA-512", "NumPy"],
    year: 2025,
    status: "Live",
    url: "",
    image: "/projects/aes.png"
  },
  {
    id: 3,
    name: "Emotion-Aware Movie Recommender",
    description: "NLP-powered recommendation engine using SBERT embeddings. Achieved 0.71 Macro F1 score with hybrid emotional state interpretation.",
    stack: ["Python", "SBERT", "spaCy", "ML"],
    year: 2025,
    status: "Live",
    url: "https://github.com/puneet-chandna",
    image: "/projects/movie.png"
  },
  {
    id: 4,
    name: "Crop Stress Detection Model",
    description: "Hybrid ML pipeline combining LSTM with attention and XGBoost for binary plant-health classification using 30-day sensor data.",
    stack: ["Python", "LSTM", "XGBoost", "Scikit-learn"],
    year: 2025,
    status: "Live",
    url: "",
    image: "/projects/crop.png"
  },
  {
    id: 5,
    name: "BookMyEvent",
    description: "Event management app processing ~10,000 ticket authentications per event with real-time QR code scanning via ESP32.",
    stack: ["Node.js", "Express", "Flutter", "MongoDB"],
    year: 2024,
    status: "Live",
    url: "https://github.com/puneet-chandna",
    image: "/projects/event.png"
  },
  {
    id: 6,
    name: "Water Brakes",
    description: "Streamlit app helping farmers optimize water management by analyzing contour maps for swale and trench placement.",
    stack: ["Streamlit", "Python", "Plotly"],
    year: 2024,
    status: "Live",
    url: "https://github.com/puneet-chandna",
    image: "/projects/water.png"
  }
];

export const experience = [
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

export const socialLinks = {
  github: "https://github.com/puneet-chandna",
  linkedin: "https://www.linkedin.com/in/puneet-chandna",
  email: "puneetchandna@zohomail.in",
  website: "https://puneetchandna.com"
};
