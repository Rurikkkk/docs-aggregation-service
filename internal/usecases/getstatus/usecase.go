package getstatus

import (
	"log"

	"docs-aggregation-service/internal/domains/taskdomain"
)

type GetStatusUsecase struct {
	taskManager TaskManager
}

func NewGetStatusUsecase(taskManager TaskManager) *GetStatusUsecase {
	return &GetStatusUsecase{taskManager: taskManager}
}

func (gsu *GetStatusUsecase) Run(taskID string) ([]taskdomain.Task, error) {
	tasks, err := gsu.taskManager.GetStatus(taskID)
	if err != nil {
		log.Printf("[GetStatusUsecase] Getting task statuses error (ID %s...)", taskID[:6])
		return []taskdomain.Task{}, err
	}
	return tasks, nil
}
