package observation_test

import (
	"fmt"
	"testing"

	"github.com/AsynkronIT/protoactor-go/actor"
	observation "github.com/artyomturkin/protoactor-go-observation"
)

const (
	ObserverName   = "test.actor.observer"
	ObservableName = "test.actor.observable"
)

//ObservableActor simple observable actor for tests
type ObservableActor struct {
	observation.Mixin

	t    *testing.T
	data string
}

func (o *ObservableActor) CurrentObservation() interface{} {
	return o.data
}

// ensure ObservableActor implements observation.Observable interface
var _ observation.Observable = (*ObservableActor)(nil)

func (o *ObservableActor) Receive(ctx actor.Context) {
	switch m := ctx.Message().(type) {
	case string:
		o.t.Logf("observable - received %s", m)
		o.data = m
		o.Notify(o.data)
	}
}

//ObserverActor simple observer actor for tests
type ObserverActor struct {
	observations *[]string

	t *testing.T
}

type Observe struct {
	PID          *actor.PID
	Observations *[]string
}

func (o *ObserverActor) Receive(ctx actor.Context) {
	switch m := ctx.Message().(type) {
	case *Observe:
		o.observations = m.Observations

		m.PID.Tell(&observation.Subscribe{
			PID: ctx.Self(),
		})
		o.t.Logf("observer - observing %s", m.PID)
	case *observation.Observation:
		*o.observations = append(*o.observations, fmt.Sprintf("%s-%s", m.Producer, m.Data))
		o.t.Logf("observer - observed %s-%s", m.Producer, m.Data)
	default:
		o.t.Logf("observer - got message of type %T", m)
	}
}
