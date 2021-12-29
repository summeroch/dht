package main

import (
	"encoding/hex"
	"encoding/json"
	"github.com/summeroch/dht-spider/pkg/basic"
	"github.com/summeroch/dht-spider/pkg/config"
	"github.com/summeroch/dht-spider/pkg/dht"
	"github.com/summeroch/dht-spider/pkg/es"
	"log"
	"os"
	"time"
)

type file struct {
	Path   []interface{} `json:"path"`
	Length int           `json:"length"`
}

type bitTorrent struct {
	Time     time.Time `json:"@timestamp"`
	InfoHash string    `json:"infohash"`
	Name     string    `json:"name"`
	Files    []file    `json:"files,omitempty"`
	Length   int       `json:"length,omitempty"`
}

func main() {
	var c = config.GetConfig()
	f := "./spider.log"
	logFile, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("%s", err)
	}
	log.SetOutput(logFile)
	log.SetFlags(0)

	w := dht.NewWire(65536, 4096, 1024)
	go func() {
		for resp := range w.Response() {
			metadata, err := dht.Decode(resp.MetadataInfo)
			if err != nil {
				continue
			}

			info := metadata.(map[string]interface{})

			if _, ok := info["name"]; !ok {
				continue
			}

			bt := bitTorrent{
				Time:     time.Now(),
				InfoHash: hex.EncodeToString(resp.InfoHash),
				Name:     info["name"].(string),
			}

			if v, ok := info["files"]; ok {
				files := v.([]interface{})
				bt.Files = make([]file, len(files))

				for i, item := range files {
					f := item.(map[string]interface{})
					bt.Files[i] = file{
						Path:   f["path"].([]interface{}),
						Length: f["length"].(int),
					}
				}
			} else if _, ok := info["length"]; ok {
				bt.Length = info["length"].(int)
			}

			if data, err := json.Marshal(bt); err != nil {
				basic.CheckError(err)
			} else {
				if res, err := es.Write(c.Elasticsearch.Index, "_doc", data); err != nil {
					basic.CheckError(err)
				} else {
					log.Printf("%s\n%s\n", data, res)
				}
			}
		}
	}()
	go w.Run()

	AppConfig := dht.NewCrawlConfig()
	AppConfig.OnAnnouncePeer = func(infoHash, ip string, port int) {
		w.Request([]byte(infoHash), ip, port)
	}
	d := dht.New(AppConfig)

	d.Run()
}
