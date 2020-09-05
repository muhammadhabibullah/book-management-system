package elasticsearch

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"

	"book-management-system/configs"
)

var (
	es   *elasticsearch.Client
	once sync.Once
)

// Init returns elastic search client
func Init() *elasticsearch.Client {
	once.Do(func() {
		cfg := configs.GetConfig().ElasticSearch

		esCfg := elasticsearch.Config{
			Addresses: []string{
				cfg.Address,
			},
		}
		if cfg.IsAuth {
			esCfg.Username = cfg.Username
			esCfg.Password = cfg.Password
		}

		var err error
		es, err = elasticsearch.NewClient(esCfg)
		if err != nil {
			log.Fatal(err)
		}

		res, err := es.Info()
		if err != nil {
			log.Fatalf("Error getting elasticsearch info: %s", err)
		}
		defer func() {
			_ = res.Body.Close()
		}()
		if res.IsError() {
			log.Fatalf("Error elasticsearch info response: %s", res.String())
		}

		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		}
	})

	return es
}
