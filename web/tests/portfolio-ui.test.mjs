import assert from 'node:assert/strict'
import { readFile } from 'node:fs/promises'
import { test } from 'node:test'

const readSource = (path) => readFile(new URL(`../${path}`, import.meta.url), 'utf8')

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
