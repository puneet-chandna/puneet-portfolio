import assert from 'node:assert/strict'
import { readFile, stat } from 'node:fs/promises'
import { test } from 'node:test'
import { runInNewContext } from 'node:vm'
import { createServer } from 'vite'
import { createElement } from 'react'
import { renderToStaticMarkup } from 'react-dom/server'

const readSource = (path) => readFile(new URL(`../${path}`, import.meta.url), 'utf8')

test('portfolio data loads with valid skill icons', async () => {
  const { skills } = await import('../src/data/portfolio.js')

  for (const skill of skills) {
    assert.equal(typeof skill.icon, 'function', skill.name)
  }
})

test('navbar exposes an accessible mobile menu', async () => {
  const navbar = await readSource('src/components/ui/Navbar.jsx')

  assert.match(navbar, /className="menu-toggle"/)
  assert.match(navbar, /aria-expanded=\{mobileMenuOpen\}/)
  assert.match(navbar, /className="mobile-menu"/)
})

test('terminal portfolio command is consistent between display and copy action', async () => {
  const contact = await readSource('src/components/sections/Contact.jsx')

  assert.match(contact, /navigator\.clipboard\.writeText\('ssh puneet\.space'\)/)
  assert.doesNotMatch(contact, /ssh puneet\.sh/)
})

test('page provides a skip link to main content', async () => {
  const app = await readSource('src/App.jsx')
  const styles = await readSource('src/styles/index.css')

  assert.match(app, /href="#main-content"/)
  assert.match(app, /<main[^>]+id="main-content"/)
  assert.match(styles, /\.skip-link/)
})

test('interactive elements have visible keyboard focus styles', async () => {
  const styles = await readSource('src/styles/index.css')

  assert.match(styles, /:focus-visible/)
  assert.match(styles, /outline:/)
  assert.match(styles, /outline-offset:/)
})

test('icon-only controls expose accessible names', async () => {
  const navbar = await readSource('src/components/ui/Navbar.jsx')
  const hero = await readSource('src/components/sections/Hero.jsx')
  const contact = await readSource('src/components/sections/Contact.jsx')

  assert.match(navbar, /aria-label="Close photo preview"/)
  assert.match(hero, /aria-label="GitHub"/)
  assert.match(hero, /aria-label="LinkedIn"/)
  assert.match(hero, /aria-label="Email"/)
  assert.match(contact, /aria-label=\{copied \? 'SSH command copied' : 'Copy SSH command'\}/)
})

test('curated projects have current repository links and existing images', async () => {
  const { projects } = await import('../src/data/portfolio.js')

  assert.equal(projects.find(project => project.name === 'Lattora').url, 'https://github.com/puneet-chandna/Lattora')
  assert.equal(projects.find(project => project.name === 'Requests Native').url, 'https://github.com/puneet-chandna/requests-native')
  assert.equal(projects.some(project => project.id === 2 || project.id === 6), false)

  for (const project of projects) {
    await stat(new URL(`../public${project.image}`, import.meta.url))
  }
})

test('project thumbnails retain their original square frames', async () => {
  const styles = await readSource('src/styles/index.css')

  assert.match(styles, /\.project-card \.project-image\s*\{[^}]*aspect-ratio:\s*1\s*\/\s*1/s)
  assert.doesNotMatch(styles, /\.project-card \.project-image\s*\{[^}]*height:\s*200px/s)
})

test('body and secondary text meet AA contrast on both page surfaces', async () => {
  const styles = await readSource('src/styles/index.css')
  const luminance = token => {
    const hex = styles.match(new RegExp(`${token}:\\s*#([a-f0-9]{6})`, 'i'))[1]
    const channels = hex.match(/../g).map(value => {
      const channel = parseInt(value, 16) / 255
      return channel <= 0.04045 ? channel / 12.92 : ((channel + 0.055) / 1.055) ** 2.4
    })
    return channels[0] * 0.2126 + channels[1] * 0.7152 + channels[2] * 0.0722
  }
  for (const background of ['--bg-dark', '--bg-card']) {
    for (const foreground of ['--text', '--text-dim']) {
      assert.ok((luminance(foreground) + 0.05) / (luminance(background) + 0.05) >= 4.5, `${foreground} on ${background}`)
    }
  }
})

test('the actual page renders accessible portfolio content without requiring WebGL', async () => {
  const server = await createServer({
    server: { middlewareMode: true, hmr: false, ws: false, watch: null },
    optimizeDeps: { noDiscovery: true, include: [] },
    cacheDir: '/tmp/puneet-portfolio-render-check',
    appType: 'custom'
  })
  try {
    const { default: App } = await server.ssrLoadModule('/src/App.jsx')
    const html = renderToStaticMarkup(createElement(App))
    assert.match(html, /PUNEET CHANDNA/)
    assert.match(html, /A software engineer with a passion for technology/)
    assert.doesNotMatch(html, /Animate background|From production workflows|Pause motion|Resume motion/)
    assert.equal((html.match(/<article/g) || []).length, 7)
    assert.match(html, /<dialog[^>]*aria-label="Photo of Puneet Chandna"/)
    assert.match(html, /autoComplete="email"/i)
    assert.match(html, /role="status"/)
    assert.doesNotMatch(html, /<canvas|app-shell/)
    const about = html.match(/<section[^>]*id="about"[\s\S]*?<\/section>/)[0]
    assert.match(about, /Product Developer at Hyr\.works/)
    assert.match(about, /Hyr’s interview assessment product/)
    assert.match(html, /multilingual transcription and translation/)
    assert.match(html, /transcript excerpts supporting each score/)
    const page = await readSource('index.html')
    assert.match(page, /class="app-shell" aria-hidden="true"/)
    const schema = JSON.parse(page.match(/<script type="application\/ld\+json">([\s\S]*?)<\/script>/)[1])
    assert.equal(schema.worksFor['@type'], 'Organization')
    assert.equal(schema.alumniOf.name, 'VIT Chennai')
    const descriptions = [...page.matchAll(/<meta\s+(?:name|property)="(?:description|og:description|twitter:description)"\s+content="([^"]+)"/g)]
    assert.equal(descriptions.length, 3)
    for (const [, description] of descriptions) {
      assert.equal(description, descriptions[0][1])
      assert.match(description, /applied AI, cloud simulation, quantitative finance/)
      assert.doesNotMatch(description, /interview assessment/i)
    }
    const { default: Scene } = await server.ssrLoadModule('/src/components/canvas/Scene.jsx')
    for (const staticBackground of [true, false]) {
      const canvas = Scene({ staticBackground }).props.children
      assert.equal(canvas.props.frameloop, staticBackground ? 'demand' : 'always')
      const sceneObjects = canvas.props.children.find(child => child?.props?.fallback === null).props.children
      assert.equal(sceneObjects.length, 3, 'Keep particles, grid and orbs in the static scene')
      for (const object of sceneObjects) assert.equal(object.props.animated, !staticBackground)
      const cameraController = canvas.props.children.find(child => child?.type?.name === 'CameraController')
      assert.equal(Boolean(cameraController), !staticBackground)
    }
  } finally {
    await server.close()
  }
})

test('the original loader waits for both page load and React, then releases input and fades away', async () => {
  const page = await readSource('index.html')
  const script = page.match(/<script>\s*([\s\S]*?)<\/script>/)[1]
  for (const loadFirst of [true, false]) {
    for (const reduced of [false, true]) {
      let appMounted = false
      let loaded = false
      let removed = false
      let disconnected = false
      let onLoad, onMutation
      const timers = []
      const root = { inert: false, hasChildNodes: () => appMounted }
      const body = { classList: { add: value => { loaded = value === 'loaded' } }, style: {} }
      runInNewContext(script, {
        document: { readyState: 'loading', body,
          getElementById: id => id === 'root' ? root : { remove: () => { removed = true } } },
        window: { matchMedia: () => ({ matches: reduced }),
          addEventListener: (_, callback) => { onLoad = callback } },
        performance: { now: () => 0 },
        setTimeout: (callback, delay) => timers.push({ callback, delay }),
        MutationObserver: class {
          constructor(callback) { onMutation = callback }
          observe() {}
          disconnect() { disconnected = true }
        }
      })
      assert.equal(root.inert, true)
      if (loadFirst) onLoad()
      else { appMounted = true; onMutation() }
      assert.equal(timers.length, 0, 'Do not reveal an unmounted or still-loading page')
      if (loadFirst) { appMounted = true; onMutation() }
      else onLoad()
      assert.equal(timers.length, 1)
      assert.equal(timers[0].delay, reduced ? 0 : 1500)
      assert.equal(disconnected, true)
      timers.shift().callback()
      assert.equal(loaded, true)
      assert.equal(root.inert, false)
      assert.equal(body.style.overflow, 'auto')
      assert.equal(timers[0].delay, reduced ? 0 : 800)
      timers.shift().callback()
      assert.equal(removed, true)
    }
  }
})
