import type { User } from '@/lib/types'
import { useOutletContext } from 'react-router'

export type AuthenticatedLayoutContext = {
  currentUser: User
}

export const useAuthLayoutContext = () => {
  return useOutletContext<AuthenticatedLayoutContext>()
}