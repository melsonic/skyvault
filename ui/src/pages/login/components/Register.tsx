import { useActionState } from 'react'
import { toast } from 'sonner'

import { InputWithFeedback } from '@/components/InputWithFeedback'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { handlePromise } from '@/lib/utils'
import { checkExistingUserByEmail, userRegister } from '@/lib/user'
import { LOCAL_STORAGE_KEYS, ROUTES } from '@/lib/constants'
import { generatePath, useNavigate } from 'react-router'

const PASSWORD_MIN_LENGTH = 8

type FormState =
  | {
      status: 'error'
      errors: {
        email: string
        password: string
      }
    }
  | {
      status: 'success'
    }

export function RegisterForm() {
  const navigate = useNavigate()

  const [state, formAction, isPending] = useActionState<FormState, FormData>(
    async (_, formData) => {
      const email = formData.get('email') as string
      const password = formData.get('password') as string
      const confirmPassword = formData.get('confirmPassword') as string

      const errors = {
        email: '',
        password: '',
      }

      if (password.length < PASSWORD_MIN_LENGTH) {
        errors.password = `Password must be at least ${PASSWORD_MIN_LENGTH} characters long`
        return { status: 'error', errors }
      }

      if (password !== confirmPassword) {
        errors.password = 'Passwords do not match'
        return { status: 'error', errors }
      }

      // TODO: need to create another API route in authcomp to handle existing user emails
      const [existingUserError, existingUser] = await handlePromise(checkExistingUserByEmail(email))

      if (!existingUserError) {
        errors.email = 'Something went wrong during registration. Please try later.'
        return { status: 'error', errors }
      }

      if (existingUser) {
        errors.email = 'Email already exists'
        return { status: 'error', errors }
      }

      // TODO: perform signUp function
      const [error, authResponse] = await handlePromise(userRegister({
        email: email,
        password: password
      }))

      if (error || !authResponse) {
        errors.email = 'Something went wrong during registration. Please try later.'
        return { status: 'error', errors }
      }

      localStorage.setItem(LOCAL_STORAGE_KEYS.ACCESS_TOKEN, authResponse.access_token)
      localStorage.setItem(LOCAL_STORAGE_KEYS.REFRESH_TOKEN, authResponse.refresh_token)

      toast.success('Registration successful')
      void navigate(generatePath(ROUTES.home))
      return { status: 'success' }
    },
    { status: 'error', errors: { email: '', password: '' } }
  )

  return (
    <form action={formAction} className="flex flex-col gap-9">
      <div className="flex flex-col gap-2.5">
        <Label htmlFor="email">Email</Label>
        <InputWithFeedback
          name="email"
          id="email"
          placeholder="user@skyvault.com"
          type="email"
          errorMessage={state.status === 'error' ? state.errors?.email : ''}
          isError={state.status === 'error' && !!state.errors?.email}
          required
        />
      </div>
      <div className="flex flex-col gap-2.5">
        <Label htmlFor="password">Password</Label>
        <InputWithFeedback
          name="password"
          id="password"
          errorMessage={state.status === 'error' ? state.errors?.password : ''}
          isError={state.status === 'error' && !!state.errors?.password}
          required
          type="password"
          helperText="Password must be at least 8 characters long"
          placeholder="********"
        />
      </div>
      <div className="flex flex-col gap-2.5">
        <Label htmlFor="confirmPassword">Confirm Password</Label>
        <InputWithFeedback
          name="confirmPassword"
          id="confirmPassword"
          required
          // just show error border if any password errors
          isError={state.status === 'error' && !!state.errors?.password}
          type="password"
          placeholder="********"
        />
      </div>
      <Button type="submit" isLoading={isPending} disabled={isPending} className="mt-2">
        Register
      </Button>
    </form>
  )
}