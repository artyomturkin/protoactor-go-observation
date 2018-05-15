package observation

import (
	"log"
	"reflect"

	"github.com/AsynkronIT/protoactor-go/actor"
)

//Middleware handles initializing observable and managing subscriptions
func Middleware(next actor.ActorFunc) actor.ActorFunc {
	fn := func(ctx actor.Context) {
		switch m := ctx.Message().(type) {
		case *actor.Started:
			next(ctx)
			if p, ok := ctx.Actor().(Observable); ok {
				p.init(ctx)
			} else {
				log.Fatalf("Actor type %v is not observable", reflect.TypeOf(ctx.Actor()))
			}
		case *Subscribe:
			if p, ok := ctx.Actor().(Observable); ok {
				p.subscribe(m.PID)
				m.PID.Tell(&Observation{
					Producer: p.Name(),
					Data:     p.CurrentObservation(),
				})
			}
		case *Unsubscribe:
			if p, ok := ctx.Actor().(Observable); ok {
				p.unsubscribe(m.PID)
			}
		default:
			next(ctx)
		}
	}
	return fn
}
