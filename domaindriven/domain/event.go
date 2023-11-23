package domain

import "time"

type Event interface {
	Id() string
	OccurredOn() time.Time
}

type EventHandler interface {
	Handle(event Event)
}

type EventPublisher interface {
	Publish(event Event) error
}
