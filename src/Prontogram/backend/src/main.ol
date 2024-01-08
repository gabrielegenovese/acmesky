include "interface.iol"
include "console.iol"
include "string_utils.iol"
// include "queue_utils.iol"

inputPort ProntogramServicePort
{
    Location: "socket://localhost:8000"
    Interfaces: IProntogramService
    Protocol: http {
        format = "json"
        // contentType = "application/json"
        osc << {
            auth_signup << {
                template = "/api/auth/signup"
                method = "put"
                statusCodes.UserAlreadyExists = 403
            }
            auth_login << {
                template = "/api/auth/login"
                method = "post"
                statusCodes.UserNotAuthorized = 401
                statusCodes.UserNotFound = 404
            }
            auth_logout << {
                template = "/api/auth/logout"
                method = "post"
                statusCodes.UserNotAuthorized = 401
                statusCodes.UserNotFound = 404
            }
            // get stored messages
            getMessages << {
                template = "/api/messages"
                method = "get"
                statusCodes.UserNotAuthorized = 401
            }
            sendMessage << {
                template = "/api/messages"
                method = "put"
                statusCodes.UserNotAuthorized = 401
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
            if ( is_defined( global.users.(UserSignUpRequest.credentials.user_id) ) ) {
                throw (UserAlreadyExists, UserSignUpRequest.credentials.user_id)
            }
            
            with(global.users.(UserSignUpRequest.credentials.user_id)) {
                .id = UserSignUpRequest.credentials.user_id;
                .display_name = UserSignUpRequest.display_name;
                .password = UserSignUpRequest.credentials.password
            }

            println@Console("New user signed up: '" + global.users.(UserSignUpRequest.credentials.user_id).id + "'")()
        }
    ]
    [
        auth_login(UserAuthCredentials)(AuthenticatedUser) {
            if ( !is_defined( global.users.(UserAuthCredentials.user_id) ) ) {
                throw (UserNotFound, UserAuthCredentials.user_id)
            }
            else if( global.users.(UserAuthCredentials.user_id).password != UserAuthCredentials.password ) {
                throw (UserNotAuthorized, UserAuthCredentials.user_id)
            }
            else {
                csets.sid = new
                with(AuthenticatedUser) {
                    .user_id = UserAuthCredentials.user_id;
                    .sid = csets.sid
                }
                isLogged = true
                println@Console("user '" + AuthenticatedUser.user_id + "' logged in")()
                println@Console("user '" + AuthenticatedUser.user_id + " started session " + AuthenticatedUser.sid )()
            }
        }
    ]
    [
        auth_logout(AuthenticatedUser)() {
            if ( !is_defined( global.users.(AuthenticatedUser.user_id) ) ) {
                throw (UserNotFound, UserAuthCredentials.user_id)
            }
            else {
                isLogged = false
            }
            println@Console("user '" + AuthenticatedUser.user_id + "' logged out")()
        }
    ]
    [
        sendMessage(SendMessageRequest)(Message){
            with(Message) {
                .id = new;
                .content = SendMessageRequest.content;
                .sender_user_id = SendMessageRequest.sender.user_id;
                .receiver_user_id = SendMessageRequest.receiver_user_id
            }
            if ( !is_defined( global.users.(Message.receiver_user_id) ) ) {
                throw (UserNotFound, Message.receiver_user_id)
            }
            inbox -> global.inbox.(Message.receiver_user_id)
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
            inbox -> global.inbox.(AuthenticatedUser.user_id)
            synchronized(inboxLock) {
                for( i = 0, i < #inbox, i++ ) {
                    if( is_defined(inbox[i])) {
                        MessageList.messages[#MessageList.messages] << inbox[i]
                    }
                }
                
                // undef(global.inbox.(AuthenticatedUser.user_id))
            }

            valueToPrettyString@StringUtils(MessageList)( s )
            println@Console("Messages " + s)()

            valueToPrettyString@StringUtils(inbox)( s )
            println@Console("Messages left " + s)()
        }
    ]
}