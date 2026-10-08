import { cleanup } from '@testing-library/react'
import { afterEach } from 'vitest'
// Adds DOM matchers such as toBeInTheDocument() and toBeDisabled().
import '@testing-library/jest-dom/vitest'

// Unmount rendered components after each test so tests do not share state.
afterEach(() => {
  cleanup()
})
