
package blog

/*
	1.0.1
*/
import (
    "context"
    "github.com/941112341/avalon/sdk/avalon/server"
)

type AvalonHandler struct {

    advances []server.Advance
    handler  BlogService
}

func NewAvalonHandler(handler BlogService) BlogService {
    return &AvalonHandler{
        advances: server.DefaultAdvance(),
        handler:  handler,
    }
}



func (a *AvalonHandler) GetBlog(ctx context.Context, request *GetBlogRequest) (r *GetBlogResponse, err error) {
    var call server.Call = func(ctx context.Context, request interface{}) (interface{}, error) {
        return a.handler.GetBlog(ctx, request.(*GetBlogRequest))
    }
    for _, advance := range a.advances {
        call = advance(call)
    }
    resp, err := call(ctx, request)
    if err != nil {
        return nil, err
    }
    return resp.(*GetBlogResponse), nil
}


func (a *AvalonHandler) ListBlogs(ctx context.Context, request *ListBlogsRequest) (r *ListBlogsResponse, err error) {
    var call server.Call = func(ctx context.Context, request interface{}) (interface{}, error) {
        return a.handler.ListBlogs(ctx, request.(*ListBlogsRequest))
    }
    for _, advance := range a.advances {
        call = advance(call)
    }
    resp, err := call(ctx, request)
    if err != nil {
        return nil, err
    }
    return resp.(*ListBlogsResponse), nil
}


func (a *AvalonHandler) SaveBlog(ctx context.Context, request *SaveBlogRequest) (r *SaveBlogResponse, err error) {
    var call server.Call = func(ctx context.Context, request interface{}) (interface{}, error) {
        return a.handler.SaveBlog(ctx, request.(*SaveBlogRequest))
    }
    for _, advance := range a.advances {
        call = advance(call)
    }
    resp, err := call(ctx, request)
    if err != nil {
        return nil, err
    }
    return resp.(*SaveBlogResponse), nil
}


func (a *AvalonHandler) DeleteBlogs(ctx context.Context, request *DeleteBlogRequest) (r *DeleteBlogResponse, err error) {
    var call server.Call = func(ctx context.Context, request interface{}) (interface{}, error) {
        return a.handler.DeleteBlogs(ctx, request.(*DeleteBlogRequest))
    }
    for _, advance := range a.advances {
        call = advance(call)
    }
    resp, err := call(ctx, request)
    if err != nil {
        return nil, err
    }
    return resp.(*DeleteBlogResponse), nil
}
