package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
)

type ValueX struct {
	Name  string
	Email string
}

func main() {
	rdb, err := initC("localhost:6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	key1 := "MyCluster"
	value1 := &ValueX{Name: "Nam", Email: "nam@gmail.com"}

	err = rdb.setKey(key1, value1, time.Minute*2)
	if err != nil {
		fmt.Println(err)
		return
	}
	value2 := &ValueX{}
	err = rdb.getKey(key1, value2)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("MyCluster")
	fmt.Println(value2.Email)
	// redisCluster()
}

type redisClient struct {
	c *redis.Client
}

type redisCluterClient struct {
	c *redis.ClusterClient
}

var client = &redisCluterClient{}
var cl = &redisClient{}

func initCl(addr []string) (*redisCluterClient, error) {
	c := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: addr,
	})

	err := c.Ping().Err()
	if err != nil {
		return nil, err
	}

	client.c = c

	return client, nil

}

func initC(addr string) (*redisClient, error) {
	c := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	err := c.Ping().Err()
	if err != nil {
		return nil, err
	}

	cl.c = c

	return cl, nil

}

func (cl *redisClient) setKey(key string, value interface{}, expiration time.Duration) error {
	json, err := json.Marshal(value)

	if err != nil {
		return err
	}

	err = cl.c.Set(key, json, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (client *redisCluterClient) setKey(key string, value interface{}, expiration time.Duration) error {
	json, err := json.Marshal(value)

	if err != nil {
		return err
	}

	err = client.c.Set(key, json, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (cl *redisClient) getKey(key string, src interface{}) error {
	val, err := cl.c.Get(key).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(val), &src)
	if err != nil {
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

func redisCluster() {
	addr := []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"}

	ReCl, err := initCl(addr)

	if err != nil {
		fmt.Println(err)
		return
	}
	key1 := "MyCluster"
	value1 := &ValueX{Name: "Nam", Email: "nam@gmail.com"}

	err = ReCl.setKey(key1, value1, time.Minute*2)
	if err != nil {
		fmt.Println(err)
		return
	}
	value2 := &ValueX{}
	err = ReCl.getKey(key1, value2)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("MyCluster")
	fmt.Println(value2.Email)

}
