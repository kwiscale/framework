package kwiscale
import (
    "log"
    "time"
    "os"
    "fmt"
    "net/http"
)

// sessions handler
var sessions map[string]map[string]interface{}

// Runs in parallel, check Session TTL
func checkSessionsTTL (){
    for id, session := range sessions {
        duration := time.Since(session["__last_access"].(time.Time))
        if duration > GetConfig().SessionTTL {
            log.Printf("Kill session %s after %v", id, duration)
            delete(sessions, id)
        }
    }
    time.Sleep(time.Second *1)
    checkSessionsTTL()
}

// generate a session uuid
func (this *RequestHandler) genSessionID() {

	f, _ := os.Open("/dev/urandom")
	b := make([]byte, 16)
	f.Read(b)
	f.Close()
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	this.SessionId = uuid

}

// checkSessid checks if sessions is created and if user have session
// if not, generate a session then set cookie
func (this *RequestHandler) checkSessid() {
	if id, err := this.Request.Cookie(GetConfig().SessID); err == nil {
		this.SessionId = id.Value
	} else {
		this.genSessionID()
	}

	c := http.Cookie{
		Name:  GetConfig().SessID,
		Value: this.SessionId,
		Path:  "/",
	}

    s := sessions[this.SessionId]
	if s == nil {
		sessions[this.SessionId] = make(map[string]interface{})
		s = sessions[this.SessionId]
	}

    resetSessionTTL(this.SessionId)
	http.SetCookie(this.Response, &c)
}

//GetSession returns the value assign to key "name"
func (this *RequestHandler) GetSession(name string) interface{} {

	this.checkSessid()

	if s := sessions[this.SessionId]; s != nil {
		if r := s[name]; r != nil {
			return r
		}
	}

	return nil
}


// SetSession set val to name (key) for the current handled connection
func (this *RequestHandler) SetSession(name string, val interface{}) {

	this.checkSessid()

	if sessions == nil {
		sessions = make(map[string]map[string]interface{})
	}

	sessions[this.SessionId][name] = val
}

// resetSessionTTL set current date to __last_access session key
// see: checkSessionsTTL function
func resetSessionTTL (id string) {
    sess := sessions[id]
    sess["__last_access"] = time.Now()
}
