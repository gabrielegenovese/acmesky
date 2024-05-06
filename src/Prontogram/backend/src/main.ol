include "interface.iol"
include "console.iol"
include "string_utils.iol"
// include "queue_utils.iol"

inputPort ProntogramServicePort
{
    Location: "socket://localhost:8080"
    Interfaces: IProntogramService
    Protocol: http {
        format = "json"
        // contentType = "application/json"
        osc << {
            auth_signup << {
                template = "/api/users"
                method = "post"
                statusCodes.UserAlreadyExists = 403
            }
            auth_login << {
                template = "/api/auth/{userId}/login"
                method = "post"
                statusCodes.UserUnauthorized = 401
                statusCodes.UserNotFound = 404
            }
            auth_logout << {
                template = "/api/auth/{userId}/logout"
                method = "post"
                statusCodes.UserUnauthorized = 401
                statusCodes.UserNotFound = 404
            }
            // get stored messages
            getMessages << {
                template = "/api/users/{userId}/messages"
                method = "get"
                statusCodes.UserUnauthorized = 401
            }
            sendMessage << {
                template = "/api/users/{receiverUserId}/messages"
                method = "post"
                statusCodes.UserUnauthorized = 401
                statusCodes.UserNotFound = 404
            }
        }
    }
    
}

execution { concurrent }

init {
    global.users = {}
    global.inbox = {}
}

cset {
    sid: AuthenticatedUser.sid
}

main {

    [
        auth_signup(UserSignUpRequest)() {
            if ( is_defined( global.users.(UserSignUpRequest.credentials.userId) ) ) {
                throw (UserAlreadyExists, UserSignUpRequest.credentials.userId)
            }
            
            with(global.users.(UserSignUpRequest.credentials.userId)) {
                .id = UserSignUpRequest.credentials.userId;
                .display_name = UserSignUpRequest.display_name;
                .password = UserSignUpRequest.credentials.password
            }

            println@Console("New user signed up: '" + global.users.(UserSignUpRequest.credentials.userId).id + "'")()
        }
    ]
    [
        auth_login(UserAuthCredentials)(AuthenticatedUser) {
            if ( !is_defined( global.users.(UserAuthCredentials.userId) ) ) {
                throw (UserNotFound, UserAuthCredentials.userId)
            }
            else if( global.users.(UserAuthCredentials.userId).password != UserAuthCredentials.password ) {
                throw (UserUnauthorized, UserAuthCredentials.userId)
            }
            else {
                csets.sid = new
                with(AuthenticatedUser) {
                    .userId = UserAuthCredentials.userId;
                    .sid = csets.sid
                }
                isLogged = true
                println@Console("user '" + AuthenticatedUser.userId + "' logged in")()
                println@Console("user '" + AuthenticatedUser.userId + " started session " + AuthenticatedUser.sid )()
            }
        }
    ]
    [
        auth_logout(AuthenticatedUser)() {
            if ( !is_defined( global.users.(AuthenticatedUser.userId) ) ) {
                throw (UserNotFound, UserAuthCredentials.userId)
            }
            else {
                isLogged = false
            }
            println@Console("user '" + AuthenticatedUser.userId + "' logged out")()
        }
    ]
    [
        sendMessage(SendMessageRequest)(Message){
            with(Message) {
                .id = new;
                .content = SendMessageRequest.content;
                .sender_userId = SendMessageRequest.sender.userId;
                .receiverUserId = SendMessageRequest.receiverUserId
            }
            if ( !is_defined( global.users.(Message.receiverUserId) ) ) {
                throw (UserNotFound, Message.receiverUserId)
            }
            inbox -> global.inbox.(Message.receiverUserId)
            synchronized(inboxLock) {
                idx = #inbox
                inbox[idx] << Message
            }

            valueToPrettyString@StringUtils(Message)( s )
            println@Console("Message " + s)()
        }
    ]
    [
        getMessages(AuthenticatedUser)(MessageList){
            inbox -> global.inbox.(AuthenticatedUser.userId)
            synchronized(inboxLock) {
                for( i = 0, i < #inbox, i++ ) {
                    if( is_defined(inbox[i])) {
                        MessageList.messages[#MessageList.messages] << inbox[i]
                    }
                }
                
                // undef(global.inbox.(AuthenticatedUser.userId))
            }

            valueToPrettyString@StringUtils(MessageList)( s )
            println@Console("Messages " + s)()

            valueToPrettyString@StringUtils(inbox)( s )
            println@Console("Messages left " + s)()
        }
    ]
}