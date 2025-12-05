import { motion } from 'framer-motion'
import { experience } from '../../data/portfolio'

export default function Experience() {
  return (
    <section className="section" id="experience">
      <motion.div
        initial={{ opacity: 0, y: 50 }}
        whileInView={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
        viewport={{ once: true }}
      >
        <h2 className="section-title">Experience</h2>
        
        <div className="timeline">
          {experience.map((exp, index) => (
            <motion.div
              key={exp.id}
              className="timeline-item"
              initial={{ opacity: 0, x: -30 }}
              whileInView={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.5, delay: index * 0.1 }}
              viewport={{ once: true }}
            >
              <div className="timeline-date">{exp.date}</div>
              <h3 className="timeline-title">{exp.title}</h3>
              <div className="timeline-company">{exp.company}</div>
              <p className="timeline-desc">{exp.description}</p>
            </motion.div>
          ))}
        </div>
      </motion.div>
    </section>
  )
}
