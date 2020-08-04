#!/usr/bin/env bash




# thrift --out ../gen --gen go:package_prefix=github.com/941112341/avalon/common/gen/ idgenerator.thrift
avalon -i idl/test.thrift -o gen/test/generator.go