package queue

type SyncDriver struct {
	Options QueueWorkerOptions
}

func (d *SyncDriver) Dispatch(job JobContract) error {
	return job.Handle()
}

func (d *SyncDriver) Process() error {
	return nil
}

func (d *SyncDriver) DeleteJob(messageHandle *string) error {
	return nil
}

func (d *SyncDriver) ClearQueue() error {
	return nil
}
