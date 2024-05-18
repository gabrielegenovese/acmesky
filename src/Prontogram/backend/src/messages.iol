include "authentication.iol"

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

interface Messages
{
    RequestResponse:
        sendMessage(SendMessageRequest)(Message) throws UserNotFound(string) UserUnauthorized(string),
        getMessages(AuthenticatedUser)(MessageList) throws UserUnauthorized(string)
}