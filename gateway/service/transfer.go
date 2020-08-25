package service

import (
	"context"
	"github.com/941112341/avalon/sdk/avalon"
	"github.com/941112341/avalon/sdk/avalon/client"
	"github.com/941112341/avalon/sdk/generic"
	"github.com/941112341/avalon/sdk/generic/invoker"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

const baseRespKey = "baseResp"

var TransferServiceInstance = &TransferService{}

var clientMap sync.Map

func Initial() error {
	return avalon.NewBeanFunc(TransferServiceInstance, func() error {
		return nil
	}).Initial()
}

type TransferTarget struct {
	Service string
	Method  string
	Base    string
	PSM     string
	Timeout string
	Retry   int
}

type Response struct {
	BaseResp *BaseResp
	Header   map[string]string
}

type BaseResp struct {
	Code    int32
	Message string
	Data    interface{}
}

type TransferService struct {
	idlGroup *generic.ThriftGroup
	IDLPath  string `viper:"conf.idlpath"`

	lastUpdate time.Time
}

// todo 可以优化
func (t *TransferService) Initial() (err error) {

	if t.idlGroup == nil {
		t.idlGroup = &generic.ThriftGroup{ModelMap: map[string]*generic.ThriftFileModel{}}
	}
	newIdlGroup := &generic.ThriftGroup{ModelMap: map[string]*generic.ThriftFileModel{}}
	fileInfos := inline.AllFIleInDirFunc(t.IDLPath, func(info os.FileInfo) bool {
		return path.Ext(info.Name()) == ".thrift"
	})
	for key, info := range fileInfos {
		fileName := inline.FileName(key)
		if info.ModTime().After(t.lastUpdate) {
			idlGroup, err := generic.NewThriftGroupBase([]string{key})
			if err != nil {
				return err
			}
			newIdlGroup.Merge([]*generic.ThriftGroup{idlGroup})
		} else {
			fileModel, ok := t.idlGroup.ModelMap[fileName]
			if ok {
				newIdlGroup.ModelMap[fileName] = fileModel
			}
		}
	}

	t.idlGroup = newIdlGroup
	t.lastUpdate = time.Now()
	return
}

func (t *TransferService) Destroy() error {
	return nil
}

func (t *TransferService) Transfer(request *http.Request) (*Response, error) {
	urlPath := request.URL.Path
	var target TransferTarget
	err := viper.UnmarshalKey(urlPath, &target)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "viper unmarshal %s", urlPath)
	}

	if time.Now().Sub(t.lastUpdate) > time.Minute {
		if err := t.Initial(); err != nil {
			return nil, inline.PrependErrorFmt(err, "init when transfer")
		}
	}
	invoke, err := invoker.CreateInvoker(t.idlGroup, target.Base, target.Service, target.Method)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "create invoker fail %s", inline.ToJsonString(target))
	}

	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "read all")
	}
	var dataMap map[string]interface{}
	inline.MustUnmarshal(string(data), &dataMap)
	dataMap["base"] = map[string]interface{}{
		"ip":   inline.GetIP(),
		"psm":  viper.GetString("conf.psm"),
		"time": time.Now().Unix(),
		"extra": map[string]string{
			"header": inline.ToJsonString(request.Header),
		},
	}
	argMap := map[string]interface{}{
		"request": dataMap,
	}

	// todo get or default
	var cli *client.AvalonClient
	o, ok := clientMap.Load(target.PSM)
	if ok {
		cli = o.(*client.AvalonClient)
	} else {
		cli = client.DefaultClientTimeout(target.PSM, target.Timeout)
		if err := cli.Initial(); err != nil {
			return nil, err
		}
		cli.Retry = target.Retry
		clientMap.Store(target.PSM, cli)
	}

	response, err := invoke.Invoke(context.Background(), cli, inline.ToJsonString(argMap))
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "request arg %s", inline.ToJsonString(argMap))
	}

	any := inline.JSONAny(response)
	respData := any.Get("success").GetInterface()
	respMap := respData.(map[string]interface{})
	delete(respMap, baseRespKey)
	baseResp := any.Get("success").Get(baseRespKey)
	code := baseResp.Get("code").ToInt32()
	message := baseResp.Get("message").ToString()
	resp := &Response{BaseResp: &BaseResp{
		Code:    code,
		Message: message,
		Data:    respMap,
	}, Header: map[string]string{}}

	extraAny := baseResp.Get("extra").Get("header")
	if extraAny.LastError() == nil {
		header := extraAny.ToString()
		inline.MustUnmarshal(header, &resp.Header)
	}
	return resp, nil
}
