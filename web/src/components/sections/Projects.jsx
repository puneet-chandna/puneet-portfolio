import { motion } from 'framer-motion'
import { projects } from '../../data/portfolio'

export default function Projects() {
  return (
    <section className="section" id="projects">
      <motion.div
        initial={{ opacity: 0, y: 50 }}
        whileInView={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
        viewport={{ once: true }}
      >
        <h2 className="section-title">Projects</h2>
        
        <div className="grid-3">
          {projects.map((project, index) => (
            <motion.div
              key={project.id}
              className="card project-card"
              initial={{ opacity: 0, y: 30 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.4, delay: index * 0.1 }}
              viewport={{ once: true }}
              whileHover={{ y: -10 }}
            >
              <div className="project-image">
                <span>{project.icon}</span>
              </div>
              <div style={{ padding: '20px' }}>
                <h3>{project.name}</h3>
                <p>{project.description}</p>
                <div className="project-tags">
                  {project.stack.map((tech, i) => (
                    <span key={i} className="tag">{tech}</span>
                  ))}
                </div>
                {project.url && (
                  <a
                    href={project.url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="btn btn-outline"
                    style={{ marginTop: '1rem', display: 'inline-block' }}
                  >
                    View Project →
                  </a>
                )}
              </div>
            </motion.div>
          ))}
        </div>
      </motion.div>
    </section>
  )
}
