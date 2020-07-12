package avalon

import (
	"fmt"
	"github.com/941112341/avalon/sdk/inline"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	os.Setenv("base", "../../base.yaml")

	cfg, err := GetConfig()
	fmt.Println(inline.ToJsonString(cfg), err)

}
