export const ROUTES = {
  login: '/',
  home: '/home',
} as const

export const TAB_VALUES = {
  LOGIN: 'login',
  REGISTER: 'register',
} as const

export const LOCAL_STORAGE_KEYS = {
  ACCESS_TOKEN: 'access-token',
  REFRESH_TOKEN: 'refresh-token'
} as const

export const TANSTACK_QUERY_KEYS = {
  USER: 'user',
  NEW_TOKEN: 'new_token'
}

const AUTH_SERVICE_BASE_URL = "http://localhost:8003"

export const API_ENDPOINTS = {
  AUTH: {
    REGISTER: `${AUTH_SERVICE_BASE_URL}/auth/register`,
    LOGIN: `${AUTH_SERVICE_BASE_URL}/auth/login`,
    LOGOUT: `${AUTH_SERVICE_BASE_URL}/auth/logout`,
    USER_BY_EMAIL: `${AUTH_SERVICE_BASE_URL}/users/exists/:email`,
    NEW_ACCESS_TOKEN: `${AUTH_SERVICE_BASE_URL}/auth/refresh`,
    IDENTIFY_USER: `${AUTH_SERVICE_BASE_URL}/users/whoami`,
    GET_USER: `${AUTH_SERVICE_BASE_URL}/user/me`,
    UPDATE_USER: `${AUTH_SERVICE_BASE_URL}/user/me`,
    DELETE_USER: `${AUTH_SERVICE_BASE_URL}/user/me`,
    RESET_PASSWORD: `${AUTH_SERVICE_BASE_URL}/auth/password-reset`,
  }
} as const
