package es

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/summeroch/dht-spider/pkg/basic"
	"github.com/summeroch/dht-spider/pkg/config"
	"io"
	"log"
)

var (
	c = config.GetConfig()
	r map[string]interface{}
)

func Write(i, t string, data []byte) (string, error) {
	cfg := elasticsearch.Config{
		Addresses: c.Elasticsearch.Addresses,
		Username:  c.Elasticsearch.Username,
		Password:  c.Elasticsearch.Password,
	}

	NewClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		basic.CheckError(err)
	}

	req := esapi.CreateRequest{
		Index:        i,
		DocumentType: t,
		DocumentID:   basic.RandStr(20),
		Body:         bytes.NewReader(data),
	}

	res, err := req.Do(context.Background(), NewClient)
	if err != nil {
		basic.CheckError(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		basic.CheckError(err)
	}(res.Body)

	return res.String(), err
}

func Search(q string) map[string]interface{} {
	cfg := elasticsearch.Config{
		Addresses: c.Elasticsearch.Addresses,
		Username:  c.Elasticsearch.Username,
		Password:  c.Elasticsearch.Password,
	}

	NewClient, err := elasticsearch.NewClient(cfg)

	var buf bytes.Buffer

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"name": q,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Printf("Error encoding query: %s", err)
	}

	res, err := NewClient.Search(
		NewClient.Search.WithContext(context.Background()),
		NewClient.Search.WithIndex("spider"),
		NewClient.Search.WithBody(&buf),
		NewClient.Search.WithTrackTotalHits(true),
		NewClient.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	return r
}
