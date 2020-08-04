package main

import (
	"context"
	"fmt"
	"github.com/941112341/avalon/common/gen/test"
	"github.com/941112341/avalon/sdk/inline"
	"os"
)

type Handler struct {
}

func (h Handler) GetCat(ctx context.Context, request *test.CatRequest) (r *test.CatResponse, err error) {
	return &test.CatResponse{
		Cats: map[int64]*test.Cat{
			/*1: {
				Age:    1,
				Name:   nil,
				Babies: nil,
			},*/
			2: {
				//Age:  20,
				//Name: inline.StringPtr("Tom"),
				Babies: []*test.LittleCat{
					{
						Cat:   nil,
						Age:   5,
						Color: 0,
						Ids: map[int64][]*test.Foo{
							17: {
								{
									Love: false,
								},
							},
						},
					},
				},
			},
		},
	}, nil
}

func (h Handler) GetLittleCat(ctx context.Context, request *test.LittleCatRequest) (r *test.LittleCatResponse, err error) {
	return &test.LittleCatResponse{LittleCat: []*test.LittleCat{
		{
			Cat: &test.Cat{
				Age:    5,
				Name:   inline.StringPtr("Jimmy"),
				Babies: nil,
			},
			Age:   20,
			Color: 1,
			Ids:   nil,
		},
	}}, nil
}

func main() {

	fmt.Println(os.Getwd())

	err := test.Run("example.jiangshihao.test", &Handler{})
	if err != nil {
		panic(err)
	}
}
