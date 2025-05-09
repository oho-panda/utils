package rd

import (
	"context"
	"fmt"
	"testing"
)

func TestRedis(t *testing.T) {
	ctx := context.Background()
	//InitRedisClient("redis host", "redis password", 6, time.Second*10)
	set := Set(ctx, "test", "111")
	if set {
		fmt.Println("redis test success")
	}
	ok, val := Get(ctx, "test")
	fmt.Println(ok, val)
	fmt.Println(Del(ctx, "test"))
	ok, val = Get(ctx, "test")
	fmt.Println(ok, val)

}
