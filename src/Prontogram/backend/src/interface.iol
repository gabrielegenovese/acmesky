
type UserAuthCredentials {
    user_id: string
    password: string
}

type UserSignUpRequest {
    .display_name: string
    .credentials: UserAuthCredentials
}

type AuthenticatedUser {
    user_id: string
    sid: string
}

type SendMessageRequest: void {
    .sender: AuthenticatedUser
    .receiver_user_id: string
    .content: string
}

type Message: void {
    .id: string
    .sender_user_id: string
    .receiver_user_id: string
    .content: string
}

type MessageList: void {
    .messages[0,*]: Message
}

interface IProntogramService
{

    RequestResponse:
        auth_signup(UserSignUpRequest)(void) throws UserAlreadyExists(string),
        auth_login(UserAuthCredentials)(AuthenticatedUser) throws UserNotFound(string) UserNotAuthorized(string),
        auth_logout(AuthenticatedUser)(void) throws UserNotAuthorized(string),
        sendMessage(SendMessageRequest)(void) throws UserNotFound(string) UserNotAuthorized(string),
        getMessages(AuthenticatedUser)(MessageList) throws UserNotAuthorized(string),
}