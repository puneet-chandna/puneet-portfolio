import { useState } from 'react'
import { motion } from 'framer-motion'

const WEB3FORMS_KEY = import.meta.env.VITE_WEB3FORMS_KEY

export default function Contact() {
  const [formData, setFormData] = useState({ name: '', email: '', message: '' })
  const [status, setStatus] = useState('')
  const [copied, setCopied] = useState(false)

  const handleSubmit = async (e) => {
    e.preventDefault()
    setStatus('sending')

    try {
      const response = await fetch('https://api.web3forms.com/submit', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          access_key: WEB3FORMS_KEY,
          ...formData,
          subject: 'New message from Portfolio Website'
        })
      })

      if (response.ok) {
        setStatus('success')
        setFormData({ name: '', email: '', message: '' })
      } else {
        setStatus('error')
      }
    } catch {
      setStatus('error')
    }
  }

  const copySSH = () => {
    navigator.clipboard.writeText('ssh puneet.sh')
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <section className="section" id="contact">
      <motion.div
        initial={{ opacity: 0, y: 50 }}
        whileInView={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
        viewport={{ once: true }}
      >
        <h2 className="section-title">Get in Touch</h2>
        
        <div className="contact-content">
          <form className="contact-form" onSubmit={handleSubmit}>
            <div className="form-group">
              <label htmlFor="name">Name</label>
              <input
                type="text"
                id="name"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                required
                placeholder="Your name"
              />
            </div>
            
            <div className="form-group">
              <label htmlFor="email">Email</label>
              <input
                type="email"
                id="email"
                value={formData.email}
                onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                required
                placeholder="your@email.com"
              />
            </div>
            
            <div className="form-group">
              <label htmlFor="message">Message</label>
              <textarea
                id="message"
                value={formData.message}
                onChange={(e) => setFormData({ ...formData, message: e.target.value })}
                required
                placeholder="Your message..."
              />
            </div>
            
            <button
              type="submit"
              className="btn btn-primary"
              disabled={status === 'sending'}
              style={{ width: '100%' }}
            >
              {status === 'sending' ? 'Sending...' : 'Send Message'}
            </button>
            
            {status === 'success' && (
              <p style={{ color: 'var(--success)', textAlign: 'center', marginTop: '1rem' }}>
                ✓ Message sent successfully!
              </p>
            )}
            {status === 'error' && (
              <p style={{ color: 'var(--accent)', textAlign: 'center', marginTop: '1rem' }}>
                Something went wrong. Please try again.
              </p>
            )}
          </form>

          {/* SSH Banner */}
          <motion.div
            className="ssh-banner"
            initial={{ opacity: 0 }}
            whileInView={{ opacity: 1 }}
            transition={{ duration: 0.6, delay: 0.3 }}
            viewport={{ once: true }}
          >
            <h3 style={{ color: 'var(--primary)', marginBottom: '0.5rem' }}>
              Try my Terminal Portfolio
            </h3>
            <p style={{ color: 'var(--text-dim)' }}>
              For the full hacker experience, connect via SSH:
            </p>
            <div className="ssh-command">
              <span>$ ssh puneet.sh</span>
              <button onClick={copySSH} title="Copy to clipboard">
                {copied ? '✓' : '📋'}
              </button>
            </div>
          </motion.div>
        </div>
      </motion.div>
    </section>
  )
}
