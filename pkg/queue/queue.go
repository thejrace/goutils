package queue

func GetDriver(options QueueWorkerOptions) QueueDriverContract {
	mapping := map[string]func(options QueueWorkerOptions) QueueDriverContract{
		"null": nullDriver,
		"sqs":  sqsDriver,
		"sync": syncDriver,
	}

	return mapping[options.Driver](options)
}

func sqsDriver(options QueueWorkerOptions) QueueDriverContract {
	return &SqsDriver{
		Options: options,
	}
}

func nullDriver(options QueueWorkerOptions) QueueDriverContract {
	return &NullDriver{
		Options: options,
	}
}

func syncDriver(options QueueWorkerOptions) QueueDriverContract {
	return &SyncDriver{
		Options: options,
	}
}

type QueueWorkerOptions struct {
	Driver    string
	QueueName string
	Manager   QueueManagerContract
}

type QueueManagerContract interface {
	ConstructJob(body string) (JobContract, error)
	MarkJobAsFailed(job JobContract, exception error) error
}
