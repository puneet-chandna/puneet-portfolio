import Scene from './components/canvas/Scene'
import Navbar from './components/ui/Navbar'
import Hero from './components/sections/Hero'
import About from './components/sections/About'
import Experience from './components/sections/Experience'
import Projects from './components/sections/Projects'
import Skills from './components/sections/Skills'
import Certifications from './components/sections/Certifications'
import Contact from './components/sections/Contact'
import './styles/index.css'

function App() {
  return (
    <>
      <Scene />
      <Navbar />
      <main className="content">
        <Hero />
        <About />
        <Experience />
        <Projects />
        <Skills />
        <Certifications />
        <Contact />
      </main>
      <footer style={{
        textAlign: 'center',
        padding: '2rem',
        color: 'var(--text)',
        borderTop: '1px solid var(--dim)'
      }}>
        <p>© 2025 Puneet Chandna. All rights reserved.</p>
      </footer>
    </>
  )
}

export default App
