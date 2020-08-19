package client

import (
	"github.com/941112341/avalon/common/gen/test"
	_cli "github.com/941112341/avalon/sdk/avalon/client"
)

var (
	TestCli *test.CatServiceClient
)

func init() {
	cli, err := _cli.NewClientOptions("example.jiangshihao.test")
	if err != nil {
		panic(err)
	}
	TestCli = test.NewCatServiceClient(cli)
}