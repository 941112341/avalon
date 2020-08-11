module github.com/941112341/avalon/example/test

go 1.14

require (
	github.com/941112341/avalon/common v0.0.0-20200803095615-2069c2472e9a
	github.com/941112341/avalon/sdk v0.0.0-20200810130135-4fae6bb39015
)

replace (
	github.com/941112341/avalon/common => ../../common
	github.com/941112341/avalon/sdk => ../../sdk
)
