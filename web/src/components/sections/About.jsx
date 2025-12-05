import { motion } from 'framer-motion'

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
              I'm a <span className="highlight">Computer Science student</span> at VIT Chennai, 
              passionate about building software that's both performant and elegant.
            </p>
            <p>
              With experience at <span className="highlight">IIT Bombay</span>, 
              <span className="highlight"> Daira Edtech</span>, and 
              <span className="highlight"> CeAT VIT</span>, I've worked on everything from 
              cloud optimization algorithms to full-stack web platforms.
            </p>
            <p>
              When I'm not coding, I'm leading blockchain workshops as 
              <span className="highlight"> Outreach Lead at DAO Community</span> or 
              mentoring students in the Hackclub.
            </p>
          </div>
          
          <motion.div 
            className="card glow-box"
            whileHover={{ scale: 1.02 }}
            transition={{ type: "spring", stiffness: 300 }}
          >
            <h3 style={{ color: 'var(--primary)', marginBottom: '1rem' }}>Quick Facts</h3>
            <ul style={{ listStyle: 'none', lineHeight: 2 }}>
              <li>🎓 B.Tech CSE @ VIT Chennai (2026)</li>
              <li>🔬 Research: VM Placement Optimization</li>
              <li>🚀 50% community growth at DAO</li>
              <li>👨‍🏫 Mentored 25+ blockchain devs</li>
              <li>🏆 5 major hackathons organized</li>
            </ul>
          </motion.div>
        </div>
      </motion.div>
    </section>
  )
}
