
package test

/*
	1.0.1
*/
import (
    "context"
    "github.com/941112341/avalon/sdk/avalon/server"
)

type AvalonHandler struct {

    advances []server.Advance
    handler  CatService
}

func NewAvalonHandler(handler CatService) CatService {
    return &AvalonHandler{
        advances: server.DefaultAdvance(),
        handler:  handler,
    }
}



func (a *AvalonHandler) GetCat(ctx context.Context, request *CatRequest) (r *CatResponse, err error) {
    var call server.Call = func(ctx context.Context, request interface{}) (interface{}, error) {
        return a.handler.GetCat(ctx, request.(*CatRequest))
    }
    for _, advance := range a.advances {
        call = advance(call)
    }
    resp, err := call(ctx, request)
    if err != nil {
        return nil, err
    }
    return resp.(*CatResponse), nil
}


func (a *AvalonHandler) GetLittleCat(ctx context.Context, request *LittleCatRequest) (r *LittleCatResponse, err error) {
    var call server.Call = func(ctx context.Context, request interface{}) (interface{}, error) {
        return a.handler.GetLittleCat(ctx, request.(*LittleCatRequest))
    }
    for _, advance := range a.advances {
        call = advance(call)
    }
    resp, err := call(ctx, request)
    if err != nil {
        return nil, err
    }
    return resp.(*LittleCatResponse), nil
}
