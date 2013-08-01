package kwiscale

import (
	//"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	//"github.com/gosexy/gettext"
	"strings"
    "time"
)

// Configuration structure
type configStruct struct {
	// static dir
	Statics string

	// template dir
	Templates string

	// Session cookie name
	SessID string

    // Session Max Age (in seconds)
    SessionTTL time.Duration

	// Default language
	Lang string
}

type RouteMap struct {
	Route   *regexp.Regexp
	Handler func() IRequestHandler
}

// handlers stack
var globalhandlers []RouteMap
var config configStruct



// GetConfig returns a pointer to the server configuration
// if no configuration, returns this config:
//		configStruct {
//			Statics:   "./statics",
//			Templates: "./templates",
//			SessID:    "SESSID",
//			Lang:      "en_US",
//          SessionTTL: time.Second * 10  * 30 ,
//      }
func GetConfig() *configStruct {
	if config.Statics == "" {
		config = configStruct{
			Statics:   "./statics",
			Templates: "./templates",
			SessID:    "SESSID",
			Lang:      "en_US",
            SessionTTL: time.Second * 10  * 30 ,
		}
		/*
		   log.Println("config gettext")
		   log.Println(gettext.LC_ALL)

		   gettext.SetLocale(gettext.LC_ALL, "")
		   gettext.BindTextdomain("messages", "./locales")
		   gettext.Textdomain("messages")
		*/
	}

	return &config
}

// AddHandler adds handlers to the handlers stack.
func AddHandler(requests ...IRequestHandler) {
	for _, r := range requests {
		field, _ := reflect.TypeOf(r).Elem().FieldByName("RequestHandler")
		route := field.Tag.Get("route")
		//log.Printf("Append route: %s", route)
		reg := regexp.MustCompile(route)
		routemap := RouteMap{reg, r.New}
		globalhandlers = append(globalhandlers, routemap)
	}
}

// Serve start to serve on given address.
//
// Example:
//      kwiscale.Server(":8080")
func Serve(address string) {
    sessions = make(map[string]map[string]interface{})
    log.Println(sessions)
    go checkSessionsTTL()
	http.Handle("/statics/", http.StripPrefix("/statics", http.FileServer(http.Dir("./statics"))))
	http.HandleFunc("/", dispatch)
	http.ListenAndServe(address, nil)
}

// dispatch request to correct handler
func dispatch(w http.ResponseWriter, r *http.Request) {

	/*	defer func() {
			if err := recover(); err != nil {
				log.Println("ERROR", err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "%v", err)
			}
		}()
	*/
	r.ParseForm()
	rcall := r.URL.Path

	for _, route := range globalhandlers {
		if res := route.Route.FindStringSubmatch(rcall); len(res) > 0 {
			callMethod(res, route, w, r)
			//TODO: Why ? Tests tells it's ok, but real handlers seems to not have responses
			//go callMethod(res, route, w, r)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)

}

// callMethod will find correct Method to call from handler
// and append params needed. It check lang, session, and so on.
// If no method are defined to respond to Method, callMethod panics
func callMethod(res []string, route RouteMap, w http.ResponseWriter, r *http.Request) {

	// call constructor
	handler := route.Handler()

	// we've got some paramters
	if len(res) > 1 {
		// params captured
		handler.setParams(w, r, res[1:])
	} else {
		handler.setParams(w, r, nil)
	}

	//gettext.BindTextdomain("messages", "./locales/" + lang)
	_ = getLang(r, handler.GetSession("LANG"))

	switch r.Method {
	case "GET":
		handler.Get()
	case "POST":
		handler.Post()
	case "DELETE":
		handler.Delete()
	case "HEAD":
		handler.Head()
	case "PUT":
		handler.Put()
	default:
		panic("Method not found: " + r.Method)
	}
}

func getLang(r *http.Request, forced interface{}) string {
	if forced != nil {
		return forced.(string)
	}

	lang := r.Header["Accept-Language"]
	if len(lang) > 0 {
		language := GetBestMatchLang(lang[0])
		return strings.Replace(language, "-", "_", -1)
	}

	return GetConfig().Lang
}
