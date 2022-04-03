package common

import (
	"context"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-redis/redis/v8"
	"time"
)

type QueueConfig struct {
	Name     string
	Address  string
	Password string
	DB       int
	Timeout  time.Duration
}

type Queue struct {
	Name     string
	Address  string
	Password string
	DB       int
	rdb      *redis.Client
	ctx      context.Context
	timeout  time.Duration
}

func CreateQueue(params *QueueConfig) (*Queue, error) {
	q := &Queue{
		Name:     params.Name,
		Address:  params.Address,
		Password: params.Password,
		DB:       params.DB,
		ctx:      context.TODO(),
	}
	if params.Timeout == 0 {
		q.timeout = time.Second * 60
	}

	q.rdb = redis.NewClient(&redis.Options{
		Addr:     params.Address,
		Password: params.Password,
		DB:       params.DB,
	})

	_, err := q.rdb.Ping(q.ctx).Result()
	if err != nil {
		return nil, err
	}

	return q, nil
}

func (this *Queue) Put(value string) error {
	_, err := this.rdb.LPush(this.ctx, this.Name, value).Result()
	return err
}

func (this *Queue) BLGet() (string, error) {
	result, err := this.rdb.BLPop(this.ctx, this.timeout, this.Name).Result()
	if err != nil {
		return "", err
	}
	return result[1], err
}
