package main

import (
  "flag"
  "fmt"
  "net"
  "time"
  "encoding/xml"
)

type Pro5ConnectInfo struct {
  Host string
  Port int
  Password string
}

func main() {
  connectInfo := Pro5ConnectInfo{}

  var listenPort = flag.Int("port", 3000, "The port that this server listens on.")
  flag.StringVar(&connectInfo.Host, "pro5-host", "localhost", "The name of the pro5 server.")
  flag.IntVar(&connectInfo.Port, "pro5-port", 9002, "The port that pro5's remote stage display is listening on.")
  flag.StringVar(&connectInfo.Password, "pro5-password", "stage", "The password to the pro5 remote stage display")

  flag.Parse()

  fmt.Printf("listen on %d\n", *listenPort)

  var pro5, err = ConnectToPro5(connectInfo)
  var httpServer, err = StartServer(*listenPort, pro5)
}

//////////////////////////////
//
func StartServer(port int, pro5 Pro5Connection) {
  httpServer.Start(staticPages)
  httpServer.ListenForWebSocket(func(conn WebSocket) {
    pro5.SendUpdates(conn)
  })
}
func (c *Pro5Connection) SendUpdates(conn WebSocket) {
  conn <- c.DisplayLayouts
  conn <- c.InitialSlide
  c.OnSlide(func(newSlide interface{}) {
    conn <- newSlide
  })
}
//
//////////////////////////////

type Pro5Connection struct {
  Info Pro5ConnectInfo
  Connection net.Conn
}

func ConnectToPro5(info Pro5ConnectInfo) (*Pro5Connection, error) {
  var result = new(Pro5Connection)
  var err error
  result.Info = info
  result.Connection, err = net.Dial("tcp", fmt.Sprintf("%s:%d", info.Host, info.Port))
  if err != nil {
    return nil, err
  }
  go result.Run()
  return result, nil
}

func (c *Pro5Connection) Run() {
  var err error
  var xmlWriter = xml.NewEncoder(c.Connection)
  var loginElement = xml.StartElement{}
  loginElement.Name.Local = "StageDisplayLogin"
  xmlWriter.EncodeElement(c.Info.Password, loginElement)
  fmt.Fprintf(c.Connection, "\r\n")

  time.Sleep(1 * time.Second)

  var b []byte
  var i int
  i, err = c.Connection.Read(b)
  if err == nil {
    fmt.Printf("read %d bytes\n", i)
  } else {
    fmt.Printf("errror: %s\n", err)
  }
}
