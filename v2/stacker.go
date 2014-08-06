package kwiscale

// stacker stacks several IHandler instances in a buffered channel that is returned.
// This channel is read in Server.Route() method.
// This way, handlers are instanciated before usage.
func stacker(factory Factory, stackSize int) <-chan IHandler {

	stack := make(chan IHandler, stackSize)

	// stack handlers in goroutine
	go func(f Factory, s chan IHandler) {
		for {
			s <- f()
		}
	}(factory, stack)
	return stack
}
