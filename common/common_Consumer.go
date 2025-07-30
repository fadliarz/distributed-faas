package common

type Consumer[M any] interface {
	PollAndProcessMessages()
}
