module github.com/941112341/avalon/example/idgenerator

go 1.14

require (
	github.com/941112341/avalon/common v0.0.0-20200803095615-2069c2472e9a
	github.com/941112341/avalon/sdk v0.0.0-20200719035616-ff7359fa4a5a
	github.com/facebookgo/ensure v0.0.0-20200202191622-63f1cf65ac4c // indirect
	github.com/facebookgo/inject v0.0.0-20180706035515-f23751cae28b
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/facebookgo/structtag v0.0.0-20150214074306-217e25fb9691 // indirect
	github.com/facebookgo/subset v0.0.0-20200203212716-c811ad88dec4 // indirect
	github.com/jinzhu/gorm v1.9.15
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.3.0
)

replace (
	github.com/941112341/avalon/common => ../../common
	github.com/941112341/avalon/sdk => ../../sdk
)
