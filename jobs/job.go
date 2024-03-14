package jobs

import "errors"

type Type string

const (
	PrintMessageJobType Type = "print_message"
	Fetch42DataJobType  Type = "fetch_42_data"
	Set42DataJobType    Type = "set_42_data"
)

type Job struct {
	Type Type   `json:"type"`
	Data string `json:"data"`
}

type Handler interface {
	Handle(*Job) (Job, error)
}

var handlerMap = map[Type]Handler{
	PrintMessageJobType: &PrintMessageJobHandler{},
	Fetch42DataJobType:  &Fetch42DataJobHandler{},
}

func Handle(job *Job) (Job, error) {
	if handler, ok := handlerMap[job.Type]; ok {
		return handler.Handle(job)
	}
	return Job{}, errors.New(string("unknown job type: " + job.Type))
}
