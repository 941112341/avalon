package inline

import (
	"fmt"
	"testing"
)

func TestUnwrap(t *testing.T) {
	type args struct {
		r       string
		content string
	}
	tests := []struct {
		name  string
		args  args
		wantS string
	}{
		{
			args: args{
				r:       "<(.*)>",
				content: "list<list<i18>>",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotS := Unwrap(tt.args.r, tt.args.content); gotS != tt.wantS {
				t.Errorf("Unwrap() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}

func TestTemplateExtract(t *testing.T) {

	maps, err := TemplateExtract("/{api}i/webhook/{hook}", "/api/webhook/idgenerator")
	if err != nil {
		panic(err)
	}

	fmt.Println(ToJsonString(maps))
}
