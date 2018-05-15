package observation_test

import (
	"fmt"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/artyomturkin/protoactor-go-observation"
)

//ExampleObservableActor simple observable actor for tests
type ExampleObservableActor struct {
	observation.Mixin

	data string
}

func (o *ExampleObservableActor) CurrentObservation() interface{} {
	return o.data
}

// ensure ObservableActor implements observation.Observable interface
var _ observation.Observable = (*ExampleObservableActor)(nil)

func (o *ExampleObservableActor) Receive(ctx actor.Context) {
	switch m := ctx.Message().(type) {
	case string:
		o.data = m
		o.Notify(o.data)
	}
}

//Printer
type Printer struct{}

func (p *Printer) Receive(ctx actor.Context) {
	if m, ok := ctx.Message().(*observation.Observation); ok {
		fmt.Printf("received observation with value %v\n", m)
	}
}

func ExampleObservable() {
	observableProps := actor.
		FromProducer(func() actor.Actor { return &ExampleObservableActor{} }).
		WithMiddleware(observation.Middleware)

	printerProps := actor.
		FromProducer(func() actor.Actor { return &Printer{} })

	observable, err := actor.SpawnNamed(observableProps, "example.observable.actor")
	if err != nil {
		fmt.Printf("failed to spawn observable actor: %v\n", err)
		panic(err)
	}

	printer, err := actor.SpawnNamed(printerProps, "printer")
	if err != nil {
		fmt.Printf("failed to spawn observer actor: %v\n", err)
		panic(err)
	}

	observable.Tell(&observation.Subscribe{
		PID: printer,
	})

	observable.Tell("hello")

	observable.Tell(&observation.Unsubscribe{
		PID: printer,
	})

	observable.Tell("hello again")

	//wait for async events to complete
	time.Sleep(1 * time.Millisecond)

	//stop all actors
	observable.GracefulStop()
	printer.GracefulStop()

	//Output:
	//received observation with value &{example.observable.actor }
	//received observation with value &{example.observable.actor hello}
}
