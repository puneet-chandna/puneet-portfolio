package data

// GetBio returns a simplified about me section
func GetBio() string {
	return `> STATUS: ONLINE
> ROLE: PRODUCT DEVELOPER @ HYR.WORKS
> EDUCATION: B.Tech CSE, VIT Chennai
> GRADUATION: 2026 (COMPLETED)

TECH ARSENAL:
  Python     JavaScript  Go
  C++        Node.js     React
  Next.js    MongoDB     PostgreSQL
  AWS/GCP    Docker      Linux/Git

CORE COMPETENCIES:
  • Full Stack Development & RESTful APIs
  • Data Structures & Algorithms
  • System Design & Optimization
  • Cryptography (AES/SHA-512)
  • Parallel Processing & Cloud Computing

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
  Developed and stabilized the Next.js V2 recruiter
  platform, fixing recurring client-reported bugs
  across frontend and backend workflows. Built the
  Hyr Live Chrome extension and integrated Hyr Agent APIs.

  Hardened tenant isolation across Next.js and Java
  with Supabase RLS, JWT authentication, and RBAC;
  resolved auth and workflow regressions through
  integration testing.

  Rebuilt the FFmpeg recording pipeline with dual
  recorders, 2-minute chunks, browser buffering, and
  verified uploads to preserve partial interviews
  during network failures. Added Grafana/Sentry
  monitoring and Telegram alerts across DigitalOcean.

  Independently researched, designed, and built
  hyr.works in Next.js, exploring 16+ prototypes
  and implementing responsive pages, technical SEO,
  and generative engine optimization (GEO).

• Backend Engineering Intern @ Apoliums Infotech India Pvt. Ltd.
  Dec 2025-Jan 2026
  Migrated legacy Node.js services to Go (Gin),
  optimized MySQL schemas, and deployed REST APIs on GCP.

• Research Intern @ CeAT, VIT Chennai
  May-Jul 2025
  CloudSim Plus eval for HO algorithm,
  scalable modular simulations

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
