namespace go base

include "base.thrift"

struct MessageRequest {
    1: map<string, string> header
    2: binary body
    3: string methodName
    4: string url

    255: base.Base base
}

struct MessageResponse {
    1: map<string, string> header
    2: binary body
    3: i32 status

    255: base.BaseResp baseResp
}

service messageService {
    MessageResponse MessageDispatcher(1: MessageRequest request)
}