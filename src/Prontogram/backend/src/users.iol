
type UserSignUpRequest: void {
    .userId: string
    .password: string
    .display_name: string
}

interface Users
{
    RequestResponse:
        user_signup(UserSignUpRequest)(void) throws UserAlreadyExists(string)
}