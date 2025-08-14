package service

import (
	"GpsTracker2/models"
	"GpsTracker2/repository"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"strings"
)

const ModelCategory = 1       // Category for searching by model name
const BrandCategory = 2       // Category for searching by brand name
const DescriptionCategory = 3 // Category for searching by description

// Create Device
func CreateDevice(userID int, device *models.Device, es *elasticsearch.Client) error {
	// ✅ Ensure user exists before creating device
	exists, err := repository.CheckUserExist(userID)
	if err != nil {
		return err
	}
	if exists == 0 {
		return errors.New("User not found")
	}

	// ✅ Save device to the database
	device.UserId = userID
	if err := repository.CreateDevice(device); err != nil {
		return err
	}

	// Index in Elasticsearch
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(device); err != nil {
		return err
	}
	res, err := es.Index("devices", &buf)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	return nil
}

func GetAllDevices(page int, size int, es *elasticsearch.Client, query string, categories []int) ([]models.Device, int, error) {
	var shouldQueries []map[string]interface{}

	if len(categories) == 0 {
		shouldQueries = append(shouldQueries, map[string]interface{}{
			"wildcard": map[string]interface{}{
				"model": map[string]interface{}{
					"value": strings.ToLower("*" + query + "*"),
				},
			},
		})
		shouldQueries = append(shouldQueries, map[string]interface{}{
			"wildcard": map[string]interface{}{
				"brand": map[string]interface{}{
					"value": strings.ToLower("*" + query + "*"),
				},
			},
		})

		shouldQueries = append(shouldQueries, map[string]interface{}{
			"wildcard": map[string]interface{}{
				"description": map[string]interface{}{
					"value": strings.ToLower("*" + query + "*"),
				},
			},
		})
	}

	for _, cat := range categories {
		switch cat {
		case ModelCategory:
			shouldQueries = append(shouldQueries, map[string]interface{}{
				"wildcard": map[string]interface{}{
					"model": map[string]interface{}{
						"value": strings.ToLower("*" + query + "*"),
					},
				},
			})
		case BrandCategory:
			shouldQueries = append(shouldQueries, map[string]interface{}{
				"wildcard": map[string]interface{}{
					"brand": map[string]interface{}{
						"value": strings.ToLower("*" + query + "*"),
					},
				},
			})
		case DescriptionCategory:
			shouldQueries = append(shouldQueries, map[string]interface{}{
				"wildcard": map[string]interface{}{
					"description": map[string]interface{}{
						"value": strings.ToLower("*" + query + "*"),
					},
				},
			})
		default:
			return nil, 0, errors.New("invalid category")
		}
	}

	esQuery := map[string]interface{}{
		"from": (page - 1) * size,
		"size": size,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should":               shouldQueries,
				"minimum_should_match": 1,
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(esQuery); err != nil {
		return nil, 0, err
	}

	res, err := es.Search(
		es.Search.WithIndex("devices"),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Println("Elasticsearch search error:", res.String())
		return nil, 0, errors.New("failed to search devices in Elasticsearch")
	}

	var esResult map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&esResult); err != nil {
		return nil, 0, err
	}
	hits := esResult["hits"].(map[string]interface{})["hits"].([]interface{})
	if len(hits) == 0 {
		return nil, 0, nil
	}

	var devices []models.Device
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]
		sourceJSON, _ := json.Marshal(source)
		var device models.Device
		if err := json.Unmarshal(sourceJSON, &device); err != nil {
			return nil, 0, err
		}
		devices = append(devices, device)
	}

	totalHits := int(esResult["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
	return devices, totalHits / size, nil
}

// Get Device by ID
func GetDeviceByID(deviceID int, es *elasticsearch.Client) (models.Device, error) {
	// 1) Search with Elasticsearch
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"device_id": deviceID,
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return models.Device{}, err
	}

	res, err := es.Search(
		es.Search.WithIndex("devices"),
		es.Search.WithBody(&buf),
	)
	if err != nil {
		return models.Device{}, err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatal("Elasticsearch search error:", res.String(), "returning from my")
		return repository.GetDeviceByID(deviceID)
	}

	// 2) Analyze response
	var esResult map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&esResult); err != nil {
		return models.Device{}, err
	}

	hits := esResult["hits"].(map[string]interface{})["hits"].([]interface{})
	if len(hits) == 0 {
		log.Println("Device not found in Elasticsearch")
		return repository.GetDeviceByID(deviceID)
	}

	// 3) Return first result
	source := hits[0].(map[string]interface{})["_source"]
	sourceJSON, _ := json.Marshal(source)

	var device models.Device
	if err := json.Unmarshal(sourceJSON, &device); err != nil {
		return models.Device{}, err
	}

	return device, nil
}

// Update Device
func UpdateDevice(deviceID int, updatedDevice models.Device, es *elasticsearch.Client) error {
	// ✅ Ensure device exists before updating
	exists, err := repository.CheckDeviceExist(deviceID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Device not found")
	}

	// Delete the old device from Elasticsearch and add the updated one
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"device_id": deviceID,
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return err
	}

	res, err := es.DeleteByQuery(
		[]string{"devices"},
		&buf,
	)

	if err != nil {
		log.Fatal("Elasticsearch delete error:", res.String())
	}

	defer res.Body.Close()
	if res.IsError() {
		log.Fatal("Elasticsearch delete error:", res.String())
		return errors.New("Failed to delete old device in Elasticsearch")
	}

	// Index the updated device in Elasticsearch
	var updatedBuf bytes.Buffer
	if err := json.NewEncoder(&updatedBuf).Encode(updatedDevice); err != nil {
		log.Fatal("Elasticsearch index error:", err)
	}
	res, err = es.Index("devices", &updatedBuf)
	if err != nil {
		log.Fatal("Elasticsearch index error:", err)
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		log.Fatal("Elasticsearch index error:", res.String())
		return errors.New("Failed to index updated device in Elasticsearch")
	}

	// ✅ Update device in the database
	return repository.UpdateDevice(deviceID, updatedDevice)
}

// Delete Device
func DeleteDevice(deviceID int, es *elasticsearch.Client) error {
	// ✅ Ensure device exists before deleting
	exists, err := repository.CheckDeviceExist(deviceID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Device not found")
	}

	// Delete from Elasticsearch
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"device_id": deviceID,
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Elasticsearch delete error: %s", err)
	}

	res, err := es.DeleteByQuery(
		[]string{"devices"},
		&buf,
	)

	if err != nil {
		log.Fatalf("Elasticsearch delete error: %s", err)
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		log.Fatalf("Elasticsearch delete error: %s", res.String())
		return errors.New("Failed to delete device in Elasticsearch")
	}

	// ✅ Delete device from the database
	return repository.DeleteDevice(deviceID)
}
