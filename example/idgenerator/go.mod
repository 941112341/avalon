module github.com/941112341/avalon/example/idgenerator

go 1.14

require (
	github.com/941112341/avalon/sdk v0.0.0-20200630084848-f769dffea5fe
)

replace (
    github.com/941112341/avalon/sdk =>  ../sdk
    github.com/941112341/avalon/common => ../common
)