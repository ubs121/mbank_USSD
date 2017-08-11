package ussd

import (
	xrpc "github.com/ubs121/rpc/xml"
	"log"
	"net/http"
	"time"
)

type (
	// ussd response to mobile operator
	UssdResponse struct {
		Session string `xml:"session"`
		Message string `xml:"message"`
		Type    string `xml:"type"`
		Msisdn  string `xml:"msisdn"`
	}

	// dummy response
	Response struct {
		Resp UssdResponse
	}
)

const (
	ParamMsisdn = "msisdn"
	ParamInput  = "input"
)

var (
	// holds all the session objects of users. This map will be cleared by sessionCleanup()
	sessionMap map[string]*Context

	// session timeout (in seconds)
	sessionTimeout time.Duration = time.Duration(60) * time.Second // 60 seconds

	MOBICOM_PREFIX = []string{"99", "75", "94", "95"}
	UNITEL_PREFIX  = []string{"88", "89", "86", "80"}
	SKYTEL_PREFIX  = []string{"91", "92", "50", "58", "96", "90"}
	GMOBILE_PREFIX = []string{"98", "55", "93", "97", "53"}
)

// xml-rpc controller
func ussdHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("xmlrpc << %s\n", r.URL.String())

	// filter
	if err := doFilter(w, r); err != nil {
		xrpc.WriteResponse(w, nil, err)
		return
	}

	// read params
	sessionId := r.FormValue("sessionId")
	msisdn := r.FormValue("source")
	input := r.FormValue("text")
	//serverName := r.RemoteAddr

	log.Println("<<", input)

	// remove 976
	if msisdn != "" && msisdn[0:3] == "976" {
		msisdn = msisdn[3:]
	}

	// check session id
	if sessionId == "" {
		xrpc.WriteResponse(w, nil, xrpc.FaultInvalidParams)
		return
	}

	ctx, exists := sessionMap[sessionId]

	if !exists {
		// create a new session
		ctx = NewContext(sessionId)
		ctx.Params[ParamMsisdn] = msisdn
		ctx.Params[ParamInput] = input

		// and put in the session map
		sessionMap[sessionId] = ctx

		log.Printf("New session created %s", sessionId)
	}

	// send to the processor
	retmsg, err := _send(ctx, input)

	log.Println(">>", retmsg)

	if err != nil {
		log.Printf("Error: %v", err)
		xrpc.WriteResponse(w, nil, err)
		return
	}

	resp := UssdResponse{Session: sessionId, Msisdn: msisdn}

	// terminate (end)
	if ctx.State == "." {
		// TODO: delete context ?
		resp.Type = "end"
		resp.Message = retmsg
		xrpc.WriteResponse(w, &Response{resp}, nil)
		return
	}

	// continue (cont)
	resp.Type = "cont"
	resp.Message = retmsg
	xrpc.WriteResponse(w, &Response{resp}, nil)
}

func doFilter(w http.ResponseWriter, r *http.Request) error {
	// TODO: IP filter

	return nil
}

// session cleanup
func sessionCleanup() {
	n := len(sessionMap)

	if n > 0 {

		for k, ss := range sessionMap {
			// timeout болсон session-үүдийг устгах
			if time.Now().Sub(ss.UpdateTime) > sessionTimeout {
				delete(sessionMap, k)
			}
		}

		log.Printf("%v session(s), %v cleaned up\n", n, n-len(sessionMap))
	}
}

func hasPrefix(msisdn string, pres []string) bool {
	for _, p := range pres {
		if msisdn[0:2] == p {
			return true
		}
	}
	return false
}

func RegisterService(mux *http.ServeMux) {
	// init session
	sessionMap = make(map[string]*Context)

	// start session cleaner
	ticker := time.NewTicker(sessionTimeout)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				// do stuff
				sessionCleanup()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	mux.HandleFunc("/mbank", ussdHandler)
}
