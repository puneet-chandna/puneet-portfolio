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
