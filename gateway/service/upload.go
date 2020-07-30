package service

import (
	"fmt"
	"github.com/941112341/avalon/common/client"
	"github.com/941112341/avalon/gateway/model"
	"github.com/941112341/avalon/gateway/registry"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/patrickmn/go-cache"
	"time"
)

type uploaderService struct {
	Repo model.UploadRepository `inject:"UploadRepository"`
}

func (u *uploaderService) Upload(request *IDLFileVo) error {
	m := newModel(u.Repo, request)
	if err := m.Upload(); err != nil {
		return inline.PrependErrorFmt(err, "upload %s", inline.ToJsonString(request))
	}
	return nil
}

func newModel(repo model.UploadRepository, request *IDLFileVo) *model.IDLFile {
	m := &model.IDLFile{
		IDLFileID: model.IDLFileID{
			PSM:  request.PSM,
			Base: request.Filename,
		},
		Content: request.Body,
		Version: 0,
		Deleted: inline.BoolPtr(false),
		Created: time.Now(),
		Updated: time.Now(),
		Repo:    repo,
	}
	return m
}

func (u *uploaderService) Get(request *IDLFileVo) (*IDLFileVo, error) {
	m := newModel(u.Repo, request)
	file, err := m.Get()
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "get fail %s", inline.ToJsonString(request))
	}

	return &IDLFileVo{
		PSM:      file.PSM,
		Filename: file.Base,
		Body:     file.Content,
	}, nil
}

type CacheUploader struct {
	cache *cache.Cache
	Proxy UploadService `inject:"UploadService"`
}

func UnionKey(request *IDLFileVo) string {
	return fmt.Sprintf("%s_%s", request.Filename, request.PSM)
}

func (c *CacheUploader) Upload(request *IDLFileVo) error {
	key := UnionKey(request)
	defer c.cache.Delete(key)
	return c.Proxy.Upload(request)
}

func (c *CacheUploader) Get(request *IDLFileVo) (*IDLFileVo, error) {
	key := UnionKey(request)
	i, _ := c.cache.Get(key)
	cacheObj, ok := i.(*IDLFileVo)
	if ok {
		return cacheObj, nil
	}
	idlFile, err := c.Proxy.Get(request)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "Get upload err: request %s", inline.ToJsonString(request))
	}
	c.cache.Set(key, idlFile, time.Hour)
	return idlFile, nil
}

func init() {
	_ = registry.Registry("UploadService", &uploaderService{})
	_ = registry.Registry("CacheUploadService", &CacheUploader{
		cache: cache.New(time.Hour, time.Hour*2),
	})
}
