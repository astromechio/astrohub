package astrolib

type ARes struct {
	ID           string
	ResponseChan chan ([]byte)
	ErrorChan    chan (error)
}

func (a *ARes) Response() ([]byte, error) {
	select {
	case res := <-a.ResponseChan:
		return res, nil
	case err := <-a.ErrorChan:
		return nil, err
	}
}
