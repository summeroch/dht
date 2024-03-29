package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/summeroch/dht-spider/pkg/es"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	f, _ := os.Create("./search.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	router := gin.New()
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("{\"@timestamp\":\"%s\",\"remoteAddr\":\"%s\",\"request_method\":\"%s\",\"request_uri\":\"%s\",\"protocol\":\"%s\",\"status\":\"%d\",\"requestTime\":\"%s\",\"httpUserAgent\":\"%s\",\"Errormessage\":\"%s\"}\n",
			param.TimeStamp.Format(time.RFC3339),
			param.ClientIP,
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())

	router.GET("/search", func(c *gin.Context) {
		name := c.Query("name")
		res := es.Search(name)
		hit := int(res["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
		took := int(res["took"].(float64))
		for num, data := range res["hits"].(map[string]interface{})["hits"].([]interface{}) {
			id := data.(map[string]interface{})["_id"].(string)
			t := data.(map[string]interface{})["_source"].(map[string]interface{})["@timestamp"].(string)
			n := data.(map[string]interface{})["_source"].(map[string]interface{})["name"].(string)
			hash := data.(map[string]interface{})["_source"].(map[string]interface{})["infohash"].(string)
			data := gin.H{
				"id":   id,
				"hits": hit,
				"took": took,
				"time": t,
				"list": num,
				"name": n,
				"hash": hash,
			}
			c.JSONP(http.StatusOK, data)
		}
	})

	err := router.Run(":8080")
	if err != nil {
		log.Fatalln(err)
	}
}
