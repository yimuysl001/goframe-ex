package bus

import (
	"context"
	eventBus "github.com/asaskevich/EventBus"
)

var localEventBus = &LocalEventBus{bus: eventBus.New(), cancelMap: make(map[string]context.CancelFunc)}

type LocalEventBus struct {
	bus       eventBus.Bus
	handler   interface{}
	cancelMap map[string]context.CancelFunc
}
