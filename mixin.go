package observation

import (
	"github.com/AsynkronIT/protoactor-go/actor"
)

type observableMixin interface {
	Name() string
	init(actor.Context)
	subscribe(*actor.PID)
	unsubscribe(*actor.PID)
	Notify(interface{})
}

//Observable interface that actor must implement to work to be observable
type Observable interface {
	observableMixin

	CurrentObservation() interface{}
}

// enforces that Mixin implements observableMixin interface
var _ observableMixin = (*Mixin)(nil)

//Mixin add support for tracking observers and notifying them about changes
type Mixin struct {
	name          string
	subscriptions []*actor.PID
}

//Name return observable name
func (m *Mixin) Name() string {
	return m.name
}

func (m *Mixin) init(context actor.Context) {
	m.name = context.Self().Id
	m.subscriptions = []*actor.PID{}
}

func (m *Mixin) subscribe(pid *actor.PID) {
	if found, _ := containsPID(m.subscriptions, pid); !found {
		m.subscriptions = append(m.subscriptions, pid)
	}
}

func (m *Mixin) unsubscribe(pid *actor.PID) {
	if found, pos := containsPID(m.subscriptions, pid); found {
		m.subscriptions = append(m.subscriptions[:pos], m.subscriptions[pos+1:]...)
	}
}

//Notify send observation data to observers
func (m *Mixin) Notify(o interface{}) {
	for _, pid := range m.subscriptions {
		pid.Tell(&Observation{
			Producer: m.name,
			Data:     o,
		})
	}
}

func containsPID(pids []*actor.PID, pid *actor.PID) (bool, int) {
	for i, p := range pids {
		if pid == p {
			return true, i
		}
	}
	return false, -1
}
