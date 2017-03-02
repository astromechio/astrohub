package astrolib

import uuid "github.com/satori/go.uuid"

type AReq struct {
	ID        string `json:"id"`
	QueueName string `json:"-"`
	Body      []byte `json:"originalRequest"`
	Response  *ARes  `json:"-"`
}

func NewReq(queueName string, body []byte) *AReq {
	newReq := &AReq{
		ID:        uuid.NewV4().String(),
		QueueName: queueName,
		Body:      body,
	}

	return newReq
}
