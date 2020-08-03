module github.com/941112341/avalon/gateway

go 1.14

require (
	github.com/941112341/avalon/common v0.0.0-20200729103654-81d651c1f161
	github.com/941112341/avalon/example/idgenerator v0.0.0-20200729103654-81d651c1f161
	github.com/941112341/avalon/sdk v0.0.0-20200719035616-ff7359fa4a5a
	github.com/facebookgo/inject v0.0.0-20180706035515-f23751cae28b
	github.com/jinzhu/gorm v1.9.15
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/t-tiger/gorm-bulk-insert/v2 v2.0.1

)

replace github.com/941112341/avalon/sdk => ../sdk

replace github.com/941112341/avalon/common => ../common
