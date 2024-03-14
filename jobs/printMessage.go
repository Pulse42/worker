package jobs

import "log"

type PrintMessageJobHandler struct {
}

func (h *PrintMessageJobHandler) Handle(job *Job) (Job, error) {
	log.Println(job.Data)
	return Job{}, nil
}
