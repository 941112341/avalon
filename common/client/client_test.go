package client

import (
	"context"
	"fmt"
	"testing"
)


func TestClient2(t *testing.T) {
	ctx := context.Background()
	ids, err := MultiIDs(ctx, 10)
	fmt.Println(ids, err)
}
