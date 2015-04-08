package main

import (
  "bytes"
  "container/list"
  "encoding/xml"
  "flag"
  "fmt"
  "io"
  "log"
  "net"
  "net/http"

  "golang.org/x/net/websocket"
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
  if err != nil {
    log.Fatal("ConnectToPro5: ", err)
  }
  StartServer(*listenPort, pro5)
}

func StartServer(port int, pro5 *Pro5Connection) {
  var mux = http.NewServeMux()
  mux.Handle("/connect", websocket.Handler(func(ws *websocket.Conn) { pro5.AddListener(ws) }))
  mux.Handle("/", http.FileServer(http.Dir("public/")))
  http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

type Pro5Connection struct {
  Info Pro5ConnectInfo
  Connection net.Conn
  Listeners list.List
  DisplayLayouts string
  LastSlide string
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
  var xmlWriter = xml.NewEncoder(c.Connection)
  var loginElement = xml.StartElement{}
  loginElement.Name.Local = "StageDisplayLogin"
  xmlWriter.EncodeElement(c.Info.Password, loginElement)
  fmt.Fprintf(c.Connection, "\r\n")

  go c.ReadEverything()
}

// <DisplayLayouts selected="Default">
//   <DisplayLayout showBorder="1E0" width="1680" identifier="Default" height="1050">
//     <Frame height="105.000000" width="336.000000" xAxis="84.000000" isVisible="YES" identifier="Clock" yAxis="0.000000"/>
//     <Frame height="105.000000" width="336.000000" xAxis="1260.000000" isVisible="YES" identifier="ElapsedTime" yAxis="0.000000"/>
//     <Frame height="105.000000" width="336.000000" xAxis="84.000000" isVisible="YES" identifier="Timer1" yAxis="787.500000"/>
//     <Frame height="105.000000" width="336.000000" xAxis="672.000000" isVisible="YES" identifier="Timer2" yAxis="787.500000"/>
//     <Frame height="105.000000" width="336.000000" xAxis="1260.000000" isVisible="YES" identifier="Timer3" yAxis="787.500000"/>
//     <Frame height="105.000000" width="336.000000" xAxis="672.000000" isVisible="YES" identifier="VideoCounter" yAxis="0.000000"/>
//     <Frame height="420.000000" width="336.000000" xAxis="1302.000000" isVisible="YES" identifier="ChordChart" yAxis="236.250000"/>
//     <Frame height="525.000000" width="672.000000" xAxis="42.000000" isVisible="YES" identifier="CurrentSlide" yAxis="131.250000" fontSize="60"/>
//     <Frame height="420.000000" width="504.000031" xAxis="756.000000" isVisible="YES" identifier="NextSlide" yAxis="183.750000" fontSize="60"/>
//     <Frame height="105.000000" width="672.000000" xAxis="42.000000" isVisible="YES" identifier="CurrentSlideNotes" yAxis="656.250000" fontSize="60"/>
//     <Frame height="105.000000" width="504.000031" xAxis="756.000000" isVisible="YES" identifier="NextSlideNotes" yAxis="603.750000" fontSize="60"/>
//     <Frame height="105.000000" width="1512.000000" xAxis="84.000000" isVisible="YES" identifier="Message" yAxis="918.750000" fontSize="60" flashColor="0.000000 1.000000 0.000000"/>
//   </DisplayLayout>
// </DisplayLayouts>
type DisplayLayouts struct {
  selected string `xml:",attr"`
  DisplayLayout DisplayLayout `xml:"DisplayLayout"`
}
type DisplayLayout struct {
  showBorder string `xml:",attr"`
  width string `xml:",attr"`
  identifier string `xml:",attr"`
  height string `xml:",attr"`
  Frame []Frame
}
type Frame struct {
  height string `xml:",attr"`
  width string `xml:",attr"`
  xAxis string `xml:",attr"`
  isVisible string `xml:",attr"`
  identifier string `xml:",attr"`
  yAxis string `xml:",attr"`
  fontSize string `xml:",attr"`
  flashColor string `xml:",attr"`
}

// <?xml version="1.0" encoding="UTF-8" standalone="no"?>
// <StageDisplayData>
//   <Fields>
//     <Field type="clock" clockFormat="1" label="Clock" identifier="Clock" alpha="1E0" red="1E0" green="1E0" blue="1E0">7:49:48 PM</Field>
//     <Field running="0" type="elapsed" label="Time Elapsed" identifier="ElapsedTime" alpha="1E0" red="1E0" green="1E0" blue="1E0">--:--:--</Field>
//     <Field running="0" type="countdown" overrun="0" label="Countdown 1" identifier="Timer1" alpha="1E0" red="1E0" green="1E0" blue="1E0">--:--:--</Field>
//     <Field running="0" type="countdown" overrun="0" label="Countdown 2" identifier="Timer2" alpha="1E0" red="1E0" green="1E0" blue="1E0">--:--:--</Field>
//     <Field running="0" type="countdown" overrun="0" label="Countdown 3" identifier="Timer3" alpha="1E0" red="1E0" green="1E0" blue="1E0">--:--:--</Field>
//     <Field type="slide" label="Current Slide" identifier="CurrentSlide" alpha="1E0" red="1E0" green="1E0" blue="1E0">You're the reason I sing
// The reason I sing
// Yes my heart will sing
// How I love You</Field>
//     <Field type="slide" label="Next Slide" identifier="NextSlide" alpha="1E0" red="1E0" green="1E0" blue="1E0">And forever I'll sing
// Forever I'll sing
// Yes my heart will sing
// How I love You</Field>
//     <Field type="slide" label="Current Slide Notes" identifier="CurrentSlideNotes" alpha="1E0" red="1E0" green="1E0" blue="1E0"/>
//     <Field type="slide" label="Next Slide Notes" identifier="NextSlideNotes" alpha="1E0" red="1E0" green="1E0" blue="1E0"/>
//     <Field type="message" label="Message" identifier="Message" alpha="1E0" red="1E0" green="1E0" blue="1E0"/>
//     <Field running="0" type="countdown" label="Video Countdown" identifier="VideoCounter" alpha="1E0" red="1E0" green="1E0" blue="1E0">--:--:--</Field>
//     <Field type="chordChart" label="Chord Chart" identifier="ChordChart"/>
//   </Fields>
// </StageDisplayData>
type StageDisplayData struct {
  // todo
}

func (c *Pro5Connection) ReadEverything() {
  var xmlReader = xml.NewDecoder(c.Connection)

  // http://blog.davidsingleton.org/parsing-huge-xml-files-with-go/
  for {
    var token, err = xmlReader.Token()
    if err != nil {
      log.Fatal("xmlReader.Token(): ", err)
    }
    switch se := token.(type) {
    case xml.StartElement:
      switch se.Name.Local {
      case "DisplayLayouts":
        c.DisplayLayouts, err = readXmlString(xmlReader, DisplayLayouts{}, &se)
        if err != nil {
          log.Fatal(err)
        } else {
          fmt.Println(c.DisplayLayouts)
        }
        sendToListeners(c, c.DisplayLayouts)
      case "StageDisplayData":
        c.LastSlide, err = readXmlString(xmlReader, StageDisplayData{}, &se)
        if err != nil {
          log.Fatal(err)
        } else {
          fmt.Println(c.LastSlide)
        }
        sendToListeners(c, c.LastSlide)
      }
    }
  }
}

func readXmlString(xmlReader *xml.Decoder, v interface{}, startElement *xml.StartElement) (string, error) {
  var err error

  var buffer bytes.Buffer
  var xmlWriter = xml.NewEncoder(&buffer)

  err = xmlWriter.EncodeToken(*startElement)
  if err != nil {
    return "", err
  }
  for depth := 1; depth > 0; {
    var token, err = xmlReader.Token()
    if err != nil {
      return "", err
    }
    switch token.(type) {
    case xml.StartElement:
      depth = depth + 1
    case xml.EndElement:
      depth = depth - 1
    }
    err = xmlWriter.EncodeToken(token)
    if err != nil {
      return "", err
    }
  }
  xmlWriter.Flush()
  return buffer.String(), nil
}

func (c *Pro5Connection) AddListener(listener io.Writer) {
  c.Listeners.PushBack(listener)
  sendToListener(c, listener, c.DisplayLayouts)
  sendToListener(c, listener, c.LastSlide)
}

func sendToListeners(c *Pro5Connection, payload string) (err error) {
  for e := c.Listeners.Front(); e != nil; e = e.Next() {
    listener, ok := e.Value.(io.Writer)
    if ok {
      err = sendToListener(c, listener, payload)
      if err != nil {
        return
      }
    }
  }
  return
}

func sendToListener(c *Pro5Connection, listener io.Writer, payload string) error {
  fmt.Printf("%T %s\n", listener, payload)
  if len(payload) > 0 {
    _, err := fmt.Fprintf(listener, "%s", payload)
    if err != nil {
      return err
    }
  }
  return nil
}
