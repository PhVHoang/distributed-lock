package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/go-redis/redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	lockScript = `
		return redis.call('SET', KEYS[1], ARGV[1], 'NX', 'PX', ARGV[2])
	`
	releaseScript = `

		if redis.call("get",KEYS[1]) == ARGV[1] then
		    return redis.call("del",KEYS[1])
		else
		    return 0
		end
	`
)

func AcquireLock(key, value string, timeoutMs int) (bool, error) {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: ":6379",
	})

	_ = rdb.FlushDB(ctx).Err()

	cmd := redis.NewScript(lockScript, []string{key}, []string{
		id, strconv.Itoa(timeoutMs),
	})

	num, err := cmd.Run(ctx, rdb, key)

	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		logx.Errorf("Error on acquring lock for %s, %s".key, err.Error())
		return false, err
	} else if num == nil {
		return false, err
	}

	return true, nil
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	//err = client.Set("name", "Hoang Pham", 0).Err()
	val, err := client.Get("name").Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(val)
}
