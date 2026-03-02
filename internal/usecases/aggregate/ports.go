package aggregate

import "time"

type FiltersRepo interface {
	GetFilters() ([]string, error)
}

type DocsRepo interface {
	GetDocsByFilters(startDate, endDate time.Time, fiscalDriveNums []string) ([]DocumentFields, error)
}

type TaskManager interface {
	StartTask(taskID string) error
	CompleteTask(taskID string, aggregation [][]string) error
	TerminateTask(taskID string, err error) error
}
