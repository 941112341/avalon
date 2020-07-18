namespace go base

struct Base {
    1: string psm
    2: string ip
    3: i64 time
    4: map<string, string> extra
    5: Base base
}