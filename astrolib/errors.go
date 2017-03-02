package astrolib

import "errors"

const (
	QueueErrorTypeDoesNotExist = iota
)

func QueueError(errType int) error {
	switch errType {
	case QueueErrorTypeDoesNotExist:
		return errors.New("Queue error: queue does not exist")
	default:
		return errors.New("Queue error: unknown error")
	}
}
