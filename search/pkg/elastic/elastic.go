package elastic

import (
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

func NewElasticClient(addresses []string) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating the client: %w", err)
	}

	res, err := es.Info()
	if err != nil {
		return nil, fmt.Errorf("error getting info response: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error: %s", res.String())
	}

	_, _ = es.Indices.Create("users", es.Indices.Create.WithBody(strings.NewReader(`{
		"mappings": {
		"properties": {
			"username": { "type": "text", "fields": { "keyword": { "type": "keyword" } } },
			"bio": { "type": "text", "analyzer": "russian" }
		}
		}
	}`)))

	_, _ = es.Indices.Create("tweets", es.Indices.Create.WithBody(strings.NewReader(`{
		"mappings": {
		"properties": {
			"content": { "type": "text", "analyzer": "russian" },
			"username": { "type": "text", "fields": { "keyword": { "type": "keyword" } } },
			"user_id": { "type": "integer" }
		}
		}
	}`)))

	return es, nil
}
