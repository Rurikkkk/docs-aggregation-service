package taskmanager

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"docs-aggregation-service/internal/domains/taskdomain"
)

type TaskManager struct {
	tasksCollection *mongo.Collection
}

func NewTaskManager(client *mongo.Client, dbName, collectionName string) *TaskManager {
	return &TaskManager{tasksCollection: client.Database(dbName).Collection(collectionName)}
}

func (tm *TaskManager) StartTask(taskID string) error {
	task := taskdomain.Task{ID: taskID, Status: taskdomain.StatusError}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := tm.tasksCollection.InsertOne(ctx, task)
	if err != nil {
		log.Printf("[TaskManager] Inserting task failed (ID %s...): %v", taskID[:6], err)
		return err
	}
	log.Printf("[TaskManager] Inserting task completed (ID %s...)", taskID[:6])
	return nil
}

func (tm *TaskManager) CompleteTask(taskID string, aggregation [][]string) error {
	err := os.MkdirAll("aggregations", 0o755)
	if err != nil {
		log.Printf("[TaskManager] Creating aggregation files directory failed (ID %s...): %v", taskID[:6], err)
		termErr := tm.TerminateTask(taskID, err)
		if termErr != nil {
			log.Printf("[TaskManager] Task terminating after completing error failed (ID %s...)", taskID[:6])
		}
		return err
	}

	filepath := filepath.Join("aggregations", fmt.Sprintf("aggregation_%s.csv", taskID))
	file, err := os.Create(filepath)
	if err != nil {
		log.Printf("[TaskManager] Creating aggregation CSV file failed (ID %s...): %v", taskID[:6], err)
		termErr := tm.TerminateTask(taskID, err)
		if termErr != nil {
			log.Printf("[TaskManager] Task terminating after completing error failed (ID %s...)", taskID[:6])
		}
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"fiscalDriveNumber", "fiscalDocumentNumber", "items"})
	if err != nil {
		log.Printf("[TaskManager] Writing CSV header failed (ID %s...): %v", taskID[:6], err)
		termErr := tm.TerminateTask(taskID, err)
		if termErr != nil {
			log.Printf("[TaskManager] Task terminating after completing error failed (ID %s...)", taskID[:6])
		}
		return err
	}
	for _, line := range aggregation {
		err := writer.Write(line)
		if err != nil {
			log.Printf("[TaskManager] Writing aggregation data failed (ID %s...): %v", taskID[:6], err)
			termErr := tm.TerminateTask(taskID, err)
			if termErr != nil {
				log.Printf("[TaskManager] Task terminating after completing error failed (ID %s...)", taskID[:6])
			}
			return err
		}
	}

	filter := bson.M{"_id": taskID}
	update := bson.M{
		"$set": bson.M{
			"status":   taskdomain.StatusCompleted,
			"filepath": filepath,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = tm.tasksCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("[TaskManager] Updating task failed (ID %s...): %v", taskID[:6], err)
		return err
	}
	log.Printf("[TaskManager] Completing task completed (ID %s...)", taskID[:6])
	return nil
}

func (tm *TaskManager) TerminateTask(taskID string, err error) error {
	filter := bson.M{"_id": taskID}
	update := bson.M{
		"$set": bson.M{
			"status": taskdomain.StatusError,
			"error":  err,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, updErr := tm.tasksCollection.UpdateOne(ctx, filter, update)
	if updErr != nil {
		log.Printf("[TaskManager] Updating task failed (ID %s...): %v", taskID[:6], err)
		return err
	}
	log.Printf("[TaskManager] Terminating task completed (ID %s...)", taskID[:6])
	return nil
}

func (tm *TaskManager) GetStatus(taskID string) ([]taskdomain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if taskID != "" {
		var task taskdomain.Task
		filter := bson.M{"_id": taskID}
		err := tm.tasksCollection.FindOne(ctx, filter).Decode(&task)
		if err != nil {
			log.Printf("[TaskManager] Finding and decoding task failed (ID %s...): %v", taskID[:6], err)
			return []taskdomain.Task{}, err
		}
		return []taskdomain.Task{task}, nil
	}

	cursor, err := tm.tasksCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("[TaskManager] Getting tasks failed: %v", err)
		return []taskdomain.Task{}, err
	}
	defer cursor.Close(ctx)

	var tasks []taskdomain.Task
	for cursor.Next(ctx) {
		var task taskdomain.Task
		err := cursor.Decode(&task)
		if err != nil {
			log.Printf("[TaskManager] Decoding task status failed: %v", err)
			return []taskdomain.Task{}, err
		}
		tasks = append(tasks, task)
	}
	if len(tasks) == 0 {
		log.Printf("[TaskManager] No tasks founded")
	}
	return tasks, nil
}
