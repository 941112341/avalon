namespace go test

include "base.thrift"


struct Cat {
    1: required i32 age
    2: optional string name

    3: list<LittleCat> babies
}

struct LittleCat {
    1: Cat cat
    2: i32 age
    3: i8 color
    4: map<i64, list< Foo> > ids
}

struct CatRequest {
    1: list<i64> id

    255: base.Base base
}

struct CatResponse {
    1: map<i64, Cat> cats

    255: base.BaseResp baseResp
}

struct Foo {
    1: bool love
}

struct LittleCatRequest {
    1: optional Cat cat
    255: base.Base base
}

struct LittleCatResponse {
    1: optional list< LittleCat > littleCat
    255: base.BaseResp baseResp
}

service CatService {
    CatResponse GetCat(1: CatRequest request)
    LittleCatResponse GetLittleCat(1: LittleCatRequest request)
}