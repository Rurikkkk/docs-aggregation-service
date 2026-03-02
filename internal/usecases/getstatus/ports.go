package getstatus

import "docs-aggregation-service/internal/domains/taskdomain"

type TaskManager interface {
	GetStatus(taskID string) ([]taskdomain.Task, error)
}
