package collector

import (
	"context"
	"fmt"
	"log"
	"time"

	"backup-keeper/internal/domain"
	"backup-keeper/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBCollector struct {
	client    *mongo.Client
	database  string
	batchSize int
}

func NewMongoDBCollector(uri, database string) (domain.Collector, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	log.Println("MongoDB collector initialized and connection verified")

	return &MongoDBCollector{
		client:    client,
		database:  database,
		batchSize: 100000,
	}, nil
}

func (c *MongoDBCollector) Collect() (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db := c.client.Database(c.database)
	collections, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to get list collections: %v", err)
	}

	filePaths := make([]string, 0)

	for _, coll := range collections {
		log.Println("Processing collection: ", coll)

		tempFilePaths, err := dumpToTempFiles(ctx, db, coll, c.batchSize)

		if err != nil {
			log.Println("Error ", err)
			continue
		}

		filePaths = append(filePaths, tempFilePaths...)
		log.Println("Done collection: ", coll)
	}

	finalZipFile := fmt.Sprintf("dump_%s.zip", time.Now().Format("20060102_150405"))
	if err := utils.ZipFiles(filePaths, finalZipFile); err != nil {
		log.Println("Error ", err)
	}

	for _, tempFile := range filePaths {
		if err := utils.DeleteFile(tempFile); err != nil {
			log.Println("⚠️ Warning: Failed to delete temp file: " + tempFile + " - error: " + err.Error())
		}
	}

	return finalZipFile, nil
}

func dumpToTempFiles(ctx context.Context, db *mongo.Database, collName string, batchSize int) ([]string, error) {
	coll := db.Collection(collName)

	cursor, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var (
		docs     []bson.M
		docCount = 0
		chunkNum = 1
	)

	filePaths := make([]string, 0)

	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}

		docs = append(docs, doc)
		docCount++

		if docCount == batchSize {
			tempFile := fmt.Sprintf("%s_part_%d_size_%d", collName, chunkNum, docCount)
			if _, err := utils.WriteBatchToJson(docs, tempFile); err != nil {
				return nil, err
			}

			filePaths = append(filePaths, tempFile)

			docs = nil
			docCount = 0
			chunkNum++
		}
	}

	if docCount > 0 {
		tempFile := fmt.Sprintf("%s_part_%d_size_%d", collName, chunkNum, docCount)
		if _, err := utils.WriteBatchToJson(docs, tempFile); err != nil {
			return nil, err
		}

		filePaths = append(filePaths, tempFile)
	}

	return filePaths, nil
}

func (c *MongoDBCollector) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.client.Disconnect(ctx)
}
