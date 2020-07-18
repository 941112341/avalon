package client

import (
	"fmt"
	"github.com/941112341/avalon/sdk/avalon"
	"github.com/941112341/avalon/sdk/idl/message/base"
	"github.com/941112341/avalon/sdk/inline"
	"testing"
)

func TestSetFieldJSON(t *testing.T) {
	type args struct {
		any       interface{}
		fieldName string
		to        interface{}
	}
	request := base.MessageRequest{}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{
				any:       &request,
				fieldName: "Base",
				to: avalon.Base{
					PSM: "123",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := inline.SetFieldJSON(tt.args.any, tt.args.fieldName, tt.args.to); (err != nil) != tt.wantErr {
				t.Errorf("SetFieldJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	fmt.Println(request)
}
