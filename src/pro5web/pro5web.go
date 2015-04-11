package pro5web

import (
  "fmt"
  "net/http"
  "runtime/pprof"

  "golang.org/x/net/websocket"
)

// Arguments for serving the web front end.
type WebServerInfo struct {
  Port int
}

type NewChannels chan<- (chan<- string)

func StartServer(info WebServerInfo, messageChannels NewChannels) {
  var mux = http.NewServeMux()
  mux.Handle("/connect", serveWebSockets(messageChannels))
  mux.HandleFunc("/_goroutines", dumpGoroutines)
  mux.Handle("/", http.FileServer(http.Dir("public/")))
  http.ListenAndServe(fmt.Sprintf(":%d", info.Port), mux)
  // This never returns.
}

func serveWebSockets(messageChannels NewChannels) websocket.Handler {
  return func(ws *websocket.Conn) {
    messages := make(chan string, 5)
    messageChannels <- messages
    for message := range messages {
      fmt.Fprintf(ws, "%s", message)
    }
  }
}

func dumpGoroutines(response http.ResponseWriter, request *http.Request) {
  response.Header().Set("Content-Type", "text/plain")
  goroutines := pprof.Lookup("goroutine")
  if goroutines == nil {
    response.WriteHeader(404)
    fmt.Fprintf(response, "Couldn't load goroutine profile.\r\n")
    return
  }
  goroutines.WriteTo(response, 1)
}
