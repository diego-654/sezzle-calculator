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
})