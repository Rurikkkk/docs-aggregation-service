package docsrepo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"docs-aggregation-service/internal/usecases/aggregate"
)

type DocsRepo struct {
	docsCollection *mongo.Collection
}

func NewDocsRepo(client *mongo.Client, dbName, collectionName string) *DocsRepo {
	return &DocsRepo{docsCollection: client.Database(dbName).Collection(collectionName)}
}

func (dr *DocsRepo) GetDocsByFilters(
	startDate,
	endDate time.Time,
	fiscalDriveNums []string,
) ([]aggregate.DocumentFields, error) {
	filter := bson.M{
		"doc.dateTime":          bson.M{"$gte": startDate, "$lte": endDate},
		"doc.fiscalDriveNumber": bson.M{"$in": fiscalDriveNums},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := dr.docsCollection.Find(ctx, filter)
	if err != nil {
		log.Printf("[DocsRepo] Reading docs from MongoDB repo failed: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []aggregate.DocumentFields
	for cursor.Next(ctx) {
		var doc aggregate.Document
		err := cursor.Decode(&doc)
		if err != nil {
			log.Printf("[DocsRepo] Decoding doc failed: %v", err)
			return nil, err
		}
		docs = append(docs, doc.Doc)
	}
	return docs, nil
}
