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
        osc << {
            signup << {
                template = "/api/auth/signup"
                method = "put"
                statusCodes.UserAlreadyExists = 403
            }
            login << {
                template = "/api/auth/login"
                method = "post"
                statusCodes.UserNotAuthorized = 401
                statusCodes.UserNotFound = 404
            }
            logout << {
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

init {
    global.users = {}
    global.inbox = {}
}

cset {
    sid: AuthenticatedUser.sid
}

main {

    while( true ) {
        scope(scope1) {

            install(
                UserNotFound => {
                    println@Console("User not found: " + scope1.UserNotFoundError)()
                },
                UserNotAuthorized => {
                    println@Console("User Not Authorized: " + scope1.UserNotAuthorizedError)()
                },
                UserAlreadyExists => {
                    println@Console("User already exists (cant create): " + scope1.UserAlreadyExistsError)()
                }
            )

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
                        with(AuthenticatedUser) {
                            .user_id = UserAuthCredentials.user_id;
                            .sid = csets.sid = new
                        }
                    }
                }
            ]
            [
                auth_logout(AuthenticatedUser)() {
                    if ( !is_defined( global.users.(AuthenticatedUser.user_id) ) ) {
                        throw (UserNotFound, UserAuthCredentials.user_id)
                    }
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

                    synchronized(inboxLock) {
                        idx = #global.inbox.(Message.receiver_user_id)
                        global.inbox.(Message.receiver_user_id)[idx] << Message
                    }

                    valueToPrettyString@StringUtils(Message)( s )
                    println@Console("Message " + s)()
                }
            ]
            [
                getMessages(AuthenticatedUser)(MessageList){
                    
                    synchronized(inboxLock) {
                        for( i = 0, i < #global.inbox.(AuthenticatedUser.user_id), i++ ) {
                            message = global.inbox.(AuthenticatedUser.user_id)[i]
                            if( is_defined(message)) {
                                MessageList.messages[#MessageList.messages] = message
                            }
                        }
                        
                        // undef(global.inbox.(AuthenticatedUser.user_id))
                    }

                    valueToPrettyString@StringUtils(MessageList)( s )
                    println@Console("Messages " + s)()

                    valueToPrettyString@StringUtils(global.inbox.(AuthenticatedUser.user_id))( s )
                    println@Console("Messages left " + s)()
                }
            ]
        }
    }
}