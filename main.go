package main

import (
	"context"
	"encoding/json"
	"flag"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type url struct {
	UrlPath string `json:"url_path"`
}

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var m map[string]string

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func getUrl(c *gin.Context) {
	s := c.Param("str")

	log.Println(s)

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	log.Println(client.Ping(ctx))

	resp, err := client.Get(ctx, s).Result()
	if err != nil {
		log.Printf("key value failed to get: %v", err)
	}

	c.IndentedJSON(http.StatusOK, resp)
}

func postUrl(c *gin.Context) {
	ByteBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("request body failed to parse")
	}
	var data url
	err = json.Unmarshal(ByteBody, &data)
	if err != nil {
		log.Println("byte body failed to parse")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	log.Println(client.Ping(ctx))

	NewUrlPath := randStringBytes(8)

	err = client.Set(ctx, NewUrlPath, data.UrlPath, 0).Err()
	if err != nil {
		log.Printf("key value failed to add: %v", err)
	}

	NewUrlPath = strings.Join([]string{"http://127.0.0.1:8000", NewUrlPath}, "/")

	c.IndentedJSON(http.StatusOK, NewUrlPath)
}

func getUrlInMemory(c *gin.Context) {
	s := c.Param("str")

	log.Println(s)

	resp, err := m[s]
	if err {
		log.Printf("key value failed to get: %v", err)
	}

	c.IndentedJSON(http.StatusOK, resp)
}

func postUrlInMemory(c *gin.Context) {
	ByteBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("request body failed to parse")
	}
	var data url
	err = json.Unmarshal(ByteBody, &data)
	if err != nil {
		log.Println("byte body failed to parse")
	}

	NewUrlPath := randStringBytes(8)

	m[NewUrlPath] = data.UrlPath

	NewUrlPath = strings.Join([]string{"http://127.0.0.1:8000", NewUrlPath}, "/")

	c.IndentedJSON(http.StatusOK, NewUrlPath)
}

func main() {
	r := gin.Default()

	var nFlag = flag.Bool("d", false, "1 or redis by default, 0 or in memory optional")
	flag.Parse()

	if *nFlag {
		m = make(map[string]string)

		r.GET("/:str", getUrlInMemory)
		r.POST("/", postUrlInMemory)
	} else {
		r.GET("/:str", getUrl)
		r.POST("/", postUrl)
	}

	r.Run(":8000")
}
