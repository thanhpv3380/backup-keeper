package collector

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"backup-keeper/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBCollector struct {
	client      *mongo.Client
	database    string
	collections []string
}

func NewMongoDBCollector(uri, database string, collections []string) (domain.Collector, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	log.Println("MongoDB collector initialized")

	return &MongoDBCollector{
		client:      client,
		database:    database,
		collections: collections,
	}, nil
}

func (c *MongoDBCollector) Collect() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a buffer to hold all collections data
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)

	// Create a map to hold all collections data
	collectionsData := make(map[string]interface{})

	for _, collection := range c.collections {
		// Get all documents from collection
		cursor, err := c.client.Database(c.database).Collection(collection).Find(ctx, bson.M{})
		if err != nil {
			return nil, fmt.Errorf("failed to query collection %s: %v", collection, err)
		}

		var results []bson.M
		if err = cursor.All(ctx, &results); err != nil {
			return nil, fmt.Errorf("failed to decode collection %s: %v", collection, err)
		}

		// Convert BSON to JSON-serializable format
		var jsonResults []map[string]interface{}
		for _, doc := range results {
			jsonResults = append(jsonResults, doc)
		}

		collectionsData[collection] = jsonResults
	}

	// Encode all data to JSON
	if err := encoder.Encode(collectionsData); err != nil {
		return nil, fmt.Errorf("failed to encode collections data: %v", err)
	}

	return buffer.Bytes(), nil
}

func (c *MongoDBCollector) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.client.Disconnect(ctx)
}
