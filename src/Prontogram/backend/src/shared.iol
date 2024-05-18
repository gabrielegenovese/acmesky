include "protocols/http.iol"

interface Shared
{
    RequestResponse:
        default(DefaultOperationHttpRequest)(void)
}