/// <reference types="vitest/config" />
import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    // Forward API calls to the Go backend so the browser only talks to
    // one origin and no CORS is needed. nginx does the same in production.
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
  test: {
    // jsdom simulates a browser (document, window) inside Node.
    environment: 'jsdom',
    setupFiles: ['./src/setupTests.ts'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'html'],
      include: ['src/**/*.{ts,tsx}'],
      exclude: ['src/main.tsx', 'src/**/*.test.{ts,tsx}', 'src/setupTests.ts'],
    },
  },
})
