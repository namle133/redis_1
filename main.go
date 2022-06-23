package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
)

func main() {
	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379",
	// 	Password: "",
	// 	DB:       1,
	// })

	// err := rdb.Set("name", "Nam", 0).Err()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// val, e := rdb.Get("name").Result()
	// if e != nil {
	// 	fmt.Println(e)
	// 	return
	// }
	// fmt.Println(val)
	redisCluster()
}

type redisCluterClient struct {
	c *redis.ClusterClient
}

var client = &redisCluterClient{}

func initCl(addr []string) (*redisCluterClient, error) {

	c := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: addr,
	})

	if err := c.Ping().Err(); err != nil {
		return nil, err
	}

	client.c = c
	return client, nil
}

func (client *redisCluterClient) setKey(key string, value interface{}, expiration time.Duration) error {

	json, err := json.Marshal(value)

	if err != nil {
		return err
	}

	if err = client.c.Set(key, json, expiration).Err(); err != nil {
		return err
	}

	return nil
}

func (client *redisCluterClient) getKey(key string, src interface{}) error {

	val, err := client.c.Get(key).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(val), &src)
	if err != nil {
		return err
	}
	return nil
}

type valueEx struct {
	Name  string
	Email string
}

func redisCluster() {

	addr := []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"}
	redisCluterClient, err := initCl(addr)
	if err != nil {
		fmt.Println("Error connect: ", err)
		return
	}
	key1 := "myKeyCluster"
	value1 := &valueEx{Name: "congpv", Email: "congpv@lozi.vn"}
	err = redisCluterClient.setKey(key1, value1, time.Minute*1)
	if err != nil {
		fmt.Printf("Error: %v", err.Error())
		return
	}
	value2 := &valueEx{}
	err = redisCluterClient.getKey(key1, value2)
	if err != nil {
		fmt.Printf("Error: %v", err.Error())
		return
	}

	fmt.Println(value2.Email)
	fmt.Println(value2.Name)
}
