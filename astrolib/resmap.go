package astrolib

import "errors"

type ResponseMap struct {
	Responses map[string]chan ([]byte)
	Errors    map[string]chan (error)
}

var sharedMap *ResponseMap

func SharedResponseMap() *ResponseMap {
	if sharedMap == nil {
		sharedMap = &ResponseMap{
			Responses: make(map[string]chan ([]byte)),
			Errors:    make(map[string]chan (error)),
		}
	}

	return sharedMap
}

func (rm *ResponseMap) AddHandler(r *ARes) error {
	rm.Responses[r.ID] = r.ResponseChan
	rm.Errors[r.ID] = r.ErrorChan

	return nil
}

func (rm *ResponseMap) SendResponse(ID string, res []byte) error {
	q := rm.Responses[ID]

	if q == nil {
		return QueueError(QueueErrorTypeDoesNotExist)
	}

	q <- res

	rm.Responses[ID] = nil
	rm.Errors[ID] = nil

	return nil
}

func (rm *ResponseMap) SendError(ID string, err error) error {
	q := rm.Errors[ID]

	if q == nil {
		return QueueError(QueueErrorTypeDoesNotExist)
	}

	q <- err

	rm.Responses[ID] = nil
	rm.Errors[ID] = nil

	return nil
}

func (rm *ResponseMap) SendErrorString(ID string, err string) error {
	q := rm.Errors[ID]

	if q == nil {
		return QueueError(QueueErrorTypeDoesNotExist)
	}

	q <- errors.New(err)

	rm.Responses[ID] = nil
	rm.Errors[ID] = nil

	return nil
}
