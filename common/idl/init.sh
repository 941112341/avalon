#!/usr/bin/env bash


file=$1

thrift --out ../gen  --gen go:package_prefix=github.com/941112341/avalon/common/gen/ ${file}.thrift
