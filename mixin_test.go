package observation_test

import (
	"testing"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/artyomturkin/protoactor-go-observation"
)

func dlHandler(t *testing.T) func(interface{}) {
	return func(evt interface{}) {
		t.Logf("dead letter %t - val %v", evt, evt)
	}
}

func TestObservableActor(t *testing.T) {
	observations := &[]string{}

	actor.WithDeadLetterSubscriber(dlHandler(t))

	observableProps := actor.
		FromProducer(func() actor.Actor { return &ObservableActor{t: t} }).
		WithMiddleware(observation.Middleware)

	observerProps := actor.
		FromProducer(func() actor.Actor { return &ObserverActor{t: t} })

	observable, err := actor.SpawnNamed(observableProps, ObservableName)
	if err != nil {
		t.Fatalf("failed to spawn observable actor: %v", err)
	}

	observer, err := actor.SpawnNamed(observerProps, ObserverName)
	if err != nil {
		t.Fatalf("failed to spawn observer actor: %v", err)
	}

	//Check if state is propagated on subscription
	observer.Tell(&Observe{
		Observations: observations,
		PID:          observable,
	})

	//Wait for consistency
	time.Sleep(20 * time.Millisecond)
	//Change state
	observable.Tell("hello")

	//Wait for consistency
	time.Sleep(20 * time.Millisecond)
	//Change state agein
	observable.Tell("hello again")

	//Wait for consistency
	time.Sleep(20 * time.Millisecond)
	//unsubscribe observer and check that it does not receive new observations
	observable.Tell(&observation.Unsubscribe{
		PID: observer,
	})

	//Wait for consistency
	time.Sleep(20 * time.Millisecond)
	//Change state agein
	observable.Tell("hello one more time")

	//wait for observable and observer to shutdown
	observable.GracefulStop()
	observer.GracefulStop()

	o := *observations
	if len(o) != 3 {
		t.Fatalf("Did not get expected number of observations. Expected 3, got %d", len(*observations))
	}
	if o[0] != "test.actor.observable-" {
		t.Errorf("observed event 0 is not empty. Got: %s", o[0])
	}
	if o[1] != "test.actor.observable-hello" {
		t.Errorf("observed event 1 is not 'hello'. Got: %s", o[1])
	}
	if o[2] != "test.actor.observable-hello again" {
		t.Errorf("observed event 2 is not 'hello again'. Got: %s", o[2])
	}
}
