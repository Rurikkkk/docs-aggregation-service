package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"docs-aggregation-service/internal/adapters/docsrepo"
	"docs-aggregation-service/internal/adapters/filtersparser"
	"docs-aggregation-service/internal/adapters/taskmanager"
	"docs-aggregation-service/internal/metrics"
	"docs-aggregation-service/internal/ui/httpserver"
	"docs-aggregation-service/internal/usecases/aggregate"
	"docs-aggregation-service/internal/usecases/getresult"
	"docs-aggregation-service/internal/usecases/getstatus"
)

type Context struct {
	MongoURL            string
	MongoDBName         string
	DocsCollectionName  string
	TasksCollectionName string
	AggregationDirpath  string
	ServerAddr          string
}

var (
	once     sync.Once
	instance *Context
)

func NewContext() *Context {
	once.Do(func() {
		instance = &Context{
			MongoURL:            getEnvOrDefault("MONGO_URL", "mongodb://localhost:27017"),
			MongoDBName:         getEnvOrDefault("MONGO_DB_NAME", "docs-aggregation-service"),
			DocsCollectionName:  getEnvOrDefault("DOCS_COLLECTION_NAME", "docs"),
			TasksCollectionName: getEnvOrDefault("TASKS_COLLECTION_NAME", "tasks"),
			AggregationDirpath:  getEnvOrDefault("AGGREGATION_DIRPATH", "aggregations"),
			ServerAddr:          getEnvOrDefault("SERVER_ADDR", ":8080"),
		}
	})
	return instance
}

func getEnvOrDefault(key, defaultEnv string) string {
	env := os.Getenv(key)
	if env == "" {
		return defaultEnv
	}
	return env
}

func (c *Context) HTTPServer() *http.Server {
	metrics := metrics.NewMetrics()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(c.MongoURL))
	if err != nil {
		log.Printf("[APP] Connecting to MongoDB failed (URL: %s)", c.MongoURL)
		return nil
	}

	filtersParser := filtersparser.NewFiltersParser()
	docsRepo := docsrepo.NewDocsRepo(client, c.MongoDBName, c.DocsCollectionName)
	taskManager := taskmanager.NewTaskManager(client, c.MongoDBName, c.TasksCollectionName)
	aggregateUsecase := aggregate.NewAggregateUsecase(filtersParser, docsRepo, taskManager, metrics)
	getStatusUsecase := getstatus.NewGetStatusUsecase(taskManager)
	getResultUsecase := getresult.NewGetResultUsecase(taskManager)

	server := httpserver.NewHTTPServer(aggregateUsecase, getStatusUsecase, getResultUsecase, metrics)
	server.Addr = c.ServerAddr
	return server
}
