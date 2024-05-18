include "protocols/http.iol"

type UserAuthCredentials: void {
    .userId: string
    .password: string
}

type UserSignUpRequest: void {
    .display_name: string
    .credentials: UserAuthCredentials
}

type AuthenticatedUser: void {
    .userId: string
    .sid: string
}

type SendMessageRequest: void {
    .sender: AuthenticatedUser
    .receiverUserId: string
    .content: string
}

type Message: void {
    .id: string
    .sender_userId: string
    .receiverUserId: string
    .content: string
}

type MessageList: void {
    .messages[0,*]: Message
}

interface IProntogramService
{

    RequestResponse:
        auth_signup(UserSignUpRequest)(void) throws UserAlreadyExists(string),
        auth_login(UserAuthCredentials)(AuthenticatedUser) throws UserNotFound(string) UserUnauthorized(string),
        auth_logout(AuthenticatedUser)(void) throws UserUnauthorized(string),
        sendMessage(SendMessageRequest)(Message) throws UserNotFound(string) UserUnauthorized(string),
        getMessages(AuthenticatedUser)(MessageList) throws UserUnauthorized(string),
        default(DefaultOperationHttpRequest)(void)
}