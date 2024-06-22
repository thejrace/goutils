package queue

import (
	"encoding/json"
	"testing"
)

func TestJobSerialization(t *testing.T) {
	body := map[string]any{
		"type":  "meditation",
		"value": "testValue",
	}

	job := getJob("TestJob", "default", body)

	want := "{\"jobName\": \"TestJob\", \"queueName\": \"default\", \"body\": {\"type\":\"meditation\",\"value\":\"testValue\"}}"
	got := job.GetBodySerialized()

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestJobConstruction(t *testing.T) {
	body := "{\"jobName\": \"TestJob\", \"queueName\": \"default\", \"body\": {\"type\":\"meditation\",\"value\":\"testValue\"}}"

	var jobHandler Job

	err := json.Unmarshal([]byte(body), &jobHandler)

	if err != nil {
		t.Errorf("could not decode job")
	}

	job := &TestJob{Job: &jobHandler}

	if job.Body["type"] != "meditation" {
		t.Errorf("got %q want %q", job.Body["type"], "meditation")
	}

	if job.Body["value"] != "testValue" {
		t.Errorf("got %q want %q", job.Body["value"], "testValue")
	}

	if job.JobName != "TestJob" {
		t.Errorf("got %q want %q", job.JobName, "TestJob")
	}

	if job.QueueName != "default" {
		t.Errorf("got %q want %q", job.QueueName, "default")
	}

}

type TestJob struct {
	*Job
}

func (t *TestJob) Handle() error {
	return nil
}

func getJob(j string, q string, b map[string]any) JobContract {
	return &TestJob{
		Job: &Job{
			JobName:   j,
			QueueName: q,
			Body:      b,
		},
	}
}
