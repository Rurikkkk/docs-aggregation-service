package getresult

import (
	"errors"
	"log"

	"docs-aggregation-service/internal/domains/taskdomain"
)

type GetResultUsecase struct {
	taskManager TaskManager
}

func NewGetResultUsecase(taskManager TaskManager) *GetResultUsecase {
	return &GetResultUsecase{taskManager: taskManager}
}

func (gru *GetResultUsecase) Run(taskID string) (string, error) {
	task, err := gru.taskManager.GetStatus(taskID)
	if err != nil {
		log.Printf("[GetResultUsecase] Getting task status failed (ID %s...)", taskID[:6])
		return "", err
	}
	if task[0].Status != taskdomain.StatusCompleted {
		log.Printf("[GetResultUsecase] Task is not completed (ID %s...)", taskID[:6])
		return "", errors.New("not completed")
	}
	return task[0].Filepath, nil
}
