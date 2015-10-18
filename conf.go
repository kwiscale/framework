package kwiscale

// Config structure that holds configuration
type Config struct {
	// Root directory where TemplateEngine will get files
	TemplateDir string
	// Port to listen
	Port string
	// Number of handler to prepare
	NbHandlerCache int
	// TemplateEngine to use (default, pango2...)
	TemplateEngine string
	// Template engine options (some addons need options)
	TemplateEngineOptions TplOptions

	// SessionEngine (default is a file storage)
	SessionEngine string
	// SessionName is the name of session, eg. Cookie name, default is "kwiscale-session"
	SessionName string
	// A secret string to encrypt cookie
	SessionSecret []byte
	// Configuration for SessionEngine
	SessionEngineOptions SessionEngineOptions

	// Static directory (to put css, images, and so on...)
	StaticDir string
	// Activate static in memory cache
	StaticCacheEnabled bool

	// StrictSlash allows to match route that have trailing slashes
	StrictSlash bool

	// Datastrore
	//DB        string
	//DBOptions DBOptions
}

// Initialize config default values if some are not defined
func initConfig(config *Config) *Config {
	if config == nil {
		config = new(Config)
	}

	if config.Port == "" {
		config.Port = ":8000"
	}

	if config.NbHandlerCache == 0 {
		config.NbHandlerCache = 5
	}

	if config.TemplateEngine == "" {
		config.TemplateEngine = "basic"
	}

	if config.TemplateEngineOptions == nil {
		config.TemplateEngineOptions = make(TplOptions)
	}

	if config.SessionEngine == "" {
		config.SessionEngine = "default"
	}
	if config.SessionName == "" {
		config.SessionName = "kwiscale-session"
	}
	if config.SessionSecret == nil {
		config.SessionSecret = []byte("A very long secret string you should change")
	}
	if config.SessionEngineOptions == nil {
		config.SessionEngineOptions = make(SessionEngineOptions)
	}

	return config
}

/*
type ymlDB struct {
	Engine string `yml:"engine"`
	Options DBOptions `yml:"options"`
}
*/

type ymlSession struct {
	Name    string               `yaml:"name,omitempty"`
	Engine  string               `yaml:"engine,omitempty"`
	Secret  []byte               `yaml:"secret,omitempty"`
	Options SessionEngineOptions `yaml:"options,omitempty"`
}

type ymlTemplate struct {
	Dir     string     `yaml:"dir,omitempty"`
	Engine  string     `yaml:"engine,omitempty"`
	Options TplOptions `yaml:"options,omitempty"`
}

type ymlRoute struct {
	Handler string `yaml:"handler"`
	Alias   string `yaml:"alias"`
}

// yamlConf is used to make yaml configuration easiest to write.
type yamlConf struct {
	Port               string              `yaml:"listen,omitempty"`
	NbHandlerCache     int                 `yaml:"nbhandler,omitempty"`
	StaticDir          string              `yaml:"staticdir,omitempty"`
	StaticCacheEnabled bool                `yaml:"staticcache,omitempty"`
	StrictSlash        bool                `yaml:"strictslash,omitempty"`
	Template           ymlTemplate         `yaml:"template,omitempty"`
	Session            ymlSession          `yaml:"session,omitempty"`
	Routes             map[string]ymlRoute `yaml:"routes"`
	//DB                 ymlDB               `yaml:"db,omitempty"`
}

// parse returns the *Config from yaml struct.
func (y yamlConf) parse() *Config {
	return &Config{
		Port:                  y.Port,
		NbHandlerCache:        y.NbHandlerCache,
		StaticDir:             y.StaticDir,
		StrictSlash:           y.StrictSlash,
		SessionEngine:         y.Session.Engine,
		SessionName:           y.Session.Name,
		SessionSecret:         y.Session.Secret,
		SessionEngineOptions:  y.Session.Options,
		TemplateDir:           y.Template.Dir,
		TemplateEngine:        y.Template.Engine,
		TemplateEngineOptions: y.Template.Options,
		//DB:                    y.DB.Engine,
		//DBOptions:             y.DB.Options,
	}
}
