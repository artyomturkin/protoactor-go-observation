package observation

import "github.com/AsynkronIT/protoactor-go/actor"

//Observation observed data and who was the producer
type Observation struct {
	Producer string
	Data     interface{}
}

//Subscribe add PID to notification pool
type Subscribe struct {
	PID *actor.PID
}

//Unsubscribe remove PID from notification pool
type Unsubscribe struct {
	PID *actor.PID
}
