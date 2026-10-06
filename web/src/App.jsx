import { Component, lazy, Suspense, useSyncExternalStore } from 'react'
import { MotionConfig } from 'framer-motion'
import Navbar from './components/ui/Navbar'
import Hero from './components/sections/Hero'
import About from './components/sections/About'
import Experience from './components/sections/Experience'
import Projects from './components/sections/Projects'
import Skills from './components/sections/Skills'
import Certifications from './components/sections/Certifications'
import Contact from './components/sections/Contact'
import './styles/index.css'

const Scene = lazy(() => import('./components/canvas/Scene'))
const staticBackgroundQuery = '(prefers-reduced-motion: reduce), (max-width: 768px), (hover: none) and (pointer: coarse)'
const readBackgroundPreference = () => window.matchMedia(staticBackgroundQuery).matches
function subscribeToBackgroundPreference(onChange) {
  const query = window.matchMedia(staticBackgroundQuery)
  query.addEventListener('change', onChange)
  return () => query.removeEventListener('change', onChange)
}

class BackgroundBoundary extends Component {
  state = { failed: false }
  static getDerivedStateFromError() { return { failed: true } }
  render() { return this.state.failed ? null : this.props.children }
}

function App() {
  const staticBackground = useSyncExternalStore(subscribeToBackgroundPreference, readBackgroundPreference, () => true)
  return (
    <MotionConfig reducedMotion="user">
      <a className="skip-link" href="#main-content">
        Skip to main content
      </a>
      <BackgroundBoundary>
        <Suspense fallback={null}><Scene staticBackground={staticBackground} /></Suspense>
      </BackgroundBoundary>
      <Navbar />
      <main className="content" id="main-content" tabIndex="-1">
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
        <p>© {new Date().getFullYear()} Puneet Chandna. All rights reserved.</p>
      </footer>
    </MotionConfig>
  )
}

export default App
