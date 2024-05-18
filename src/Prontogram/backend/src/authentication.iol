type UserAuthCredentials: void {
    .userId: string
    .password: string
}

type AuthenticatedUser: void {
    .userId: string
    .sid: string
}

interface Authentication
{
    RequestResponse:
        auth_login(UserAuthCredentials)(AuthenticatedUser) throws UserNotFound(string) UserUnauthorized(string),
        auth_logout(AuthenticatedUser)(void) throws UserUnauthorized(string)
}