package httpserver

import (
	"io"
	"time"

	"docs-aggregation-service/internal/domains/taskdomain"
)

type AggregateUsecase interface {
	Run(startDate, endDate time.Time, filtersReader io.ReadSeeker) (string, error)
}

type GetStatusUsecase interface {
	Run(taskID string) ([]taskdomain.Task, error)
}

type GetResultUsecase interface {
	Run(taskID string) (string, error)
}
