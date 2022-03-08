package filetransfer

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(addr, password string, db int) (*RedisStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	pong, err := client.Ping().Result()
	log.Printf("redis ping result: '%s'", pong)
	if err != nil {
		return nil, fmt.Errorf("problem connect to redis: %v", err)
	}
	return &RedisStore{client}, nil
}

func (r RedisStore) SaveUploadData(taskId string, data UploadData) {
	panic("implement me")
}

func (r RedisStore) GetUploadDataRemove(taskId string) *UploadData {
	panic("implement me")
}

func (r RedisStore) IsUploadTaskExist(taskId string) bool {
	panic("implement me")
}

func (r RedisStore) SaveDownloadData(taskId string, data DownloadData) {
	panic("implement me")
}

func (r RedisStore) GetDownloadDataRemove(taskId string) *DownloadData {
	panic("implement me")
}

func (r RedisStore) IsDownloadTaskExist(taskId string) bool {
	panic("implement me")
}
