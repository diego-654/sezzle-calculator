import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { calculate } from '../api/calculator'
import { ApiError } from '../api/types'
import { Calculator } from './Calculator'

// Replace the API client with a mock: these tests check the UI, not the network.
vi.mock('../api/calculator')
const calculateMock = vi.mocked(calculate)

beforeEach(() => {
  calculateMock.mockReset()
})

describe('Calculator', () => {
  it('shows the result returned by the API', async () => {
    const user = userEvent.setup()
    calculateMock.mockResolvedValue(5)
    render(<Calculator />)

    await user.type(screen.getByLabelText('First number (a)'), '2')
    await user.type(screen.getByLabelText('Second number (b)'), '3')
    await user.click(screen.getByRole('button', { name: 'Calculate' }))

    expect(calculateMock).toHaveBeenCalledWith('add', 2, 3)
    expect(await screen.findByText('5')).toBeInTheDocument()
  })

  it('calculates when the user presses Enter', async () => {
    const user = userEvent.setup()
    calculateMock.mockResolvedValue(6)
    render(<Calculator />)

    await user.selectOptions(screen.getByLabelText('Operation'), 'multiply')
    await user.type(screen.getByLabelText('First number (a)'), '2')
    await user.type(screen.getByLabelText('Second number (b)'), '3{Enter}')

    expect(calculateMock).toHaveBeenCalledWith('multiply', 2, 3)
    expect(await screen.findByText('6')).toBeInTheDocument()
  })

  it('shows the API error message', async () => {
    const user = userEvent.setup()
    calculateMock.mockRejectedValue(new ApiError('DIVISION_BY_ZERO', 'division by zero'))
    render(<Calculator />)

    await user.selectOptions(screen.getByLabelText('Operation'), 'divide')
    await user.type(screen.getByLabelText('First number (a)'), '1')
    await user.type(screen.getByLabelText('Second number (b)'), '0')
    await user.click(screen.getByRole('button', { name: 'Calculate' }))

    expect(await screen.findByRole('alert')).toHaveTextContent('division by zero')
  })

  it('disables the button until both numbers are valid', async () => {
    const user = userEvent.setup()
    render(<Calculator />)
    const button = screen.getByRole('button', { name: 'Calculate' })

    expect(button).toBeDisabled()
    await user.type(screen.getByLabelText('First number (a)'), '2')
    expect(button).toBeDisabled()
    await user.type(screen.getByLabelText('Second number (b)'), '3')
    expect(button).toBeEnabled()
  })

  it('hides the second number for square root', async () => {
    const user = userEvent.setup()
    calculateMock.mockResolvedValue(3)
    render(<Calculator />)

    await user.selectOptions(screen.getByLabelText('Operation'), 'sqrt')
    expect(screen.queryByLabelText('Second number (b)')).not.toBeInTheDocument()

    await user.type(screen.getByLabelText('First number (a)'), '9')
    await user.click(screen.getByRole('button', { name: 'Calculate' }))

    expect(calculateMock).toHaveBeenCalledWith('sqrt', 9, undefined)
  })

  it('shows a loading state while waiting for the API', async () => {
    const user = userEvent.setup()
    calculateMock.mockReturnValue(new Promise(() => {})) // never resolves
    render(<Calculator />)

    await user.type(screen.getByLabelText('First number (a)'), '2')
    await user.type(screen.getByLabelText('Second number (b)'), '3')
    await user.click(screen.getByRole('button', { name: 'Calculate' }))

    expect(screen.getByRole('button', { name: 'Calculating…' })).toBeDisabled()
  })
})
