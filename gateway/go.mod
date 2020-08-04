module github.com/941112341/avalon/gateway

go 1.14

require (
	github.com/941112341/avalon/common v0.0.0-20200803095615-2069c2472e9a
	github.com/941112341/avalon/example/idgenerator v0.0.0-20200803095615-2069c2472e9a
	github.com/941112341/avalon/sdk v0.0.0-20200803095615-2069c2472e9a
	github.com/facebookgo/inject v0.0.0-20180706035515-f23751cae28b
	github.com/jinzhu/gorm v1.9.15
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/t-tiger/gorm-bulk-insert/v2 v2.0.1

)

replace github.com/941112341/avalon/sdk => ../sdk

replace github.com/941112341/avalon/common => ../common

replace github.com/941112341/avalon/example/idgenerator => ../example/idgenerator
