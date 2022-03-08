package filetransfer

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"strings"
	"time"
)

const uploadSuffix = "upload"
const downloadSuffix = "download"

type redisStore struct {
	client *redis.Client
}

func NewRedisStore(addr, password string, db int) (DataStore, error) {
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
	return &redisStore{client}, nil
}

func (r redisStore) SaveUploadData(taskId string, data UploadData) {
	if taskId == "" {
		return
	}
	key := r.createUploadKey(taskId)
	uploadJSONData := r.data2Json(data)
	r.client.Set(key, uploadJSONData, 10*time.Minute)
}

func (r redisStore) GetUploadDataRemove(taskId string) *UploadData {
	key := r.createUploadKey(taskId)
	uploadJSONData, err := r.client.Get(key).Result()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		log.Printf("problem get data: %v", err)
		return nil
	}
	var uploadData UploadData
	err = json.NewDecoder(strings.NewReader(uploadJSONData)).Decode(&uploadData)
	if err != nil {
		log.Printf("problem decode data: %v", err)
	}
	r.client.Del(key)
	return &uploadData
}

func (r redisStore) IsUploadTaskExist(taskId string) bool {
	_, err := r.client.Get(r.createUploadKey(taskId)).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		log.Printf("problem get data: %v", err)
		return false
	}
	return true
}

func (r redisStore) SaveDownloadData(taskId string, data DownloadData) {
	if taskId == "" {
		return
	}
	key := r.createDownloadKey(taskId)
	downloadJSONData := r.data2Json(data)
	r.client.Set(key, downloadJSONData, 10*time.Minute)
}

func (r redisStore) GetDownloadDataRemove(taskId string) *DownloadData {
	key := r.createDownloadKey(taskId)
	downloadJSONData, err := r.client.Get(key).Result()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		log.Printf("problem get data: %v", err)
		return nil
	}
	var downloadData DownloadData
	err = json.NewDecoder(strings.NewReader(downloadJSONData)).Decode(&downloadData)
	if err != nil {
		log.Printf("problem decode data: %v", err)
	}
	r.client.Del(key)
	return &downloadData
}

func (r redisStore) IsDownloadTaskExist(taskId string) bool {
	_, err := r.client.Get(r.createDownloadKey(taskId)).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		log.Printf("problem get data: %v", err)
		return false
	}
	return true
}

// 合成上传任务的key
func (redisStore) createUploadKey(taskId string) string {
	return fmt.Sprintf("%s:%s", uploadSuffix, taskId)
}

// 合成下载任务的key
func (redisStore) createDownloadKey(taskId string) string {
	return fmt.Sprintf("%s:%s", downloadSuffix, taskId)
}

// po转换成json
func (r redisStore) data2Json(data interface{}) string {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("problem encode upload data to json")
	}
	return string(bytes)
}
