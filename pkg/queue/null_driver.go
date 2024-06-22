package queue

type NullDriver struct {
	Options QueueWorkerOptions
}

func (d *NullDriver) Dispatch(job JobContract) error {
	return nil
}

func (d *NullDriver) Process() error {
	return nil
}

func (d *NullDriver) DeleteJob(messageHandle *string) error {
	return nil
}

func (d *NullDriver) ClearQueue() error {
	return nil
}
