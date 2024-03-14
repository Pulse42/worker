package main

import (
	"encoding/json"
	"log"
	"worker/consumerProducers"
	"worker/jobs"
)

const consumerProducerType = consumerProducers.RedisConsumerProducerType

var consumerParameters = `{
	"Addr": "localhost:7000",
	"Password": "",
	"DB": 0,
	"Channel": "work"
}`

func main() {
	consumerProducer := consumerProducers.Factory(consumerProducerType)
	err := consumerProducer.Subscribe(consumerParameters)
	if err != nil {
		log.Fatal(err)
	}
	defer consumerProducer.Close()

	messageChannel := consumerProducer.MessageChannel()
	errorChannel := consumerProducer.ErrorChannel()

	var job jobs.Job
	for {
		select {
		case msg := <-messageChannel:
			if err = json.Unmarshal([]byte(msg), &job); err != nil {
				log.Println(err)
				continue
			}
			job, err = jobs.Handle(&job)
			if err != nil {
				log.Println(err)
				continue
			}
			if job.Type == "" {
				continue
			}
			jobAsBytes, err := json.Marshal(job)
			if err != nil {
				log.Println(err)
				continue
			}
			if err = consumerProducer.Send(jobAsBytes); err != nil {
				log.Println(err)
				continue
			}
		case err = <-errorChannel:
			log.Println(err)
		}
	}
}
