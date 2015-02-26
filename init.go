package kwiscale

// initialize some plugins
func init() {
	RegisterTemplateEngine("basic", &Template{})
	RegisterSessionEngine("default", &SessionStore{})
}
