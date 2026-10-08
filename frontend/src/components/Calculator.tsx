import { useState } from 'react'
import type { FormEvent } from 'react'
import { calculate } from '../api/calculator'
import { ApiError, type Operation } from '../api/types'

const OPERATIONS: { value: Operation; label: string }[] = [
  { value: 'add', label: 'Add (a + b)' },
  { value: 'subtract', label: 'Subtract (a − b)' },
  { value: 'multiply', label: 'Multiply (a × b)' },
  { value: 'divide', label: 'Divide (a ÷ b)' },
  { value: 'power', label: 'Power (a ^ b)' },
  { value: 'sqrt', label: 'Square root (√a)' },
  { value: 'percent', label: 'Percent (b% of a)' },
]

// parseNumber returns the input as a finite number, or null if it is not one.
function parseNumber(text: string): number | null {
  if (text.trim() === '') return null
  const n = Number(text)
  return Number.isFinite(n) ? n : null
}

export function Calculator() {
  const [a, setA] = useState('')
  const [b, setB] = useState('')
  const [operation, setOperation] = useState<Operation>('add')
  const [result, setResult] = useState<number | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  // Derived values: computed on every render instead of stored in state.
  const needsB = operation !== 'sqrt'
  const numA = parseNumber(a)
  const numB = needsB ? parseNumber(b) : null
  const isValid = numA !== null && (!needsB || numB !== null)

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (numA === null || (needsB && numB === null)) return

    setLoading(true)
    setError(null)
    setResult(null)
    try {
      setResult(await calculate(operation, numA, numB ?? undefined))
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Something went wrong. Please try again.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <form className="calculator" onSubmit={handleSubmit} noValidate>
      <label htmlFor="a">First number (a)</label>
      <input
        id="a"
        type="number"
        step="any"
        inputMode="decimal"
        value={a}
        onChange={(e) => setA(e.target.value)}
        required
      />

      <label htmlFor="operation">Operation</label>
      <select
        id="operation"
        value={operation}
        onChange={(e) => setOperation(e.target.value as Operation)}
      >
        {OPERATIONS.map((op) => (
          <option key={op.value} value={op.value}>
            {op.label}
          </option>
        ))}
      </select>

      {needsB && (
        <>
          <label htmlFor="b">Second number (b)</label>
          <input
            id="b"
            type="number"
            step="any"
            inputMode="decimal"
            value={b}
            onChange={(e) => setB(e.target.value)}
            required
          />
        </>
      )}

      <button type="submit" disabled={!isValid || loading}>
        {loading ? 'Calculating…' : 'Calculate'}
      </button>

      {!isValid && <p className="hint">Enter {needsB ? 'both numbers' : 'a number'} to calculate.</p>}

      {/* aria-live makes screen readers announce the result when it changes. */}
      <p className="result" aria-live="polite">
        {result !== null && (
          <>
            Result: <output>{result}</output>
          </>
        )}
      </p>

      {error && (
        <p className="error" role="alert">
          {error}
        </p>
      )}
    </form>
  )
}
