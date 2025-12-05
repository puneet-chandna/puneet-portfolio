# Puneet Portfolio

A dual-interface portfolio featuring both a **terminal-based SSH experience** and a **modern 3D web application**.

![Terminal + Web Portfolio](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![React](https://img.shields.io/badge/React-20232A?style=flat&logo=react&logoColor=61DAFB)
![Three.js](https://img.shields.io/badge/Three.js-000000?style=flat&logo=three.js&logoColor=white)

---

## 🖥️ Terminal Portfolio

An SSH-accessible portfolio built with Go and the [Charm](https://charm.sh/) stack. Features a retro **Tron aesthetic** with boot sequence animations, 3-pane TUI layout, and a contact form.

### Features

- 🚀 Boot sequence animation (spinner → progress bar → ACCESS GRANTED)
- 📊 3-pane responsive layout (Menu, Viewport, Inspector)
- ⌨️ Vim-style navigation (j/k, arrow keys)
- 📨 Contact form with Web3Forms integration
- 🎨 Tron theme (cyan glow, deep black, electric blue)

### Tech Stack

- **Go** - Core language
- **Bubble Tea** - TUI framework
- **Lipgloss** - Terminal styling
- **Wish** - SSH server

### Run Locally

```bash
cd terminal
source .env  # Load Web3Forms key
./portfolio

# In another terminal:
ssh localhost -p 2222
```

### Build

```bash
cd terminal
go build -o portfolio ./cmd/portfolio
```

---

## 🌐 Web Portfolio

A modern React portfolio with immersive **3D animations** using Three.js and React Three Fiber.

### Features

- ✨ 3D particle field background
- 🔲 Tron-style glowing grid floor
- 💫 Floating orbs with bloom effects
- 🎬 Framer Motion scroll animations
- 📱 Fully responsive design
- 📨 Contact form with Web3Forms

### Tech Stack

- **React** + **Vite** - Frontend framework
- **Three.js** - 3D graphics
- **React Three Fiber** - React renderer for Three.js
- **@react-three/drei** - R3F helpers
- **@react-three/postprocessing** - Bloom, chromatic aberration
- **Framer Motion** - UI animations
- **GSAP** - Complex animations

### Run Locally

```bash
cd web
npm install
npm run dev
```

### Build for Production

```bash
cd web
npm run build
```

---

## 📁 Project Structure

```
puneet-portfolio/
├── terminal/                    # SSH Terminal Portfolio
│   ├── cmd/portfolio/main.go   # SSH server entry point
│   ├── internal/
│   │   ├── tui/                # TUI components
│   │   │   ├── model.go        # App state
│   │   │   ├── update.go       # Key handling
│   │   │   ├── view.go         # Rendering
│   │   │   ├── styles.go       # Tron theme
│   │   │   └── sender.go       # Web3Forms integration
│   │   └── data/               # Content (bio, projects)
│   ├── deploy/                 # systemd service file
│   └── .env                    # WEB3FORMS_KEY
│
├── web/                        # React 3D Portfolio
│   ├── src/
│   │   ├── components/
│   │   │   ├── canvas/         # 3D scene components
│   │   │   ├── sections/       # Page sections
│   │   │   └── ui/             # UI components
│   │   ├── data/               # Portfolio content
│   │   └── styles/             # CSS
│   └── .env                    # VITE_WEB3FORMS_KEY
│
└── README.md
```

---

## 🔐 Environment Variables

Create `.env` files with your Web3Forms API key:

**terminal/.env**

```
WEB3FORMS_KEY=your-key-here
```

**web/.env**

```
VITE_WEB3FORMS_KEY=your-key-here
```

---

## 🚀 Deployment

### Terminal (VPS)

1. Build the binary: `go build -o portfolio ./cmd/portfolio`
2. Copy to VPS and set up systemd service (see `terminal/deploy/portfolio.service`)
3. Configure firewall to allow port 22 or 2222

### Web (Vercel/Netlify)

1. Connect your GitHub repo
2. Set build command: `cd web && npm run build`
3. Set output directory: `web/dist`
4. Add environment variable: `VITE_WEB3FORMS_KEY`

---

## 📬 Contact

- **Email**: puneetchandna@zohomail.in
- **GitHub**: [puneet-chandna](https://github.com/puneet-chandna)
- **LinkedIn**: [puneet-chandna2004](https://linkedin.com/in/puneet-chandna2004)
- **Website**: [puneetchandna.com](https://puneetchandna.com)

---

## 📄 License

MIT License - feel free to use this as inspiration for your own portfolio!
