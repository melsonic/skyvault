export type User = {
    name?: string;
    email?: string;
    password?: string;
    gender?: string;
    date_of_birth?: Date;
}

export type RegisterRequest = {
    email: string;
    password: string;
    name?: string;
    gender?: string;
    date_of_birth?: Date;
}

export type LoginRequest = {
    email: string;
    password: string;
}

export type AuthenticationResponse = {
    access_token: string;
    refresh_token: string;
}

export type NewAccessTokenResponse = {
    access_token: string;
}

export type MessageResponse = {
    message: string;
}

export type UserExistByEmailResponse = {
    user_exists: boolean;
}