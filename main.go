package main

import (
	"flag"
	"log"
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"encoding/json"
	"time"
)

type RPCRequest struct {
	Method  string
	Params  []interface{}
	Id      int
	Jsonrpc string
}

type RPCResult struct {
	Id      uint        `json:"id"`
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
}

type HydraUser struct {
	EmailAddress string `json:"emailAddress"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}

type UserEvent struct {
	Email      string
	Event      string
	OccurredAt time.Time
}

func main() {
	portFlag := flag.Int("PORT", 9923, "tcp port to listen on")
	flag.Parse()

	log.Printf("hydra stub started on port %d...\n", *portFlag)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *portFlag), newRouter()); err != nil {
		panic(err)
	}
}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", makeRpcHandler()).Methods("POST")
	return r
}

func makeRpcHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		log.Printf("request: %s\n", string(body[:]))
		if err != nil {
			log.Printf("error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var rpcRequest RPCRequest
		if err := json.Unmarshal(body, &rpcRequest); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		} else {
			log.Printf("error: %v\n", err)
		}

		log.Printf("request: %s params: %s method: %s\n", rpcRequest.Method, rpcRequest.Params, rpcRequest.Method)

		var result RPCResult
		switch rpcRequest.Method {
		case "userExists":
			result = RPCResult{
				Id:      1,
				Jsonrpc: "2.0",
				Result:  true,
			}
		case "getUserDetails":
			result = RPCResult{
				Id:      1,
				Jsonrpc: "2.0",
				Result:  HydraUser{EmailAddress: "ralf.wirdemann@gmail.com", FirstName: "Ralf", LastName: "Wirdemann"},
			}
		case "getUserEvents":
			result = RPCResult{
				Id:      1,
				Jsonrpc: "2.0",
				//Result:  []UserEvent{{Email: fmt.Sprintf("ralf%d@gmail.com", rand.Intn(100-1)+1), Event: "created", OccurredAt: time.Now()}},
				Result:  []UserEvent{{Email: "ralf.wirdemann@gmail.com", Event: "created", OccurredAt: time.Now()}},
			}
		}

		b, _ := json.Marshal(result)
		log.Printf("Result: %s", string(b[:]))

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}
