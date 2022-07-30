package main

import (
	"context"
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

func AcquireLock(key, value string, timeoutMs int, rdb redis.Client) (bool, error) {
	ctx := context.Background()

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

	AcquireLock("your_name", "name", 10000, client)

}
