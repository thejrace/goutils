package queue

import (
	"encoding/json"
	"fmt"
)

type JobContract interface {
	Handle() error
	GetJobName() string
	GetMessageId() string
	GetQueueName() string
	GetBodySerialized() string
	SetMessageId(messageId string)
}

type Job struct {
	MessageId *string        `json:"messageId"`
	JobName   string         `json:"jobName"`
	QueueName string         `json:"queueName"`
	Body      map[string]any `json:"body"`
}

func (j *Job) SetBody(body map[string]any) {
	j.Body = body
}

func (j *Job) SetJobName(jobName string) {
	j.JobName = jobName
}

func (j *Job) SetMessageId(messageId string) {
	j.MessageId = &messageId
}

func (j *Job) SetQueueName(queueName string) {
	j.QueueName = queueName
}

func (j *Job) GetBodySerialized() string {
	bodyJ, _ := json.Marshal(j.Body)

	return fmt.Sprintf("{\"jobName\": \"%s\", \"queueName\": \"%s\", \"body\": %s}", j.JobName, j.QueueName, bodyJ)
}

func (j *Job) GetJobName() string {
	return j.JobName
}

func (j *Job) GetQueueName() string {
	return j.QueueName
}

func (j *Job) GetMessageId() string {
	return *j.MessageId
}
