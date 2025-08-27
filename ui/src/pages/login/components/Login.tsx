import { useActionState } from 'react'
import { toast } from 'sonner'

import { InputWithFeedback } from '@/components/InputWithFeedback'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { handlePromise } from '@/lib/utils'
import { userLogin } from '@/lib/user'
import { LOCAL_STORAGE_KEYS, ROUTES } from '@/lib/constants'
import { generatePath, useNavigate } from 'react-router'

type FormState =
  | {
      status: 'error'
      errors: {
        email: string
      }
    }
  | {
      status: 'success'
    }

export function LoginForm() {
  const navigate = useNavigate()

  const [state, formAction, isPending] = useActionState<FormState, FormData>(
    async (_, formData) => {
      const errors = {
        email: '',
      }

      const email = formData.get('email') as string
      const password = formData.get('password') as string

      const [error, authResponse] = await handlePromise(userLogin({
        email: email,
        password: password
      }))

      if (error || !authResponse) {
        errors.email = 'Something went wrong.'
        return { status: 'error', errors }
      }

      localStorage.setItem(LOCAL_STORAGE_KEYS.ACCESS_TOKEN, authResponse.access_token)
      localStorage.setItem(LOCAL_STORAGE_KEYS.REFRESH_TOKEN, authResponse.refresh_token)

      toast.success('Login successful')

      void navigate(generatePath(ROUTES.home))
      return { status: 'success' }
    },
    { status: 'error', errors: { email: '' } }
  )

  return (
    <form className="flex flex-col gap-9" action={formAction}>
      <div className="flex flex-col gap-2.5">
        <Label htmlFor="email">Email</Label>
        <InputWithFeedback
          name="email"
          id="email"
          placeholder="user@skyvault.com"
          type="email"
          errorMessage={state.status === 'error' ? state.errors.email : ''}
          isError={state.status === 'error' && !!state.errors.email}
        />
      </div>
      <div className="flex flex-col gap-2.5">
        <Label htmlFor="password">Password</Label>
        <InputWithFeedback
          name="password"
          id="password"
          isError={state.status === 'error' && !!state.errors.email}
          type="password"
          helperText="Password must be at least 6 characters long"
          placeholder="********"
        />
      </div>
      <Button type="submit" isLoading={isPending} disabled={isPending} className="mt-2">
        Login
      </Button>
    </form>
  )
}