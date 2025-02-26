package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
)

var client *elasticsearch.Client

func init() {
	var err error
	client, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{Config.Elasticsearch.Host},
		Username:  Config.Elasticsearch.Username,
		Password:  Config.Elasticsearch.Password,
	})
	if err != nil {
		log.Fatalf("ES客户端创建失败: %v", err)
	}
}

func CreateIndex(indexName string, mappings *map[string]interface{}) error {
	res, err := client.Indices.Exists(
		[]string{indexName},
		client.Indices.Exists.WithContext(context.Background()),
	)
	if err != nil {
		return err
	}
	// 索引存在, 直接返回
	if res.StatusCode == 200 {
		return nil
	}

	mappingJSON, _ := json.Marshal(mappings)
	res, err = client.Indices.Create(indexName, client.Indices.Create.WithBody(bytes.NewReader(mappingJSON)))
	if err != nil {
		return err
	}
	if res.IsError() {
		return errors.New(res.String())
	}
	return nil
}

func IndexDocument(indexName string, docID string, doc []byte) error {
	res, err := client.Index(
		indexName,
		bytes.NewReader(doc),
		client.Index.WithDocumentID(docID),
		client.Index.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	if res.IsError() {
		return errors.New(res.String())
	}
	return nil
}

func SearchProducts(query string) ([]*map[string]interface{}, error) {
	indexName := "products"
	searchQuery := fmt.Sprintf(`{
       "query": {
          "multi_match": {
             "query": "%s",
             "fields": ["name", "description"]
          }
       }
    }`, query)
	queryReader := bytes.NewReader([]byte(searchQuery))
	res, err := client.Search(
		client.Search.WithIndex(indexName),
		client.Search.WithBody(queryReader),
		client.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, errors.New(res.String())
	}
	var result map[string]interface{}
	_ = json.NewDecoder(res.Body).Decode(&result)
	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	var sources []*map[string]interface{}
	for _, hit := range hits {
		doc := hit.(map[string]interface{})
		source := doc["_source"].(map[string]interface{})
		sources = append(sources, &source)
	}

	return sources, nil
}
