package entity

import (
	"gopkg.in/redis.v5"
	"fmt"
)

type DB struct {
	client *redis.Client
}

func NewDB() *DB {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})
	return &DB{client}
}


func(db *DB) HealthCheck() error {
	_, err :=  db.client.Ping().Result()
	return err
}

func(db *DB) Save(redirect *Redirect) (res string, err error) {
	res, err = db.client.Set(createRedirectKey(redirect.From), redirect.To, 0).Result()
	return
}

func(db *DB) Get(key string) (*Redirect, error) {
	res, err := db.client.Get(key).Result()
	if err != nil {
		return nil, err
	}
	return &Redirect{key, res}, nil
}

func(db *DB) GetAll(offset, limit int) (res []*Redirect, err error) {
	if limit > 100 || limit < 0 {
		limit = 100
	}
	keys, _, err := db.client.Scan(uint64(offset), "source:*", int64(limit)).Result()
	var result []*Redirect
	for _, key := range keys {
		res, err := db.client.Get(key).Result()
		if err != nil {
			fmt.Println("Error reading from Redis", err)
			continue
		}
		result = append(result, &Redirect{key, res})
	}
	return result, nil
}

func(db *DB) Delete(key string) error {
	_, err := db.client.Del(createRedirectKey(key)).Result()
	return err
}

func createRedirectKey(from string) string {
	return fmt.Sprintf("source:%s", from)
}
