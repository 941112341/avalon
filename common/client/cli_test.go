package client

import (
	"fmt"
	"os"
	"testing"
)


func TestCacheIDClient_GenID(t *testing.T) {

	os.Setenv("base", "../../example/idgenerator/base.yaml")
	cli := NewCacheIDClient(5)
	for i := 0; i < 1; i++ {
		fmt.Println(cli.GenID())
	}
}
