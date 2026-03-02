package taskdomain

type Task struct {
	ID       string     `bson:"_id" json:"id"`
	Status   TaskStatus `bson:"status" json:"status"`
	Error    error      `bson:"error,omitempty" json:"error,omitempty"`
	Filepath string     `bson:"filepath,omitempty" json:"filepath,omitempty"`
}

type TaskStatus string

const (
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
	StatusError     TaskStatus = "error"
)
