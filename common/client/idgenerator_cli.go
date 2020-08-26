package client

import (
	"context"
	"fmt"
	"github.com/941112341/avalon/common/gen/idgenerator"
	_cli "github.com/941112341/avalon/sdk/avalon/client"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/bwmarrin/snowflake"
	"sync"
)

var (
	client *idgenerator.IDGeneratorClient

	cacheClient *CacheIDClient
)

func init() {
	cli := _cli.DefaultClient("avalon.test.idgenerator")
	err := cli.Initial()
	if err != nil {
		panic(err)
	}
	client = idgenerator.NewIDGeneratorClient(cli)

	cacheClient = NewCacheIDClient(0)
}


func MultiIDs(ctx context.Context, cnt int) (ids []int64, err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		rerr, ok := r.(error)
		if !ok {
			inline.WithFields("recover", r).Errorln("convert err fail")
		}
		err = rerr
	}()
	resp, err := client.GenIDs(ctx, &idgenerator.IDRequest{
		Count: int32(cnt),
	})
	fmt.Println(resp, err)
	if err != nil {
	    return nil, err
	}
	if resp.BaseResp.Code != 0 {
		return nil, fmt.Errorf("rpc err[%d:%s]", resp.BaseResp.Code, resp.BaseResp.Message)
	}
	return resp.IDs, nil
}


type CacheIDClient struct {

	node *snowflake.Node
	cache []int64
	lock sync.Mutex
}

func NewCacheIDClient(num int) *CacheIDClient {
	node, _ := snowflake.NewNode(int64(num))
	return &CacheIDClient{node: node}
}

func (c *CacheIDClient) Length() int {
	return len(c.cache)
}

func (c *CacheIDClient) canAssign(cnt int) bool {
	return c.Length() >= cnt
}

func (c *CacheIDClient) Assign(cnt int) []int64 {
	if cnt > 100 {
		cnt = 100
	}
	if cnt < 1 {
		cnt = 1
	}
	c.lock.Lock()
	defer c.lock.Unlock()

	for !c.canAssign(cnt) {
		ids, err := MultiIDs(context.Background(), 100)
		if err != nil {
			list := make([]int64, 0)
			for i := 0; i < cnt; i++ {
				elem := c.node.Generate().Int64()
				list = append(list, elem)
			}
			return list
		}
		c.cache = append(c.cache, ids...)
	}
	r := c.cache[:cnt]
	c.cache = c.cache[cnt:]
	return r
}


func (c *CacheIDClient) GenID() int64 {
	return c.Assign(1)[0]
}

func GenID() int64 {
	return cacheClient.GenID()
}