package astrolib

type AQueue interface {
	QueueRequest(r *AReq) (*ARes, error)
	GetRequest(name string) (*AReq, error)
}

var sharedQueue *RequestQueue

type RequestQueue struct {
	Queues map[string]chan (*AReq)
}

func SharedRequestQueue() *RequestQueue {
	if sharedQueue == nil {
		sharedQueue = &RequestQueue{
			Queues: make(map[string]chan (*AReq)),
		}
	}

	return sharedQueue
}

func (rq *RequestQueue) QueueRequest(r *AReq) (*ARes, error) {
	reqID := r.ID

	res := &ARes{
		ID:           reqID,
		ResponseChan: make(chan []byte),
		ErrorChan:    make(chan error),
	}

	r.Response = res

	go rq.addRequestToQueue(r)

	return res, nil
}

func (rq *RequestQueue) addRequestToQueue(r *AReq) {
	q := rq.Queues[r.QueueName]

	if q != nil {
		q <- r
	} else {
		rq.Queues[r.QueueName] = make(chan *AReq)
		rq.Queues[r.QueueName] <- r
	}
}

func (rq *RequestQueue) GetRequest(name string) (*AReq, error) {
	q := rq.Queues[name]

	if q == nil {
		return nil, QueueError(QueueErrorTypeDoesNotExist)
	}

	req := <-q

	go rq.addResponseToMap(req.Response)

	return req, nil
}

func (rq *RequestQueue) addResponseToMap(r *ARes) error {
	resMap := SharedResponseMap()

	err := resMap.AddHandler(r)

	return err
}
