# Auth Service

This is an authentication & authorization service built from scratch in Go.

## Features

| Feature Category         | Description |
|--------------------------|-------------|
| User Registration & Login| Create new user accounts and authenticate existing users. |
| User Logout              | End user access and invalidate refresh tokens. |
| Profile Management       | Retrieve, update, or remove the authenticated user's profile. |
| Account Recovery         | Restore account access by resetting forgotten passwords via email. |
| Authorization            | Secure API access using JWT-based authentication. |



## API Endpoints

### 1. Register a New User
POST `/auth/register`

Request Body:
```json
{
    "email": "test@skyvault.com",
    "password": "password",
    "name": "skyvault",
    "gender": "male",
    "date_of_birth": "13/10/2001"
}
```

Response:
```json
{
    "access_token": "<access-token>",
    "refresh_token": "<refresh-token>"
}
```

### 2. Login
POST `/auth/login`

Request Body:
```json
{
    "email": "test@skyvault.com",
    "password": "password"
}
```

Response:
```json
{
    "access_token": "<access-token>",
    "refresh_token": "<refresh-token>"
}
```

### 3. Logout
POST `/auth/logout`

Headers:
```
Authorization: Bearer <access-token>
```

Response:
```json
{
    "message": "user logged out successfully!"
}
```

### 4. Request New Access Token
POST `/auth/refresh`

Request Body:
```json
{
    "refresh_token": "<refresh-token>"
}
```

Response:
```json
{
    "access_token": "<access-token>"
}
```

### 5. Password Reset
POST `/auth/password-reset`

Request Body:
```json
{
    "email": "user@skyvault.com"
}
```

Response:
```json
{
    "message": "email sent"
}
```

### 6. Get Current User (Auth Info)
GET `/users/whoami`

Headers:
```
Authorization: Bearer <access-token>
```

Response:
```json
{
    "email": "user@skyvault.com",
    "name": "skyvault"
}
```

### 7. Get Current User Profile
GET `/user/me`

Headers:
```
Authorization: Bearer <access-token>
```

Response:
```json
{
    "email": "user@skyvault.com",
    "name": "skyvault",
    "gender": "male",
    "date_of_birth": "2007-08-15T14:30:00Z"
}
```

### 8. Update Current User Profile
PUT `/user/me`

Headers:
```
Authorization: Bearer <access-token>
```

Response:
```json
{
    "email": "user@skyvault.com",
    "name": "skyvault",
    "gender": "male",
    "date_of_birth": "2007-08-15T14:30:00Z"
}
```

### 9. Delete Current User Profile
DELETE `/user/me`

Headers:
```
Authorization: Bearer <access-token>
```

Response:
```json
{
    "message": "user profile deleted successfully"
}
```
