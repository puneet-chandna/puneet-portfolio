package data

// GetBio returns a simplified about me section
func GetBio() string {
	return `> STATUS: ONLINE
> ROLE: PRODUCT DEVELOPER @ HYR.WORKS
> EDUCATION: B.Tech CSE, VIT Chennai
> GRADUATION: 2026 (COMPLETED)

Software engineer building web products, backend systems and research tools. Explore applied AI, cloud simulation, quantitative finance and native software.

CURRENT WORK:
  Product Developer at Hyr.works (Zofa AI Solutions Pvt. Ltd.) since March 2026. Helped build configurable, multi-stage AI interviews and an AI agent that finds candidates and adds them to the recruitment pipeline.

  Independently built and deployed Hyr's interview assessment product to evaluate candidate answers and interviewer performance from recorded interviews, with AI scores supported by transcript excerpts. Owned the backend, database, dashboard, interview integrations and production deployment. Also independently researched, designed and built hyr.works.

TECH ARSENAL:
  Python     JavaScript  Go
  C++        Node.js     React
  Next.js    MongoDB     PostgreSQL
  AWS/GCP    Docker      Linux/Git
  Applied AI            AI Agents
  Machine Learning      PyTorch

CORE COMPETENCIES:
  • Full Stack Development & RESTful APIs
  • Data Structures & Algorithms
  • System Design & Optimization
  • Cryptography (AES/SHA-512)
  • Parallel Processing & Cloud Computing
  • Applied AI & AI Agents

CERTIFICATION:
  Google Cloud Professional Cloud Architect

> "Building software that's both performant
   and elegant."
`
}

// GetExperience returns work experience/mission log
func GetExperience() string {
	return `MISSION LOG:

• Product Developer @ Hyr.works (Zofa AI Solutions Pvt. Ltd.)
  Mar 2026-Present
  Helped build applied AI systems for Hyr's recruiting platform, including configurable, multi-stage AI interviews and an AI agent that searches for candidates and adds them to the recruitment pipeline.

  Independently built and deployed Hyr's interview assessment product to help reviewers evaluate candidate answers and interviewer performance from recorded interviews. Applied assessment research to AI scoring, with transcript excerpts supporting each score. Owned the backend, database, dashboard, scheduling and recording integrations, multilingual transcription and translation, and production deployment.

  Developed and stabilized Hyr's V2 recruiter platform in Next.js, resolving recurring client-reported bugs across frontend and backend workflows. Built the Hyr Live Chrome extension and integrated Hyr Agent APIs.

  Hardened tenant isolation across Next.js and Java services with Supabase row-level security, JWT authentication and role-based access controls; resolved authentication and workflow regressions through integration testing.

  Rebuilt the FFmpeg recording pipeline with dual recorders, 2-minute chunks, browser buffering and verified uploads to preserve partial interviews during network failures. Added Grafana/Sentry monitoring and Telegram alerts across DigitalOcean production services.

  Independently researched, designed and built hyr.works in Next.js, exploring 16+ prototypes and implementing responsive pages, technical SEO and generative engine optimization (GEO).

• Backend Engineering Intern @ Apoliums Infotech India Pvt. Ltd.
  Dec 2025-Jan 2026
  Migrated legacy Node.js services to Go (Gin),
  optimized MySQL schemas, and deployed REST APIs on GCP.

• Research Intern @ CeAT, VIT Chennai
  May-Jul 2025
  Developed a CloudSim Plus framework for VM placement optimization using Hippopotamus Optimization, with scalable modular simulations. Authored a research paper on simulation results.

• Full Stack Developer @ Daira Edtech
  Dec 2024 - Feb 2025
  Built AMS & progress modules, RBAC,
  refactored backend for reliability; data models
  improved data retrieval efficiency by 30%

• Web Dev Intern @ IIT Bombay
  Sept-Oct 2024
  Cut payload size by 90%, added Joi
  validation for stronger APIs
`
}
