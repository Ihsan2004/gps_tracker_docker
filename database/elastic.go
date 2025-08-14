package database

import (
	"log"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// ConnectElastic Docker ortamı için güncellendi
func ConnectElastic() *elasticsearch.Client {
	// Environment variable'dan URL al, yoksa localhost kullan
	esURL := getEnv("ELASTICSEARCH_URL", "http://localhost:9200")

	cfg := elasticsearch.Config{
		Addresses:     []string{esURL},
		RetryOnStatus: []int{502, 503, 504, 429},
		RetryBackoff: func(i int) time.Duration {
			return time.Duration(i) * 100 * time.Millisecond
		},
		MaxRetries: 5,
	}

	// Retry logic ile bağlantı kur
	var Es *elasticsearch.Client
	var err error

	for i := 0; i < 5; i++ {
		Es, err = elasticsearch.NewClient(cfg)
		if err == nil {
			// Bağlantıyı test et
			info, err := Es.Info()
			if err == nil {
				info.Body.Close()
				log.Printf("Connected to Elasticsearch successfully at: %s", esURL)
				return Es
			}
		}

		log.Printf("Elasticsearch connection attempt %d failed: %v", i+1, err)
		time.Sleep(time.Duration(i+1) * 2 * time.Second)
	}

	log.Fatal("Elasticsearch connection failed after retries:", err)
	return nil
}
