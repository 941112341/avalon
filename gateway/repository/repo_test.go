package repository

import (
	"fmt"
	"github.com/941112341/avalon/gateway/initial"
	"github.com/941112341/avalon/pkg/mygorm"
	"github.com/941112341/avalon/sdk/inline"
	"testing"
)

func TestMapperAll(t *testing.T) {
	err := initial.InitAllForTest()
	if err != nil {
		panic(err)
	}

	repo := mapperRepository{}
	list, err := repo.AllMapper()
	if err != nil {
		panic(err)
	}
	fmt.Println(inline.ToJsonString(list))
}

func TestDel(t *testing.T) {
	err := initial.InitAllForTest()
	if err != nil {
		panic(err)
	}

	repo := mapperRepository{}
	err = repo.DelMapper(MapperList{{Model: mygorm.Model{ID: 1}}})
	if err != nil {
		panic(err)
	}

}

func TestAdd(t *testing.T) {
	err := initial.InitAllForTest()
	if err != nil {
		panic(err)
	}

	repo := mapperRepository{}
	err = repo.AddMapper(MapperList{
		{
			Model:   mygorm.Model{},
			URL:     "hello",
			Type:    0,
			Domain:  "www.jiangshihao.com",
			PSM:     "a.b.c",
			Base:    "idgenerator",
			Method:  "genID",
			Version: "",
		},
		{
			Model:   mygorm.Model{},
			URL:     "hello/api",
			Type:    1,
			Domain:  "www.jiangshihao.com",
			PSM:     "a.b.c",
			Base:    "base",
			Method:  "base",
			Version: "",
		},
	})
	if err != nil {
		panic(err)
	}

}

func TestInsert(t *testing.T) {
	err := initial.InitAllForTest()
	if err != nil {
		panic(err)
	}

	repo := uploadRepository{}
	err = repo.Insert(&UploadVo{
		Model: mygorm.Model{},
		UploadUnionKey: UploadUnionKey{
			UploadGroupKey: UploadGroupKey{
				PSM:     "base",
				Version: "",
			},
			Base: "base",
		},
		Content: "xxxxx",
	})
	if err != nil {
		panic(err)
	}

}

func TestFindOne(t *testing.T) {
	err := initial.InitAllForTest()
	if err != nil {
		panic(err)
	}

	repo := uploadRepository{}

	vo, err := repo.FindByKey(&UploadUnionKey{
		UploadGroupKey: UploadGroupKey{
			PSM:     "base",
			Version: "NPWNRUXEyOZrIrCasywyQvzPiFwMURMn",
		},
		Base: "base",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(inline.ToJsonString(vo))
}

func TestFindGroup(t *testing.T) {
	err := initial.InitAllForTest()
	if err != nil {
		panic(err)
	}

	repo := uploadRepository{}

	vo, err := repo.FindGroup(&UploadGroupKey{
		PSM:     "base",
		Version: "",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(inline.ToJsonString(vo))
}

func TestBatchInsert(t *testing.T) {
	err := initial.InitAllForTest()
	if err != nil {
		panic(err)
	}

	repo := uploadRepository{}

	err = repo.BatchInsert([]*UploadVo{
		{
			Model: mygorm.Model{},
			UploadUnionKey: UploadUnionKey{
				UploadGroupKey: UploadGroupKey{
					PSM:     "base",
					Version: "SscaguqptTXUwWslxVqRvhNQqUDTwFFw",
				},
				Base: "base",
			},
			Content: "xxxxx",
		},
		{
			Model: mygorm.Model{},
			UploadUnionKey: UploadUnionKey{
				UploadGroupKey: UploadGroupKey{
					PSM:     "base",
					Version: "SscaguqptTXUwWslxVqRvhNQqUDTwFFw",
				},
				Base: "base",
			},
			Content: "xxxxx",
		},
	})
	if err != nil {
		panic(err)
	}

}

func TestDeleted(t *testing.T) {
	err := initial.InitAllForTest()
	if err != nil {
		panic(err)
	}

	repo := uploadRepository{}
	err = repo.DeleteGroup(&UploadGroupKey{
		PSM:     "base",
		Version: "0",
	})
	if err != nil {
		panic(err)
	}

}
