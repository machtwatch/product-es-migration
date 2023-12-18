package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"product-es-migration/config"
	"product-es-migration/domain/product"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/dustin/go-humanize"
	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
)

var (
	numWorkers int = runtime.NumCPU()
	flushBytes int = 5e+6
)

var (
	countSuccessful uint64
)

func bulk(client *elasticsearch.Client, products []product.Product) {

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         config.ELASTICSEARCH_PRODUCT_INDEX, // The default index name
		Client:        client,                             // The Elasticsearch client
		NumWorkers:    numWorkers,                         // The number of worker goroutines
		FlushBytes:    int(flushBytes),                    // The flush threshold in bytes
		FlushInterval: 30 * time.Second,                   // The periodic flush interval
	})

	if err != nil {
		log.Fatalf("Error creating the indexer: %s", err)
	}

	for _, product := range products {
		data, err := json.Marshal(product)
		if err != nil {
			log.Fatalf("Cannot encode product %d: %s", product.Id, err)
		}

		// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
		//
		// Add an item to the BulkIndexer
		//
		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				// Action field configures the operation to perform (index, create, delete, update)
				Action: "index",

				// DocumentID is the (optional) document ID
				DocumentID: strconv.Itoa(int(product.Id)),

				// Body is an `io.Reader` with the payload
				Body: bytes.NewReader(data),

				// OnSuccess is called for each successful operation
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					atomic.AddUint64(&countSuccessful, 1)
				},

				// OnFailure is called for each failed operation
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("ERROR: %s", err)
					} else {
						log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			},
		)
		if err != nil {
			log.Fatalf("Unexpected error: %s", err)
		}

	}

	// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	// Close the indexer
	//
	if err := bi.Close(context.Background()); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}
	// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

	biStats := bi.Stats()

	// // Report the results: number of indexed docs, number of errors, duration, indexing rate
	// //
	log.Println(strings.Repeat("▔", 65))
	start := time.Now().UTC()
	dur := time.Since(start)

	if biStats.NumFailed > 0 {
		log.Fatalf(
			"Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			humanize.Comma(int64(biStats.NumFailed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	} else {
		log.Printf(
			"Sucessfuly indexed [%s] documents in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	}
}
