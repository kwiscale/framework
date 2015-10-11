package kwiscale

// initialize some plugins
func init() {
	RegisterTemplateEngine("basic", BuiltInTemplate{})
	RegisterSessionEngine("default", &CookieSessionStore{})
}
