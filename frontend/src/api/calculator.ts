import { type ErrorResponse, type Operation, type SuccessResponse, ApiError } from './types'



// calculate calls POST /api/v1/{op} and returns the result,
// or throws an ApiError for any kind of failure.
export async function calculate(op: Operation, a: number, b?: number): Promise<number> {
  let response: Response
  try {
    response = await fetch(`/api/v1/${op}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ a, b }),
    })
  } catch {
    // fetch only rejects when the request never got a response.
    throw new ApiError('NETWORK_ERROR', 'Could not reach the server. Please try again.')
  }

  const body: unknown = await response.json().catch(() => null)

  if (!response.ok) {
    if (isErrorResponse(body)) {
      throw new ApiError(body.error.code, body.error.message)
    }
    throw new ApiError('UNEXPECTED_RESPONSE', `Unexpected server response (${response.status}).`)
  }

  if (!isSuccessResponse(body)) {
    throw new ApiError('UNEXPECTED_RESPONSE', 'The server returned an invalid result.')
  }
  return body.result
}

function isSuccessResponse(body: unknown): body is SuccessResponse {
  return (
    typeof body === 'object' && body !== null && 'result' in body && typeof body.result === 'number'
  )
}

function isErrorResponse(body: unknown): body is ErrorResponse {
  return typeof body === 'object' && body !== null && 'error' in body
}