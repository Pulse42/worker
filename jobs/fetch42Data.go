package jobs

import (
	"encoding/json"
	apiclient "github.com/Millefeuille42/42APIClient"
	"os"
)

type Fetch42DataJobHandlerData struct {
	Login string `json:"login"`
}

type Fetch42DataJobHandler struct {
	isInit bool
	client apiclient.APIClient
}

func (h *Fetch42DataJobHandler) init() {
	if h.isInit {
		return
	}

	h.client = apiclient.APIClient{
		Url:    os.Getenv("PULSE__WORKER__42_API_URL"),
		Uid:    os.Getenv("PULSE__WORKER__42_API_UID"),
		Secret: os.Getenv("PULSE__WORKER__42_API_SECRET"),
	}
}

func (h *Fetch42DataJobHandler) Handle(job *Job) (Job, error) {
	var data Fetch42DataJobHandlerData
	if err := json.Unmarshal([]byte(job.Data), &data); err != nil {
		return Job{}, err
	}

	h.init()
	if err := h.client.Auth(); err != nil {
		return Job{}, err
	}

	user, err := h.client.GetUser(data.Login)
	if err != nil {
		return Job{}, err
	}

	userAsBytes, err := json.Marshal(user)
	if err != nil {
		return Job{}, err
	}
	return Job{
		Type: Set42DataJobType,
		Data: string(userAsBytes),
	}, nil
}
