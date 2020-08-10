package common

import (
	"context"
	"errors"
	"github.com/941112341/avalon/common/gen/test"
	"reflect"
	"testing"
)

func TestConvertResponse(t *testing.T) {
	var CatResponse *test.CatResponse
	type args struct {
		ctx      context.Context
		request  interface{}
		response interface{}
		err      error
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			args: args{
				ctx:      nil,
				request:  nil,
				response: CatResponse,
				err:      errors.New("err happen"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertResponse(tt.args.request, tt.args.response, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertResponse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
