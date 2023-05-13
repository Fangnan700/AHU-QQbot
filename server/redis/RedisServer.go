package redis

import (
	"Kira-qbot/model"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"gopkg.in/yaml.v3"
	"os"
)

var (
	client *redis.Client
)

func init() {
	var config model.Config

	configFileName := "config/Config.yml"
	configFile, _ := os.Open(configFileName)
	decoder := yaml.NewDecoder(configFile)
	_ = decoder.Decode(&config)

	client = redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: config.RedisPass,
		DB:       0,
	})
}

func AddUser(user model.User) {
	jsonBytes, _ := json.Marshal(user)
	_ = client.Set(fmt.Sprintf("%d", user.UserId), string(jsonBytes), 0)
}

func DeleteUser(user model.User) {
	err := client.Del(fmt.Sprintf("%d", user.UserId)).Err()
	if err != nil {
		LogError(err)
	}
}

func GetUser(UserId int64) model.User {
	var user model.User
	jsonByte, _ := client.Get(fmt.Sprintf("%d", UserId)).Result()
	_ = json.Unmarshal([]byte(jsonByte), &user)
	return user
}
