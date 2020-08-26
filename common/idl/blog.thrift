namespace go blog

include 'base.thrift'

struct Blog {
    1: string id
    2: string imageURL
    3: string content
    4: list<string> categories
    5: string name
    6: string description
    7: string updated
}

struct GetBlogRequest {
    1: string id
    255: base.Base base
}

struct GetBlogResponse {
    1: optional Blog blog
    255: base.BaseResp baseResp
}

struct ListBlogsRequest {
    1: string keyword
    2: i32 page
    3: i32 size

    255: base.Base base
}

struct ListBlogsResponse {
    1: list<Blog> blogs
    2: i64 count
    255: base.BaseResp baseResp
}

struct SaveBlogRequest {
    1: optional string id
    2: string content

    255: base.Base base
}

struct SaveBlogResponse {
    1: string id
    2: string content

    255: base.BaseResp baseResp
}

struct DeleteBlogRequest {
    1: list<string> ids

    255: base.Base base
}

struct DeleteBlogResponse {
    1: list<string> ids

    255: base.BaseResp baseResp
}

service blogService {
    GetBlogResponse GetBlog(1: GetBlogRequest request)
    ListBlogsResponse ListBlogs(1: ListBlogsRequest request)
    SaveBlogResponse SaveBlog(1: SaveBlogRequest request)
    DeleteBlogResponse DeleteBlogs(1: DeleteBlogRequest request)
}