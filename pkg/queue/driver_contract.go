package queue

type QueueDriverContract interface {
	Dispatch(job JobContract) error
	Process() error
	DeleteJob(messageHandle *string) error
	ClearQueue() error
}
