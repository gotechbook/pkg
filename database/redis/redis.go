package redis

import (
	"context"
	"github.com/gotechbook/pkg/logger"
	"github.com/redis/go-redis/v9"
	"time"
)

type Client struct {
	Db redis.UniversalClient
}

func NewClient(endpoints []string, username string, password string, dataBase, dialTimeout, readTimeout, writeTimeout int) *Client {
	logger.Debugf("connecting redis")
	opt := redis.UniversalOptions{
		Addrs:        endpoints,
		DB:           dataBase,
		DialTimeout:  time.Second * time.Duration(dialTimeout),
		ReadTimeout:  time.Second * time.Duration(readTimeout),
		WriteTimeout: time.Second * time.Duration(writeTimeout),
		Username:     username,
		Password:     password,
		//命令执行失败时的重试策略
		MaxRetries:      0,                      //命令执行失败时，最多重试多少次，默认为0即不重试
		MinRetryBackoff: 8 * time.Millisecond,   //每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		MaxRetryBackoff: 512 * time.Millisecond, //每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔
	}
	client := redis.NewUniversalClient(&opt)
	if err := client.Ping(context.TODO()).Err(); err != nil {
		logger.Fatalf("failed to connect redis: %v", err)
	}
	return &Client{client}
}
