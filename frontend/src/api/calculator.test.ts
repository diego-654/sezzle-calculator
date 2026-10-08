import { afterEach, describe, expect, it, vi } from 'vitest'
import { calculate } from './calculator'
import { ApiError } from './types'

// mockFetch replaces the global fetch with one that returns the given response.
function mockFetch(status: number, body: unknown) {
  const fetchMock = vi.fn().mockResolvedValue(
    new Response(JSON.stringify(body), {
      status,
      headers: { 'Content-Type': 'application/json' },
    }),
  )
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('calculate', () => {
  it('posts both operands and returns the result', async () => {
    const fetchMock = mockFetch(200, { result: 5 })

    await expect(calculate('add', 2, 3)).resolves.toBe(5)
    expect(fetchMock).toHaveBeenCalledWith('/api/v1/add', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ a: 2, b: 3 }),
    })
  })

  it('omits b for unary operations', async () => {
    const fetchMock = mockFetch(200, { result: 3 })

    await calculate('sqrt', 9)
    expect(fetchMock.mock.calls[0][1].body).toBe('{"a":9}')
  })

  it('throws the backend error code and message', async () => {
    mockFetch(422, { error: { code: 'DIVISION_BY_ZERO', message: 'division by zero' } })

    const promise = calculate('divide', 1, 0)
    await expect(promise).rejects.toBeInstanceOf(ApiError)
    await expect(promise).rejects.toMatchObject({
      code: 'DIVISION_BY_ZERO',
      message: 'division by zero',
    })
  })

  it('throws NETWORK_ERROR when the server cannot be reached', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new TypeError('Failed to fetch')))

    await expect(calculate('add', 1, 2)).rejects.toMatchObject({ code: 'NETWORK_ERROR' })
  })

  it('throws UNEXPECTED_RESPONSE when an error body is not JSON', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response('<html>Bad Gateway</html>', { status: 502 })))

    await expect(calculate('add', 1, 2)).rejects.toMatchObject({ code: 'UNEXPECTED_RESPONSE' })
  })

  it('throws UNEXPECTED_RESPONSE when a success body has no result', async () => {
    mockFetch(200, { value: 5 })

    await expect(calculate('add', 2, 3)).rejects.toMatchObject({ code: 'UNEXPECTED_RESPONSE' })
  })
})
