import { API_ENDPOINTS, LOCAL_STORAGE_KEYS } from "./constants";
import type { AuthenticationResponse, LoginRequest, MessageResponse, NewAccessTokenResponse, RegisterRequest, User, UserExistByEmailResponse } from "@/lib/types";
import { generatePath } from 'react-router'

// User tries to go to home page
// Get's unauthorized
// It requests a new access token
// In case it succeds, goes to /home
// Else goes to /login

export async function fetchUserIdentity(accessToken: string, refreshToken: string|null): Promise<User|null> {
    const response = await fetch(API_ENDPOINTS.AUTH.IDENTIFY_USER, {
        headers: {
            'Authorization': `Bearer ${accessToken}`
        }
    })
    if(response.status === 401) {
        if(refreshToken) await newAccessToken(refreshToken)
        return null
    }
    if(response.status !== 200) {
        throw new Error("Error in fetch user identity")
    }
    const user: User = await response.json()
    return user
}

export async function userRegister(user: RegisterRequest): Promise<AuthenticationResponse> {
    const response = await fetch(API_ENDPOINTS.AUTH.REGISTER, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(user)
    })
    if(response.status !== 200) {
        throw new Error("Error in User Registration")
    }
    const authTokens: AuthenticationResponse = await response.json()
    return authTokens
}

export async function userLogin(user: LoginRequest): Promise<AuthenticationResponse> {
    const response = await fetch(API_ENDPOINTS.AUTH.LOGIN, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(user)
    })
    if(response.status !== 200) {
        throw new Error("Error in User Login")
    }
    const authTokens: AuthenticationResponse = await response.json()
    return authTokens
}

export async function checkExistingUserByEmail(email: string): Promise<boolean> {
    const response = await fetch(generatePath(API_ENDPOINTS.AUTH.USER_BY_EMAIL, {email: email}))
    if(response.status !== 200) {
        throw new Error("Error in fetch user identity")
    }
    const data: UserExistByEmailResponse = await response.json()
    return data.user_exists
}


export async function userLogout(accessToken: string): Promise<MessageResponse> {
    const response = await fetch(API_ENDPOINTS.AUTH.LOGOUT, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${accessToken}`
        },
    })
    if(response.status !== 200) {
        throw new Error("Error in User Logout")
    }
    const messageResponse: MessageResponse = await response.json()
    return messageResponse
}

export async function fetchUserProfile(accessToken: string, refreshToken: string): Promise<User> {
    const response = await fetch(API_ENDPOINTS.AUTH.GET_USER, {
        headers: {
            'Authorization': `Bearer ${accessToken}`
        },
    })
    if(response.status !== 200) {
        if(response.status === 401 && refreshToken) {
            await newAccessToken(refreshToken)
        } else {
            throw new Error("Error in fetch user identity")
        }
    }
    const userProfile: User = await response.json()
    return userProfile
}

export async function updateUserProfile(accessToken: string, user: User, refreshToken: string): Promise<User> {
    const response = await fetch(API_ENDPOINTS.AUTH.UPDATE_USER, {
        method: 'PUT',
        headers: {
            'Authorization': `Bearer ${accessToken}`
        },
        body: JSON.stringify(user)
    })
    if(response.status !== 200) {
        if(response.status === 401 && refreshToken) {
            await newAccessToken(refreshToken)
        } else {
            throw new Error("Error in fetch user identity")
        }
    }
    const userProfile: User = await response.json()
    return userProfile
}

export async function deleteUserProfile(accessToken: string, user: User, refreshToken: string|null): Promise<MessageResponse> {
    const response = await fetch(API_ENDPOINTS.AUTH.DELETE_USER, {
        method: 'DELETE',
        headers: {
            'Authorization': `Bearer ${accessToken}`
        },
        body: JSON.stringify(user)
    })
    if(response.status !== 200) {
        if(response.status === 401 && refreshToken) {
            await newAccessToken(refreshToken)
        } else {
            throw new Error("Error in fetch user identity")
        }
    }
    const messageResponse: MessageResponse = await response.json()
    return messageResponse
}

export async function passWordResetRequest(user: User): Promise<MessageResponse> {
    const response = await fetch(API_ENDPOINTS.AUTH.RESET_PASSWORD, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(user)
    })
    if(response.status !== 200) {
        throw new Error("Error in requesting Password Reset")
    }
    const messageResponse: MessageResponse = await response.json()
    return messageResponse
}

export async function newAccessToken(refreshToken: string) {
    const response = await fetch(API_ENDPOINTS.AUTH.NEW_ACCESS_TOKEN, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            refresh_token: refreshToken
        })
    })

    if(response.status != 200) {
        throw new Error("Error in requesting new access token")
    }
    const accessToken: NewAccessTokenResponse = await response.json()
    localStorage.setItem(LOCAL_STORAGE_KEYS.ACCESS_TOKEN, accessToken.access_token)
}