import { motion } from 'framer-motion'
import { createElement } from 'react'
import {
  FaChalkboardTeacher,
  FaChartLine,
  FaFlask,
  FaGraduationCap,
  FaTrophy,
} from 'react-icons/fa'

const quickFacts = [
  { icon: FaGraduationCap, text: 'B.Tech CSE @ VIT Chennai (2026)' },
  { icon: FaFlask, text: 'Research: VM Placement Optimization' },
  { icon: FaChartLine, text: '50% community growth at DAO' },
  { icon: FaChalkboardTeacher, text: 'Mentored 25+ blockchain devs' },
  { icon: FaTrophy, text: '5 major hackathons organized' },
]

export default function About() {
  return (
    <section className="section" id="about">
      <motion.div
        initial={{ opacity: 0, y: 50 }}
        whileInView={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
        viewport={{ once: true }}
      >
        <h2 className="section-title">About Me</h2>
        
        <div className="about-content">
          <div className="about-text">
            <p>
              I am <span className="highlight">Puneet Chandna</span>, a software engineer and problem solver based in India, currently a final year student at <span className="highlight">VIT Chennai</span>. 
            </p>
            <p>
              My professional experience includes roles as a Research Intern at <span className="highlight">CeAT VIT</span> and Full Stack Developer Intern at <span className="highlight">Daira EdTech</span>, where I honed my skills in solving real-world engineering challenges.
            </p>
            <p>
              Beyond engineering, I have been actively involved in student tech communities. As the <span className="highlight">Blockchain Lead at HackClub VIT</span> and a core member of the <span className="highlight">DAO Community</span> (Outreach & Web Lead), I’ve organized hackathons, guided student teams, and helped foster a culture of hands-on learning.
            </p>
            <p>
              Outside of technology, I enjoy science-fiction films, minimalist design, and long conversations around productivity and personal growth.I also enjoy the challenge of maintaining a balance between my culinary explorations and fitness goals.
            </p>
            <p>I’m always open to exploring new ideas or collaborating on meaningful projects.</p>
            <p>
              Feel free to reach out at <a href="mailto:puneetchandna@zohomail.in" style={{ color: 'var(--primary)', textDecoration: 'none' }}>puneetchandna@zohomail.in</a>
            </p>
          </div>
          
          <motion.div 
            className="card glow-box"
            whileHover={{ scale: 1.02 }}
            transition={{ type: "spring", stiffness: 300 }}
          >
            <h3 style={{ color: 'var(--primary)', marginBottom: '1rem' }}>Quick Facts</h3>
            <ul className="quick-facts-list">
              {quickFacts.map(({ icon, text }) => (
                <li className="quick-fact" key={text}>
                  <span className="quick-fact-icon" aria-hidden="true">
                    {createElement(icon)}
                  </span>
                  <span>{text}</span>
                </li>
              ))}
            </ul>
          </motion.div>
        </div>
      </motion.div>
    </section>
  )
}
