import { motion } from 'framer-motion'
import { skills } from '../../data/portfolio'

export default function Skills() {
  return (
    <section className="section" id="skills">
      <motion.div
        initial={{ opacity: 0, y: 50 }}
        whileInView={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
        viewport={{ once: true }}
      >
        <h2 className="section-title">Skills</h2>
        
        <div className="skills-grid">
          {skills.map((skill, index) => (
            <motion.div
              key={skill.name}
              className="skill-item"
              initial={{ opacity: 0, scale: 0.8 }}
              whileInView={{ opacity: 1, scale: 1 }}
              transition={{ duration: 0.3, delay: index * 0.05 }}
              viewport={{ once: true }}
              whileHover={{ scale: 1.1, borderColor: skill.color }}
            >
              <skill.icon className="icon" style={{ color: skill.color }} />
              <span style={{ color: skill.color }}>{skill.name}</span>
            </motion.div>
          ))}
        </div>
      </motion.div>
    </section>
  )
}
