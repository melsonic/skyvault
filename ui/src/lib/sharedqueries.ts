import { queryOptions } from "@tanstack/react-query";
import { TANSTACK_QUERY_KEYS } from '@/lib/constants'
import { fetchUserIdentity } from '@/lib/user'

export function userQueryOptions(accessToken: string, refreshToken: string|null) {
    return queryOptions({
        queryKey: [TANSTACK_QUERY_KEYS.USER, accessToken],
        queryFn: () => fetchUserIdentity(accessToken, refreshToken),
        staleTime: Infinity
    })
}
