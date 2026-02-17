// eslint-disable-next-line no-unused-vars
import { motion } from 'framer-motion'
import { featuredCert, courseCerts, credlyProfile } from '../../data/portfolio'

export default function Certifications() {
  // Duplicate for seamless infinite scroll
  const marqueeItems = [...courseCerts, ...courseCerts]

  return (
    <section className="section" id="certifications">
      <motion.div
        initial={{ opacity: 0, y: 50 }}
        whileInView={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
        viewport={{ once: true }}
      >
        <h2 className="section-title">Certifications</h2>

        {/* ── Spotlight: Featured Certification ── */}
        <motion.a
          href={featuredCert.url}
          target="_blank"
          rel="noopener noreferrer"
          className="cert-spotlight"
          initial={{ opacity: 0, scale: 0.95 }}
          whileInView={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.7, delay: 0.1 }}
          viewport={{ once: true }}
        >
          <div className="cert-spotlight-glow" />
          <div className="cert-spotlight-inner">
            <div className="cert-badge-wrap">
              <img
                src={featuredCert.badge}
                alt={featuredCert.name}
                className="cert-badge-img"
              />
            </div>
            <div className="cert-spotlight-info">
              <span className="cert-type-tag">{featuredCert.type}</span>
              <h3 className="cert-spotlight-title">{featuredCert.name}</h3>
              <p className="cert-spotlight-issuer">{featuredCert.issuer}</p>
            </div>
          </div>
        </motion.a>

        {/* ── Marquee: Course Completions ── */}
        <motion.div
          className="cert-marquee-section"
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.3 }}
          viewport={{ once: true }}
        >
          <p className="cert-marquee-label">Courses & Skill Badges</p>
          <div className="cert-marquee-container">
            <div className="cert-marquee-track">
              {marqueeItems.map((cert, i) => (
                <div
                  key={`${cert.name}-${i}`}
                  className="cert-pill"
                  style={{
                    '--pill-color': cert.color
                  }}
                >
                  <span className="cert-pill-dot" />
                  <span className="cert-pill-name">{cert.name}</span>
                  <span className="cert-pill-issuer">{cert.issuer}</span>
                </div>
              ))}
            </div>
          </div>
        </motion.div>

        {/* ── Credly CTA ── */}
        <motion.div
          className="cert-cta"
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.4 }}
          viewport={{ once: true }}
        >
          <a
            href={credlyProfile}
            target="_blank"
            rel="noopener noreferrer"
            className="btn btn-outline"
          >
            View all credentials on Credly →
          </a>
        </motion.div>
      </motion.div>
    </section>
  )
}
