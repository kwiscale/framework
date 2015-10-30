package kwiscale

import "log"

// Deprecated functions that will disapear

// GetApp returns the app that holds this handler.
func (b *BaseHandler) GetApp() *App {
	log.Println("[WARN] GetAp() is deprecated, please use App() method instead.")
	return b.App()
}
