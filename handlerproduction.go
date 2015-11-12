package kwiscale

import "reflect"

// handlerManager is used to manage handler production.
type handlerManager struct {

	// the handler type to produce
	handler string

	// record closers
	closer chan int

	// chan that provides handler
	producer chan WebHandler
}

// newWebHandler produce a WebHandler from registry.
func (manager handlerManager) newWebHandler() WebHandler {
	defer func() {
		if err := recover(); err != nil {
			Error(err, handlerRegistry)
		}
	}()
	return reflect.New(handlerRegistry[manager.handler]).Interface().(WebHandler)
}

// produce returns the producer chan.
func (manager handlerManager) produce() <-chan WebHandler {
	return manager.producer
}

// produceHandlers continuously generates new handlers.
// The number of handlers to generate in cache is set by
// Config.NbHandlerCache.
func (manager handlerManager) produceHandlers() {
	// forever produce handlers until closer is filled
	for {
		select {
		case manager.producer <- manager.newWebHandler():
			Log("Produced handler ", manager.handler)
		case <-manager.closer:
			Log("Handler production closed")
			break
		}
	}
}
