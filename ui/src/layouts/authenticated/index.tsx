import { ROUTES, LOCAL_STORAGE_KEYS } from '@/lib/constants'
import { useQuery } from '@tanstack/react-query'
import { generatePath, Outlet } from 'react-router'
import { useNavigate } from 'react-router'
import { userQueryOptions } from '@/lib/sharedqueries'
import { Loader2 } from 'lucide-react'

export function AuthenticatedLayout() {
    const navigate = useNavigate()

    const accessToken = localStorage.getItem(LOCAL_STORAGE_KEYS.ACCESS_TOKEN)
    const refreshToken = localStorage.getItem(LOCAL_STORAGE_KEYS.REFRESH_TOKEN)

    if(!accessToken) {
        void navigate(generatePath(ROUTES.login))
        return
    }

    const {isPending, isError, data} = useQuery(userQueryOptions(accessToken, refreshToken));

    if(isPending) {
        return (
            <div className="h-screen w-screen flex flex-1 items-center justify-center">
                <Loader2 className="size-10 animate-spin" />
            </div>
        )
    }

    if(isError || !data) {
        void navigate(generatePath(ROUTES.login))
        return
    }

    console.log(data)

    return <Outlet context={{ user: data }} />
}