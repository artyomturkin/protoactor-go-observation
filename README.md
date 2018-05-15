# Proto Actor [Go] - Observation

[![Go Report Card](https://goreportcard.com/badge/github.com/artyomturkin/protoactor-go-observation)](https://goreportcard.com/report/github.com/artyomturkin/protoactor-go-observation)

Go package that provides a mechanism for actors to receive notifications about changes in other actors.

Idea based on [Microsoft Orleans Observers](https://dotnet.github.io/orleans/Documentation/Core-Features/Observers.html)

## Get started

Install package:

```
go get github.com/artyomturkin/protoactor-go-observation
```

Actor, that will be observed, must implement `observation.Mixin` interface:
```go
type ExampleObservableActor struct {
    observation.Mixin

    data string
}

func (o *ExampleObservableActor) CurrentObservation() interface{} {
    return o.data
}
```

And send notifications about state changes:
```go
func (o *ExampleObservableActor) Receive(ctx actor.Context) {
    switch m := ctx.Message().(type) {
    case string:
        o.data = m
        o.Notify(o.data)
    }
}
```

To subscribe observer actor to notifications from observable actor, tell observable actor `&observation.Subscribe` message with `PID` of observer actor:
```go
observable.Tell(&observation.Subscribe{
    PID: observer,
})
```

Observer will receive observation data as a `*observation.Observation` object:
```go
func (p *Observer) Receive(ctx actor.Context) {
    switch m := ctx.Message().(type) {
    case *observation.Observation:
        fmt.Printf("received observation %v\n", m)
    }
}
```

## Full Example
```go
package main

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

func main() {
    observableProps := actor.
        FromProducer(func() actor.Actor {
            return &ExampleObservableActor{} 
        }).
        WithMiddleware(observation.Middleware)

    printerProps := actor.
        FromProducer(func() actor.Actor {
            return &Printer{} 
        })

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
}
```
Outputs:
```
received observation with value &{example.observable.actor }
received observation with value &{example.observable.actor hello}
```