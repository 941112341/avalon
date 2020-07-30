package initial

import (
	"fmt"
	"github.com/941112341/avalon/example/idgenerator/registry"
	"github.com/941112341/avalon/example/idgenerator/service"
	"os"
	"testing"
)

type TestHandler struct {
	S service.GenIdsService `inject:"GenIdsService"`
}

func TestInitAll(t *testing.T) {
	os.Setenv("base", "../base.yaml")
	os.Setenv("conf", "../conf/config.yaml")

	var testHandler TestHandler
	_ = registry.Registry("", &testHandler)
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitAll(tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("InitAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	fmt.Println(testHandler)
}
