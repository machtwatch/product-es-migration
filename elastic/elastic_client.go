package elastic

import (
	"product-es-migration/config"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
)

func connect() (*elasticsearch.Client, error) {

	// uncomment logger to print out all api calls
	cfg := elasticsearch.Config{
		Addresses: []string{config.ELASTICSEARCH_CLUSTER},
		Username:  config.ELASTICSEARCH_USERNAME,
		Password:  config.ELASTICSEARCH_PASSWORD,
		// Logger:    &estransport.ColorLogger{Output: os.Stdout, EnableRequestBody: true, EnableResponseBody: true}, // print out http client
	}

	es, err := elasticsearch.NewClient(cfg)

	return es, err
}
