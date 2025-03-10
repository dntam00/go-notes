package pipeline

type request struct {
}

type response struct {
}

type PipeliningQueue interface {
	EnqueueRequest(req request) response
	NextRequestToSend() (req request, ok bool)
	ReplyToNextRequest(resp response)
}
