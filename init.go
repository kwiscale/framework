package kwiscale

func init() {
	RegisterTemplateEngine("basic", &Template{})
	RegisterSessionEngine("default", &SessionStore{})
}
