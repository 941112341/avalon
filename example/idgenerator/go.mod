module github.com/941112341/avalon/example/idgenerator

go 1.14

require (
	github.com/941112341/avalon/common v0.0.0-20200630071424-990dcb7f6e17
	github.com/941112341/avalon/sdk v0.0.0-20200719035616-ff7359fa4a5a
	github.com/facebookgo/inject v0.0.0-20180706035515-f23751cae28b
	github.com/facebookgo/structtag v0.0.0-20150214074306-217e25fb9691 // indirect
	github.com/jinzhu/gorm v1.9.15
	github.com/pkg/errors v0.9.1
)

replace (
	github.com/941112341/avalon/common => ../../common
	github.com/941112341/avalon/sdk => ../../sdk
)
