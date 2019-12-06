package redis

import (
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func GetMaxGroupID() int {
	if value, err := redisClient.Get("maxgourpid").Result(); err != nil {
		fmt.Println("no maxgroupid, ", err)
		return -1
	} else {
		i, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("parse value error ", err)
			return -1
		}
		return i
	}
}

func AddMaxGroupID() int {
	if value, err := redisClient.Get("maxgroupid").Result(); err != nil {
		fmt.Println("no maxgroupid, ", err)
		return -1
	} else {
		i, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("parse value error ", err)
			return -1
		}
		i += 1
		redisClient.Set("maxgroupid", strconv.Itoa(i), 0)
		return i
	}
}

func GetNewEventID() int {
	if value, err := redisClient.Get("maxeventid").Result(); err != nil {
		fmt.Println("no maxeventid, ", err)
		return -1
	} else {
		i, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("parse value error ", err)
			return -1
		}
		i += 1
		redisClient.Set("maxeventid", strconv.Itoa(i), 0)
		return i
	}
}

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
