import { useState } from 'react'
import { LoginForm } from './components/Login'
import { RegisterForm } from './components/Register'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { LOCAL_STORAGE_KEYS, ROUTES, TAB_VALUES } from '@/lib/constants'
import { userQueryOptions } from '@/lib/sharedqueries'
import { generatePath, useNavigate } from 'react-router'
import { useQuery } from '@tanstack/react-query'
import { Loader2 } from 'lucide-react'

export function LoginPage() {
  const [tab, setTab] = useState<typeof TAB_VALUES.LOGIN | typeof TAB_VALUES.REGISTER>(
    TAB_VALUES.LOGIN
  )
  const navigate = useNavigate()

  const accessToken = localStorage.getItem(LOCAL_STORAGE_KEYS.ACCESS_TOKEN)
  const refreshToken = localStorage.getItem(LOCAL_STORAGE_KEYS.REFRESH_TOKEN)
  if(accessToken) {
    const {isPending, isSuccess, data} = useQuery(userQueryOptions(accessToken, refreshToken))
    if(isPending) {
        return (
            <div className="h-screen w-screen flex flex-1 items-center justify-center">
                <Loader2 className="size-10 animate-spin" />
            </div>
        )
    }

    if(isSuccess && data) {
      void navigate(generatePath(ROUTES.home))
      return
    }
  }

  return (
    <div className="h-screen flex flex-1 items-center justify-center p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="flex flex-col items-center gap-1">
          <CardTitle className="text-primary flex items-center justify-center gap-2 text-center text-2xl">
            SkyVault
          </CardTitle>
        </CardHeader>
        <CardContent>
          <Tabs
            value={tab}
            onValueChange={(value) => setTab(value as (typeof TAB_VALUES)[keyof typeof TAB_VALUES])}
            className="w-full"
          >
            <TabsList className="grid w-full grid-cols-2">
              <TabsTrigger value={TAB_VALUES.LOGIN}>Login</TabsTrigger>
              <TabsTrigger value={TAB_VALUES.REGISTER}>Register</TabsTrigger>
            </TabsList>

            <TabsContent value={TAB_VALUES.LOGIN} className="pt-4">
              <LoginForm />
            </TabsContent>

            <TabsContent value={TAB_VALUES.REGISTER} className="pt-4">
              <RegisterForm />
            </TabsContent>
          </Tabs>
        </CardContent>
      </Card>
    </div>
  )
}