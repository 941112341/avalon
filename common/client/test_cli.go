package client

import (
	"github.com/941112341/avalon/common/gen/test"
	_cli "github.com/941112341/avalon/sdk/avalon/client"
)

var (
	TestCli *test.CatServiceClient
)

func init() {
	cli := _cli.DefaultClient("example.jiangshihao.test")
	err := cli.Initial()
	if err != nil {
		panic(err)
	}
	TestCli = test.NewCatServiceClient(cli)
}