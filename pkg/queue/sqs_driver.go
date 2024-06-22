package queue

import (
	"fmt"
	"time"

	"goutils/pkg/cc"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
)

type SqsDriver struct {
	Options QueueWorkerOptions
}

func (d *SqsDriver) Dispatch(job JobContract) error {
	sess, err := session.NewSession()

	if err != nil {
		return err
	}

	svc := sqs.New(sess)

	queue := job.GetQueueName()

	urlResult, _ := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queue,
	})
	queueURL := urlResult.QueueUrl

	// Todo: handle regular and fifo queues separetly
	msgGroupId := "group"

	msgDeDuplicationId := uuid.New().String()

	_, err = svc.SendMessage(&sqs.SendMessageInput{
		MessageGroupId:         &msgGroupId,
		MessageDeduplicationId: &msgDeDuplicationId,
		MessageBody:            aws.String(job.GetBodySerialized()),
		QueueUrl:               queueURL,
	})

	if err != nil {
		return err
	}

	return nil
}

func (d *SqsDriver) Process() error {
	sess, err := session.NewSession()

	if err != nil {
		return err
	}

	svc := sqs.New(sess)

	urlResult, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &d.Options.QueueName,
	})

	if err != nil {
		return err
	}

	queueURL := urlResult.QueueUrl
	var timeout int64 = 30

	msgResult, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            queueURL,
		MaxNumberOfMessages: aws.Int64(10),
		VisibilityTimeout:   &timeout,
	})

	if err != nil {
		return err
	}

	for _, job := range msgResult.Messages {
		d.handleJob(job)
	}

	return nil
}

func (d *SqsDriver) handleJob(job *sqs.Message) {
	start := time.Now().UnixMilli()

	jobHandler, err := d.Options.Manager.ConstructJob(*job.Body)

	if err != nil {
		// Invalid job payload

		return
	}

	jobHandler.SetMessageId(*job.MessageId)
	cc.Warn(fmt.Sprintf("[%s][%s] Working", jobHandler.GetMessageId(), jobHandler.GetJobName()))

	err = jobHandler.Handle()

	if err != nil {
		markError := d.Options.Manager.MarkJobAsFailed(jobHandler, err)

		if markError != nil {
			cc.Err(fmt.Sprintf("Failed to mark failed job `%s`", err.Error()))
		}

		_ = d.DeleteJob(job.ReceiptHandle)

		return
	}

	_ = d.DeleteJob(job.ReceiptHandle)

	duration := time.Now().UnixMilli() - start
	cc.Ok(fmt.Sprintf("[%s][%s] Done.......%dms", jobHandler.GetMessageId(), jobHandler.GetJobName(), duration))
}

func (d *SqsDriver) DeleteJob(messageHandle *string) error {
	sess, err := session.NewSession()

	if err != nil {
		return err
	}

	svc := sqs.New(sess)

	queue := d.Options.QueueName

	urlResult, _ := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queue,
	})

	queueURL := urlResult.QueueUrl

	_, err = svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      queueURL,
		ReceiptHandle: messageHandle,
	})

	if err != nil {
		cc.Err(fmt.Sprintf("Can not delete job from queue `%s`", err.Error()))

		return err
	}

	return nil
}

func (d *SqsDriver) ClearQueue() error {
	sess, err := session.NewSession()

	if err != nil {
		return err
	}

	svc := sqs.New(sess)

	urlResult, _ := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &d.Options.QueueName,
	})

	_, err = svc.PurgeQueue(&sqs.PurgeQueueInput{
		QueueUrl: urlResult.QueueUrl,
	})

	if err != nil {
		return err
	}

	cc.Ok(fmt.Sprintf("'%s' queue cleared", d.Options.QueueName))

	return nil
}
