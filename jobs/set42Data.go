package jobs

import "log"

type Set42DAtaJobHandler struct {
}

func (h *Set42DAtaJobHandler) Handle(job *Job) ([]byte, error) {
	log.Println(job.Data)
	return nil, nil
}
