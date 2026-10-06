import { useState, useEffect, useRef } from 'react'
import { FaBars, FaTimes } from 'react-icons/fa'

const navItems = [
  { href: '#about', label: 'About' },
  { href: '#experience', label: 'Experience' },
  { href: '#projects', label: 'Projects' },
  { href: '#skills', label: 'Skills' },
  { href: '#certifications', label: 'Certifications' },
  { href: '#contact', label: 'Contact' },
]

export default function Navbar() {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)
  const photoRef = useRef(null)
  const menuRef = useRef(null)

  useEffect(() => {
    if (!mobileMenuOpen) return
    const closeOnEscape = (event) => {
      if (event.key === 'Escape') {
        setMobileMenuOpen(false)
        menuRef.current?.focus()
      }
    }
    window.addEventListener('keydown', closeOnEscape)
    return () => window.removeEventListener('keydown', closeOnEscape)
  }, [mobileMenuOpen])

  return (
    <>
      <nav className="navbar" aria-label="Main navigation">
        <div className="nav-identity">
          <button type="button" className="portrait-button" aria-label="View photo of Puneet Chandna"
            onClick={() => photoRef.current.showModal()}>
            <img src="/puneet.webp" alt="" width="32" height="32" />
          </button>
          <a href="#home" className="logo">Puneet Chandna</a>
        </div>
        <ul className="nav-links">
          {navItems.map(item => (
            <li key={item.href}><a href={item.href}>{item.label}</a></li>
          ))}
        </ul>
        <button ref={menuRef} type="button" className="menu-toggle"
          aria-label={mobileMenuOpen ? 'Close navigation menu' : 'Open navigation menu'}
          aria-controls="mobile-menu" aria-expanded={mobileMenuOpen}
          onClick={() => setMobileMenuOpen(open => !open)}>
          {mobileMenuOpen ? <FaTimes aria-hidden="true" /> : <FaBars aria-hidden="true" />}
        </button>
      </nav>
      <nav id="mobile-menu" className="mobile-menu" aria-label="Mobile navigation" hidden={!mobileMenuOpen}>
        {navItems.map(item => (
          <a key={item.href} href={item.href} onClick={() => setMobileMenuOpen(false)}>{item.label}</a>
        ))}
      </nav>
      <dialog ref={photoRef} className="photo-dialog" aria-label="Photo of Puneet Chandna"
        onClick={event => { if (event.target === event.currentTarget) photoRef.current.close() }}>
        <button type="button" className="photo-close" autoFocus aria-label="Close photo preview"
          onClick={() => photoRef.current.close()}><span aria-hidden="true">×</span></button>
        <img src="/puneet.webp" alt="Puneet Chandna" width="640" height="640" />
      </dialog>
    </>
  )
}
