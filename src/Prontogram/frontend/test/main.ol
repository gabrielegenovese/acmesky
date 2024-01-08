include "../../backend/src/interface.iol"
include "console.iol"
include "string_utils.iol"

outputPort ProntogramService
{
    Location: "socket://localhost:8000"
    Interfaces: IProntogramService
    Protocol: http {
        format = "json"
        //contentType = "application/json"
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


init {
    with(newUserA){
        .display_name = "A test user";
        .credentials << {
            user_id = "A"
            password = "123A"
        }
    }

    with(newUserB){
        .display_name = "B test user";
        .credentials << {
            user_id = "B"
            password = "123B"
        }
    }
}

define test_signup {
    println@Console("test_signup()")()

    auth_signup@ProntogramService(newUserA)()
    valueToPrettyString@StringUtils(newUserA)( sA )
    println@Console("A Signed up as\n" + sA)()

    auth_signup@ProntogramService(newUserB)()
    valueToPrettyString@StringUtils(newUserB)( sB )
    println@Console("B Signed up as\n" + sB)()
}

define test_login {
    println@Console("test_login()")()


    println@Console("logging in as\n" + newUserB.credentials.user_id)()
    auth_login@ProntogramService(newUserB.credentials)(AuthenticatedUserB)
    valueToPrettyString@StringUtils(AuthenticatedUserB)( sB )
    println@Console("Logged in as\n" + sB)()

    println@Console("logging in as\n" + newUserA.credentials.user_id)()
    auth_login@ProntogramService(newUserA.credentials)(AuthenticatedUserA)
    valueToPrettyString@StringUtils(AuthenticatedUserA)( sA )
    println@Console("Logged in as\n" + sA)()


}

define test_read_empty_inbox {
    println@Console("test_read_empty_inbox()")()

    valueToPrettyString@StringUtils(AuthenticatedUserA)( sA )
    println@Console("get messages as\n" + sA)()
    getMessages@ProntogramService(AuthenticatedUserA)(MessageListA)
    if( #MessageListA.messages > 0 ) {
        valueToPrettyString@StringUtils(MessageListA)( sA )
        println@Console("(A) Messages list:\n" + sA)()
        throw (TestFailed, "test_read_empty_inbox")
    }

    valueToPrettyString@StringUtils(AuthenticatedUserB)( sB )
    println@Console("get messages as\n" + sB)()
    getMessages@ProntogramService(AuthenticatedUserB)(MessageListB)
    if( #MessageListB.messages > 0 ) {
        valueToPrettyString@StringUtils(MessageListB)( sB )
        println@Console("(B) Messages list:\n" + sB)()
        throw (TestFailed, "test_read_empty_inbox")
    }
}

define test_sendMessage {
    println@Console("test_sendmessages()")()

    {
        with(Amsg2B) {
            .sender -> AuthenticatedUserA;
            .receiver_user_id = "B";
            .content = "Hi B, i'm A"
        }

        sendMessage@ProntogramService(Amsg2B)(AFirstMessage)
        getMessages@ProntogramService(AuthenticatedUserB)(MessageListB)

        if( #MessageListB.messages != 1 ) {
            valueToPrettyString@StringUtils(MessageListB)( sB )
            println@Console("(B) Messages list:\n" + sB)()
            throw (TestFailed, "test_sendMessage")
        }
    }
    ;
    {
        with(Bmsg2A) {
            .sender -> AuthenticatedUserB;
            .receiver_user_id = "A";
            .content = "Hi A, i'm B"
        }

        sendMessage@ProntogramService(Bmsg2A)(BFirstMessage)
        getMessages@ProntogramService(AuthenticatedUserA)(MessageListA)

        if( #MessageListA.messages != 1 ) {
            valueToPrettyString@StringUtils(MessageListA)( sA )
            println@Console("(A) Messages list:\n" + sA)()
            throw (TestFailed, "test_sendMessage")
        }
    }

}

main {

    test_signup

    test_login

    test_read_empty_inbox

    test_sendMessage

}