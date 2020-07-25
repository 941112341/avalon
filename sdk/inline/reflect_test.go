package inline

import (
	"errors"
	"fmt"
	"testing"
)

type Hello struct {
	Base Base
}

type Base struct {
	PSM      string
	HostPort string
}

func TestSetField(t *testing.T) {
	var hello Hello
	type args struct {
		any       interface{}
		fieldName string
		to        interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{
				any:       &hello,
				fieldName: "Base",
				to: Base{
					PSM:      "123",
					HostPort: "4567",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetField(tt.args.any, tt.args.fieldName, tt.args.to); (err != nil) != tt.wantErr {
				t.Errorf("SetField() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	fmt.Println(hello)
}

func TestLog(t *testing.T) {
	err := errors.New("?")

	WithFields("err", err).Errorln("?")
}
