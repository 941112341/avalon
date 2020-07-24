namespace go idgenerator

include "base.thrift"

struct IDRequest {
    1: i32 count

    255: base.Base base
}

struct IDResponse {
    1: list<i64> IDs

    255: base.BaseResp baseResp
}

service IDGenerator {
    IDResponse GenIDs(1: IDRequest request)

}