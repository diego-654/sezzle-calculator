// Operation is every path the backend exposes under /api/v1.
export type Operation =
    | 'add'
    | 'subtract'
    | 'multiply'
    | 'divide'
    | 'power'
    | 'sqrt'
    | 'percent'

// SuccessResponse is the body of a successful call: {"result": 5}.
export type SuccessResponse = { result: number }

// ErrorResponse is the body of a failed call: {"error": {"code": "...", "message": "..."}}.
export type ErrorResponse = { error: { code: string; message: string } }

// ApiError carries the stable error code so the UI can react to it.
export class ApiError extends Error {
    readonly code: string

    constructor(code: string, message: string) {
        super(message)
        this.name = 'ApiError'
        this.code = code
    }
}