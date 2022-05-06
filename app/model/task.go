package model

import "time"

const (
	TASK_STAT_PENDING = "pending"
	TASK_STAT_RUNNING = "running"
	TASK_STAT_SUCCESS = "success"
	TASK_STAT_FAILURE = "failure"
)

type Task struct {
	Id   string    `json:"id"`
	Name string    `json:"name"`
	Repo string    `json:"repo"`
	Pack string    `json:"pack"`
	User string    `json:"user"`
	Stat string    `json:"stat"`
	Hook string    `json:"hook"`
	File string    `json:"file"`
	Proc int       `json:"proc"`
	Time time.Time `json:"time"`
	Tags []string  `json:"tags"`
	Args []string  `json:"args"`
	Logs []Note    `json:"logs"`
}
