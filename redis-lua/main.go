package main

import (
	"context"
	"github.com/redis/rueidis"
	"play-around/common"
)

func main() {
	redis := common.InitRedis()
	acqat := rueidis.NewLuaScript(`local r = redis.call("SET",KEYS[1],ARGV[1],"NX","PXAT",ARGV[2]);redis.call("GET",KEYS[1]);return r`)

	resp := acqat.Exec(context.Background(), redis, []string{"test-01"}, []string{"value-01", "1824660520047"})

	if err := resp.Error(); rueidis.IsRedisNil(err) {
		println("Error: ", err.Error())
	} else {
		println(resp.String())
	}
}
