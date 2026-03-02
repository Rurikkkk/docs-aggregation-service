package aggregate

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"

	"docs-aggregation-service/internal/metrics"
)

type AggregateUsecase struct {
	filtersRepo FiltersRepo
	docsRepo    DocsRepo
	taskManager TaskManager
	metrics     *metrics.Metrics
}

func NewAggregateUsecase(
	filtersRepo FiltersRepo,
	docsRepo DocsRepo,
	taskManager TaskManager,
	metrics *metrics.Metrics,
) *AggregateUsecase {
	return &AggregateUsecase{
		filtersRepo: filtersRepo,
		docsRepo:    docsRepo,
		taskManager: taskManager,
		metrics:     metrics,
	}
}

func (au *AggregateUsecase) Run(startDate, endDate time.Time) (string, error) {
	taskID := uuid.New().String()
	err := au.taskManager.StartTask(taskID)
	if err != nil {
		log.Printf("[AggregationUsecase] task startup failed (ID %s...)", taskID[:6])
		return "", err
	}

	go func() {
		log.Printf("[Task] ID %s...: started", taskID[:6])

		filters, err := au.filtersRepo.GetFilters()
		if err != nil {
			log.Printf("[Task] ID %s...: reading filters failed", taskID[:6])
			termErr := au.taskManager.TerminateTask(taskID, err)
			if termErr != nil {
				log.Printf("[TASK] ID %s...: task finished with error, but terminating failed", taskID[:6])
			}
			log.Printf("[Task] ID %s...: finished with error", taskID[:6])
			au.metrics.TasksFailedTotal.Inc()
			au.metrics.TasksInProgress.Dec()
			return
		}

		docs, err := au.docsRepo.GetDocsByFilters(startDate, endDate, filters)
		if err != nil {
			log.Printf("[Task] ID %s...: getting documents failed", taskID[:6])
			termErr := au.taskManager.TerminateTask(taskID, err)
			if termErr != nil {
				log.Printf("[TASK] ID %s...: task finished with error, but terminating failed", taskID[:6])
			}
			log.Printf("[Task] ID %s...: task finished with error", taskID[:6])
			au.metrics.TasksFailedTotal.Inc()
			au.metrics.TasksInProgress.Dec()
			return
		}

		aggregation := make([][]string, 0)
		for _, doc := range docs {
			items := make([]string, len(doc.Items))
			for _, item := range doc.Items {
				items = append(items, item.Name)
			}
			docData := make([]string, 0)
			docData = append(
				docData,
				strings.TrimSpace(doc.FiscalDriveNumber),
				fmt.Sprint(doc.FiscalDocumentNumber),
				strings.Join(items, "|"),
			)
			aggregation = append(aggregation, docData)
		}

		complErr := au.taskManager.CompleteTask(taskID, aggregation)
		if complErr != nil {
			log.Printf("[TASK] ID %s...: task completed (docs: %d), but saving results failed", taskID[:6], len(aggregation))
			return
		}
		log.Printf("[Task] ID %s...: completed (docs: %d)", taskID[:6], len(aggregation))
		au.metrics.TasksCompletedTotal.Inc()
		au.metrics.TasksInProgress.Dec()
	}()

	au.metrics.TasksStartedTotal.Inc()
	au.metrics.TasksInProgress.Inc()
	return taskID, nil
}
