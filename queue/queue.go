package queue

// TODO: interface for queue

type KaiXinRing interface {
	Offer(v interface{}) bool
	Poll() (interface{}, bool)
	Size() uint64
}
