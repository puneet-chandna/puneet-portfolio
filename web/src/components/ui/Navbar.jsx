import { useState, useEffect } from 'react'
import { motion } from 'framer-motion'

export default function Navbar() {
  const [scrolled, setScrolled] = useState(false)

  useEffect(() => {
    const handleScroll = () => {
      setScrolled(window.scrollY > 50)
    }
    window.addEventListener('scroll', handleScroll)
    return () => window.removeEventListener('scroll', handleScroll)
  }, [])

  return (
    <motion.nav
      className="navbar"
      initial={{ y: -100 }}
      animate={{ y: 0 }}
      transition={{ duration: 0.5 }}
      style={{
        background: scrolled ? 'rgba(10, 10, 15, 0.95)' : 'rgba(10, 10, 15, 0.8)'
      }}
    >
      <a href="#home" className="logo">
        <img 
          src="/puneet.png" 
          alt="Puneet Chandna" 
          style={{ 
            width: '32px', 
            height: '32px', 
            borderRadius: '50%',
            marginRight: '10px',
            objectFit: 'cover'
          }} 
        />
        <span>Puneet Chandna</span>
      </a>
      
      <ul className="nav-links">
        <li><a href="#about">About</a></li>
        <li><a href="#experience">Experience</a></li>
        <li><a href="#projects">Projects</a></li>
        <li><a href="#skills">Skills</a></li>
        <li><a href="#contact">Contact</a></li>
      </ul>
    </motion.nav>
  )
}
