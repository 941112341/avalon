#!/usr/bin/env bash


thrift --out ../../common/gen  --gen go ../../common/idl/idgenerator.thrift


# thrift --out ../gen --gen go:package_prefix=github.com/941112341/avalon/common/gen/ idgenerator.thrift
# avalon -i idl/idgenerator.thrift -o gen/idgenerator/generator2.go