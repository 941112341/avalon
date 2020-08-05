package service

import (
	"fmt"
	"github.com/941112341/avalon/gateway/model"
	"github.com/941112341/avalon/gateway/registry"
	"github.com/941112341/avalon/gateway/repository"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/patrickmn/go-cache"
	"time"
)

type uploaderService struct {
	Repo repository.UploadRepository `inject:"UploadRepository"`
}

func (u *uploaderService) Upload(request *UploadVoVo) error {
	m := newModel(request)
	if err := m.Upload(u.Repo); err != nil {
		return inline.PrependErrorFmt(err, "upload %s", inline.ToJsonString(request))
	}
	return nil
}

func newModel(request *UploadVoVo) *model.UploadVo {
	m := &model.UploadVo{
		UploadVoID: model.UploadVoID{
			PSM:  request.PSM,
			Base: request.Filename,
		},
		Content: request.Body,
		Version: 0,
		Deleted: inline.BoolPtr(false),
		Created: time.Now(),
		Updated: time.Now(),
	}
	return m
}

func (u *uploaderService) Get(request *UploadVoVo) (*UploadVoVo, error) {
	m := newModel(request)
	file, err := m.Get(u.Repo)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "get fail %s", inline.ToJsonString(request))
	}

	return &UploadVoVo{
		PSM:      file.PSM,
		Filename: file.Base,
		Body:     file.Content,
	}, nil
}

type CacheUploader struct {
	cache *cache.Cache
	Proxy UploadService `inject:"UploadService"`
}

func UnionKey(request *UploadVoVo) string {
	return fmt.Sprintf("%s_%s", request.Filename, request.PSM)
}

func (c *CacheUploader) Upload(request *UploadVoVo) error {
	key := UnionKey(request)
	defer c.cache.Delete(key)
	return c.Proxy.Upload(request)
}

func (c *CacheUploader) Get(request *UploadVoVo) (*UploadVoVo, error) {
	key := UnionKey(request)
	i, _ := c.cache.Get(key)
	cacheObj, ok := i.(*UploadVoVo)
	if ok {
		return cacheObj, nil
	}
	UploadVo, err := c.Proxy.Get(request)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "Get upload err: request %s", inline.ToJsonString(request))
	}
	c.cache.Set(key, UploadVo, time.Hour)
	return UploadVo, nil
}

func init() {
	_ = registry.Registry("UploadService", &uploaderService{})
	_ = registry.Registry("CacheUploadService", &CacheUploader{
		cache: cache.New(time.Hour, time.Hour*2),
	})
}
